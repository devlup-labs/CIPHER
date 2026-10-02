// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";
import {ISettlementEngine} from "../interfaces/ISettlementEngine.sol";
import {IEntropySource} from "../interfaces/IEntropySource.sol";
import {IPaymentChannel} from "../interfaces/IPaymentChannel.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {IDisputeResolver} from "../interfaces/IDisputeResolver.sol";
import {ProtocolConfig} from "../config/ProtocolConfig.sol";
import {RaffleMath} from "../randomness/RaffleMath.sol";

/**
 * @title SettlementEngine
 * @notice Orchestrates round lifecycles, commit-reveal raffle draws, winning ticket claims, batch ticket settlements, and partial fallback claims.
 */
contract SettlementEngine is ISettlementEngine, ReentrancyGuard, Pausable {
    ProtocolConfig public immutable config;
    IEntropySource public immutable entropySource;
    IPaymentChannel public immutable paymentChannel;
    IProviderRegistry public immutable providerRegistry;
    IDisputeResolver public disputeResolver;

    // key = keccak256(abi.encodePacked(channelKey, roundId))
    mapping(bytes32 => Round) public rounds;

    // Nullifier for used ticket hashes in normal settlement
    mapping(bytes32 => bool) public usedTickets;

    // Nullifier for used ticket hashes in fallback mode
    mapping(bytes32 => bool) public usedFallbackTickets;

    address public owner;

    modifier onlyOwner() {
        if (msg.sender != owner) revert UnauthorizedCaller();
        _;
    }

    modifier onlyDisputeResolver() {
        if (address(disputeResolver) != address(0) && msg.sender != address(disputeResolver)) {
            revert UnauthorizedCaller();
        }
        _;
    }

    constructor(
        address _config,
        address _entropySource,
        address _paymentChannel,
        address _providerRegistry
    ) {
        require(_config != address(0), "SettlementEngine: Zero config");
        require(_entropySource != address(0), "SettlementEngine: Zero entropy");
        require(_paymentChannel != address(0), "SettlementEngine: Zero channel");
        require(_providerRegistry != address(0), "SettlementEngine: Zero registry");

        config = ProtocolConfig(_config);
        entropySource = IEntropySource(_entropySource);
        paymentChannel = IPaymentChannel(_paymentChannel);
        providerRegistry = IProviderRegistry(_providerRegistry);
        owner = msg.sender;
    }

    function setDisputeResolver(address _disputeResolver) external onlyOwner {
        disputeResolver = IDisputeResolver(_disputeResolver);
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    function _roundKey(bytes32 channelKey, uint256 roundId) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(channelKey, roundId));
    }

    /**
     * @notice Provider commits a new round on-chain before the draw block exists.
     */
    function commitRoundSecret(
        address sender,
        uint256 roundId,
        uint256 tau,
        uint256 faceValue,
        bytes32 recipientRandHash
    ) external override whenNotPaused {
        address recipient = msg.sender;

        if (tau == 0 || tau > config.MAX_ROUND_SIZE() || faceValue == 0) {
            revert InvalidRoundParams();
        }

        bytes32 channelKey = paymentChannel.getChannelKey(sender, recipient);
        bytes32 roundKey = _roundKey(channelKey, roundId);

        Round storage round = rounds[roundKey];
        if (round.status != RoundStatus.Uncommitted) {
            revert RoundAlreadyCommitted();
        }

        round.tau = tau;
        round.faceValue = faceValue;
        round.recipientRandHash = recipientRandHash;
        round.commitBlock = block.number;
        round.status = RoundStatus.Open;

        emit RoundCommitted(sender, recipient, roundId, tau, faceValue, block.number);
    }

    /**
     * @notice Internal settlement execution logic reused by single and batch settlement.
     */
    function _settleRoundInternal(
        IEntropySource.RoundTicket calldata ticket,
        bytes calldata senderSig,
        bytes32 secret
    ) internal {
        bytes32 channelKey = paymentChannel.getChannelKey(ticket.sender, ticket.recipient);
        bytes32 roundKey = _roundKey(channelKey, ticket.roundId);
        Round storage round = rounds[roundKey];

        if (round.status == RoundStatus.Uncommitted) revert RoundNotCommitted();
        if (round.status == RoundStatus.Voided) revert RoundIsVoided();
        if (round.status != RoundStatus.Open) revert RoundNotOpen();

        if (address(disputeResolver) != address(0)) {
            if (disputeResolver.isDisputePending(ticket.sender, ticket.recipient, ticket.roundId)) {
                revert RoundDisputePending();
            }
        }

        if (ticket.faceValue != round.faceValue || ticket.localIndex >= round.tau) {
            revert InvalidRoundParams();
        }

        bytes32 tHash = entropySource.roundTicketHash(ticket);
        if (usedTickets[tHash]) revert TicketAlreadyUsed();

        if (ticket.sender != entropySource.recoverSigner(tHash, senderSig)) {
            revert IEntropySource.InvalidSignature();
        }

        if (!entropySource.verifySecretCommitment(secret, round.recipientRandHash)) {
            revert IEntropySource.InvalidSecret();
        }

        (uint256 winnerIndex, ) = entropySource.computeWinnerIndex(
            round.commitBlock,
            secret,
            ticket.roundId,
            round.tau
        );

        if (ticket.localIndex != winnerIndex) {
            revert NotWinningIndex();
        }

        usedTickets[tHash] = true;
        round.status = RoundStatus.Settled;

        paymentChannel.payoutProvider(ticket.sender, ticket.recipient, ticket.faceValue);

        emit RoundSettled(ticket.sender, ticket.recipient, ticket.roundId, ticket.localIndex, ticket.faceValue);
    }

    /**
     * @notice Settles a single winning round ticket.
     */
    function settleRound(
        IEntropySource.RoundTicket calldata ticket,
        bytes calldata senderSig,
        bytes32 secret
    ) external override nonReentrant whenNotPaused {
        _settleRoundInternal(ticket, senderSig, secret);
    }

    /**
     * @notice Batch settles multiple winning round tickets in a single transaction for gas optimization.
     */
    function settleRoundBatch(
        IEntropySource.RoundTicket[] calldata tickets,
        bytes[] calldata senderSigs,
        bytes32[] calldata secrets
    ) external override nonReentrant whenNotPaused {
        if (tickets.length != senderSigs.length || tickets.length != secrets.length) {
            revert ArrayLengthMismatch();
        }

        for (uint256 i = 0; i < tickets.length; i++) {
            _settleRoundInternal(tickets[i], senderSigs[i], secrets[i]);
        }
    }

    /**
     * @notice Transitions a stalled or prematurely ended round into partial fallback mode.
     */
    function markRoundPartial(address sender, uint256 roundId) external override whenNotPaused {
        address recipient = msg.sender;
        bytes32 channelKey = paymentChannel.getChannelKey(sender, recipient);
        bytes32 roundKey = _roundKey(channelKey, roundId);
        Round storage round = rounds[roundKey];

        if (round.status != RoundStatus.Open) revert RoundNotOpen();

        bool isUnlock = paymentChannel.isUnlockInitiated(sender, recipient);
        bool isStale = block.number > round.commitBlock + config.ROUND_STALE_TIMEOUT();

        if (!isUnlock && !isStale) revert RoundNotStaleEnough();

        round.status = RoundStatus.PartialFallback;
        emit RoundMarkedPartial(sender, recipient, roundId);
    }

    /**
     * @notice Settles an individual ticket in partial fallback mode using independent win probability.
     */
    function claimFallbackTicket(
        IEntropySource.RoundTicket calldata ticket,
        bytes calldata senderSig,
        bytes32 secret
    ) external override nonReentrant whenNotPaused {
        bytes32 channelKey = paymentChannel.getChannelKey(ticket.sender, ticket.recipient);
        bytes32 roundKey = _roundKey(channelKey, ticket.roundId);
        Round storage round = rounds[roundKey];

        if (round.status != RoundStatus.PartialFallback) revert RoundNotPartial();

        bytes32 tHash = entropySource.roundTicketHash(ticket);
        if (usedFallbackTickets[tHash]) revert TicketAlreadyUsed();

        if (ticket.sender != entropySource.recoverSigner(tHash, senderSig)) {
            revert IEntropySource.InvalidSignature();
        }

        if (!entropySource.verifySecretCommitment(secret, round.recipientRandHash)) {
            revert IEntropySource.InvalidSecret();
        }

        (uint256 targetBlockNumber) = round.commitBlock + config.CONFIRMATION_DELAY();
        bytes32 bh = blockhash(targetBlockNumber);
        if (bh == bytes32(0)) revert IEntropySource.BlockhashUnavailable();

        if (!RaffleMath.isFallbackWinner(bh, ticket.senderNonce, secret, ticket.winProb)) {
            revert NotWinningIndex();
        }

        usedFallbackTickets[tHash] = true;
        paymentChannel.payoutProvider(ticket.sender, ticket.recipient, ticket.faceValue);

        emit FallbackTicketClaimed(ticket.sender, ticket.recipient, ticket.roundId, ticket.localIndex, ticket.faceValue);
    }

    /**
     * @notice Voids a round when a provider fails an on-chain dispute.
     */
    function voidRound(
        address sender,
        address recipient,
        uint256 roundId,
        bytes32 reason
    ) external override onlyDisputeResolver {
        bytes32 channelKey = paymentChannel.getChannelKey(sender, recipient);
        bytes32 roundKey = _roundKey(channelKey, roundId);
        Round storage round = rounds[roundKey];

        round.status = RoundStatus.Voided;
        emit RoundVoided(sender, recipient, roundId, reason);
    }

    function getRoundStatus(
        address sender,
        address recipient,
        uint256 roundId
    ) external view override returns (RoundStatus) {
        bytes32 channelKey = paymentChannel.getChannelKey(sender, recipient);
        bytes32 roundKey = _roundKey(channelKey, roundId);
        return rounds[roundKey].status;
    }

    function previewWinnerIndex(
        address sender,
        address recipient,
        uint256 roundId,
        bytes32 secret
    ) external view override returns (uint256) {
        bytes32 channelKey = paymentChannel.getChannelKey(sender, recipient);
        bytes32 roundKey = _roundKey(channelKey, roundId);
        Round storage round = rounds[roundKey];

        if (round.status == RoundStatus.Uncommitted) revert RoundNotCommitted();
        if (!entropySource.verifySecretCommitment(secret, round.recipientRandHash)) {
            revert IEntropySource.InvalidSecret();
        }

        (uint256 winnerIndex, ) = entropySource.computeWinnerIndex(
            round.commitBlock,
            secret,
            roundId,
            round.tau
        );

        return winnerIndex;
    }
}
