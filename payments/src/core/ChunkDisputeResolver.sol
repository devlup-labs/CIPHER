// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IDisputeResolver} from "../interfaces/IDisputeResolver.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {ISettlementEngine} from "../interfaces/ISettlementEngine.sol";
import {IPaymentChannel} from "../interfaces/IPaymentChannel.sol";
import {ProtocolConfig} from "../config/ProtocolConfig.sol";
import {RaffleMath} from "../randomness/RaffleMath.sol";

/**
 * @title ChunkDisputeResolver
 * @notice Handles on-chain chunk delivery disputes, Merkle proof checks, and delinquent provider slashing.
 */
contract ChunkDisputeResolver is IDisputeResolver {
    ProtocolConfig public immutable config;
    IProviderRegistry public immutable providerRegistry;
    ISettlementEngine public immutable settlementEngine;
    IPaymentChannel public immutable paymentChannel;

    // key = keccak256(abi.encodePacked(sender, recipient, roundId, localIndex))
    mapping(bytes32 => Dispute) public disputes;

    // key = keccak256(abi.encodePacked(sender, recipient, roundId)) -> pending count
    mapping(bytes32 => uint256) public pendingDisputeCount;

    bytes32 public constant EXPIRED_DISPUTE_REASON = keccak256("EXPIRED_DISPUTE_SLASH");

    constructor(
        address _config,
        address _providerRegistry,
        address _settlementEngine,
        address _paymentChannel
    ) {
        require(_config != address(0), "ChunkDisputeResolver: Zero config");
        require(_providerRegistry != address(0), "ChunkDisputeResolver: Zero registry");
        require(_settlementEngine != address(0), "ChunkDisputeResolver: Zero settlement");
        require(_paymentChannel != address(0), "ChunkDisputeResolver: Zero channel");

        config = ProtocolConfig(_config);
        providerRegistry = IProviderRegistry(_providerRegistry);
        settlementEngine = ISettlementEngine(_settlementEngine);
        paymentChannel = IPaymentChannel(_paymentChannel);
    }

    function _disputeKey(
        address sender,
        address recipient,
        uint256 roundId,
        uint256 localIndex
    ) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(sender, recipient, roundId, localIndex));
    }

    function _roundDisputeKey(
        address sender,
        address recipient,
        uint256 roundId
    ) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(sender, recipient, roundId));
    }

    /**
     * @notice Client raises a chunk dispute on an open round.
     */
    function raiseChunkDispute(
        address recipient,
        uint256 roundId,
        uint256 localIndex,
        bytes32 evidenceHash
    ) external override {
        address sender = msg.sender;
        bytes32 dKey = _disputeKey(sender, recipient, roundId, localIndex);
        Dispute storage dispute = disputes[dKey];

        if (dispute.status != DisputeStatus.None) revert DisputeAlreadyOpen();

        ISettlementEngine.RoundStatus rStatus = settlementEngine.getRoundStatus(sender, recipient, roundId);
        if (rStatus != ISettlementEngine.RoundStatus.Open) {
            revert UnauthorizedCaller();
        }

        uint256 deadline = block.number + config.DISPUTE_RESPONSE_WINDOW();

        dispute.status = DisputeStatus.Open;
        dispute.evidenceHash = evidenceHash;
        dispute.raisedBlock = block.number;
        dispute.deadlineBlock = deadline;

        bytes32 rdKey = _roundDisputeKey(sender, recipient, roundId);
        pendingDisputeCount[rdKey] += 1;

        emit DisputeRaised(sender, recipient, roundId, localIndex, deadline);
    }

    /**
     * @notice Provider clears dispute by submitting challenged chunk and Merkle proof.
     */
    function respondToDispute(
        address sender,
        uint256 roundId,
        uint256 localIndex,
        bytes calldata chunk,
        bytes32[] calldata merkleProof,
        bytes32 merkleRoot
    ) external override {
        address recipient = msg.sender;
        bytes32 dKey = _disputeKey(sender, recipient, roundId, localIndex);
        Dispute storage dispute = disputes[dKey];

        if (dispute.status != DisputeStatus.Open) revert DisputeNotOpen();
        if (block.number > dispute.deadlineBlock) revert DisputeWindowElapsed();

        bool isValid = RaffleMath.verifyMerkleProof(chunk, localIndex, merkleProof, merkleRoot);
        if (!isValid) revert InvalidMerkleProof();

        dispute.status = DisputeStatus.ResolvedHonest;

        bytes32 rdKey = _roundDisputeKey(sender, recipient, roundId);
        if (pendingDisputeCount[rdKey] > 0) {
            pendingDisputeCount[rdKey] -= 1;
        }

        emit DisputeResolvedHonest(sender, recipient, roundId, localIndex);
    }

    /**
     * @notice Resolves expired dispute: slashes provider stake and voids the round.
     */
    function resolveExpiredDispute(
        address recipient,
        uint256 roundId,
        uint256 localIndex
    ) external override {
        // Anyone can call after deadline
        address sender = msg.sender;
        bytes32 dKey = _disputeKey(sender, recipient, roundId, localIndex);
        Dispute storage dispute = disputes[dKey];

        if (dispute.status != DisputeStatus.Open) revert DisputeNotOpen();
        if (block.number <= dispute.deadlineBlock) revert DisputeWindowNotElapsed();

        dispute.status = DisputeStatus.ResolvedSlashed;

        bytes32 rdKey = _roundDisputeKey(sender, recipient, roundId);
        if (pendingDisputeCount[rdKey] > 0) {
            pendingDisputeCount[rdKey] -= 1;
        }

        uint256 providerStake = providerRegistry.getProviderStake(recipient);
        uint256 slashAmount = (providerStake * config.SLASH_BPS()) / 10_000;

        uint256 slashedAmount = providerRegistry.slashProvider(recipient, slashAmount, EXPIRED_DISPUTE_REASON);
        settlementEngine.voidRound(sender, recipient, roundId, EXPIRED_DISPUTE_REASON);

        emit DisputeResolvedSlashed(sender, recipient, roundId, localIndex, slashedAmount);
    }

    function getDisputeStatus(
        address sender,
        address recipient,
        uint256 roundId,
        uint256 localIndex
    ) external view override returns (DisputeStatus) {
        bytes32 dKey = _disputeKey(sender, recipient, roundId, localIndex);
        return disputes[dKey].status;
    }

    function isDisputePending(
        address sender,
        address recipient,
        uint256 roundId
    ) external view override returns (bool) {
        bytes32 rdKey = _roundDisputeKey(sender, recipient, roundId);
        return pendingDisputeCount[rdKey] > 0;
    }
}
