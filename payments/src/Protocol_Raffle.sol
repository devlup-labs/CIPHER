// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract CipherBroker is ReentrancyGuard {
    // constants
    uint256 public constant CONFIRMATION_DELAY = 15;

    uint256 public constant MAX_BLOCKHASH_WINDOW = 256;

    uint256 public constant MIN_PROVIDER_STAKE = 1 ether;

    uint256 public constant UNBONDING_PERIOD = 7200; // ~1 day @ 12s blocks

    uint256 public constant CHANNEL_UNLOCK_PERIOD = 7200; // ~1 day

    uint256 public constant DISPUTE_RESPONSE_WINDOW = 1200; // ~4h @ 12s blocks

    uint256 public constant ROUND_STALE_TIMEOUT = 3600; // ~12h @ 12s blocks

    uint256 public constant MAX_ROUND_SIZE = 10_000;

    uint256 public constant SLASH_BPS = 2000; // 20%

    //enums

    enum RoundStatus {
        Uncommitted,
        Open,
        Settled,
        PartialFallback,
        Voided
    }

    enum DisputeStatus {
        None,
        Open,
        ResolvedHonest,
        ResolvedSlashed
    }

    // structs
    struct Provider {
        uint256 stake;
        uint256 unstakeReleaseBlock; // 0 = not unbonding // this stores the block until which the provider's stake cant be released, since if the provider is able to take back the stake immediately at any time,
        bool registered;
    }

    struct Channel {
        uint256 totalDeposited; // total ever deposited by sender
        uint256 settledAmount; // total paid out via winning tickets
        uint256 withdrawnAmount; // total reclaimed by sender after unlock
        uint256 unlockBlock; // 0 = not unlocking  // this stores the block until which the channel cant be closed, since if the client initiates the closing of the channel, the provider must have considerable time to redeem all the winning tkts, that he might have, and until then the sender has to be present to pay for those winning tkts, and once it is settled only then might the client withdraw
    }

    struct Ticket {
        address sender;
        address recipient;
        uint256 faceValue; // wei paid if this ticket wins
        uint256 winProb; // threshold out of type(uint256).max
        uint256 senderNonce;
        uint256 creationBlock;
        bytes32 recipientRandHash; // = keccak256(secret), fixed at service start
    }

    // the new structs for round_raffle

    struct Round {
        uint256 tau; // number of chunks in the round
        uint256 faceValue; // payout for winning ticket
        bytes32 recipientRandHash; // keccak256(secret)
        uint256 commitBlock; // block where round was committed
        RoundStatus status;
    }

    struct RoundTicket {
        address sender; // client
        address recipient; // provider
        uint256 roundId; // round this ticket belongs to
        uint256 localIndex; // chunk index inside this round
        uint256 faceValue; // must equal Round.faceValue
        uint256 winProb; // used only by fallback path
        uint256 senderNonce; // ticket uniqueness / fallback randomness
    }

    struct Dispute {
        DisputeStatus status;
        bytes32 evidenceHash;
        uint256 raisedBlock;
        uint256 deadlineBlock;
    }

    // mappings

    mapping(address => Provider) public providers;
    mapping(bytes32 => Channel) public channels; // key = _channelKey(sender, recipient) // this is just a map of identifying the channels based on a hash given to each channel, and that hash is being created using the sender, recipient of the channel
    mapping(bytes32 => bool) public usedTickets; // nullifier map, key = ticketHash

    // Key:
    //     _roundKey(channelKey, roundId)
    //
    // The round itself is scoped to a channel + roundId pair.
    mapping(bytes32 => Round) public rounds;

    // Disputes will be keyed by a hash representing the dispute.
    //
    // The exact dispute-key helper / lifecycle logic comes in
    // later milestones. For M1 we only create the storage layer.
    mapping(bytes32 => Dispute) public disputes;

    // Fallback tickets will eventually need their own nullifier
    // because multiple tickets from one partial round may win.
    //
    // This is implemented but not used yet in Milestone 1.
    mapping(bytes32 => bool) public usedFallbackTickets;

    // Events

    event ProviderRegistered(address indexed provider, uint256 totalStake);

    event ProviderUnstakeRequested(
        address indexed provider,
        uint256 releaseBlock
    );

    event ProviderStakeWithdrawn(address indexed provider, uint256 amount);

    event ChannelOpened(
        address indexed sender,
        address indexed recipient,
        uint256 amount
    );
    event ChannelDeposited(
        address indexed sender,
        address indexed recipient,
        uint256 amount
    );
    event ChannelUnlockInitiated(
        address indexed sender,
        address indexed recipient,
        uint256 unlockBlock
    );
    event ChannelWithdrawn(
        address indexed sender,
        address indexed recipient,
        uint256 amount
    );

    event TicketClaimed(
        bytes32 indexed ticketId,
        address indexed sender,
        address indexed recipient,
        uint256 faceValue
    );

    // round events

    event RoundCommitted(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId,
        uint256 tau,
        uint256 faceValue,
        uint256 commitBlock
    );

    event RoundSettled(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId,
        uint256 winningLocalIndex,
        uint256 faceValue
    );

    event RoundMarkedPartial(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId
    );

    event FallbackTicketClaimed(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId,
        uint256 localIndex,
        uint256 faceValue
    );

    event RoundVoided(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId,
        bytes32 reason
    );

    // dispute events

    event DisputeRaised(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId,
        uint256 localIndex,
        uint256 deadlineBlock
    );

    event DisputeResolvedHonest(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId,
        uint256 localIndex
    );

    event ProviderSlashed(
        address indexed provider,
        uint256 amount,
        bytes32 indexed reasonId
    );

    // Errors

    error InsufficientStake();
    error NotRegistered();
    error UnstakeAlreadyRequested();
    error UnstakeNotRequested();
    error UnbondingNotComplete();
    error ZeroDeposit();
    error UnlockAlreadyInitiated();
    error UnlockNotInitiated();
    error UnlockNotComplete();
    error NothingToWithdraw();
    error TicketAlreadyUsed();
    error NotRecipient();
    error InvalidSecret();
    error InvalidSignature();
    error TooEarly();
    error TicketExpired();
    error BlockhashUnavailable();
    error NotAWinner();
    error InsufficientChannelBalance();
    error TransferFailed();

    // round errors

    error RoundAlreadyCommitted();
    error RoundNotCommitted();
    error InvalidRoundParams();
    error RoundNotYetDrawable();
    error RoundDrawExpired();
    error NotWinningIndex();
    error RoundNotOpen();
    error RoundNotPartial();
    error RoundDisputePending();
    error RoundIsVoided();
    error RoundNotStaleEnough();

    // dispute errors

    error DisputeAlreadyOpen();
    error DisputeNotOpen();
    error DisputeWindowElapsed();
    error DisputeWindowNotElapsed();
    error InvalidMerkleProof();
    error NotChannelParty();

    // this basically prevents from entering the same function again until the funciton execution has completed fully
    // very similar to mutex in go, just that in go it was preventing from parallel execution, here there is noo parallel execution, but here also it has a similar problem of recursive execution

    // provider functions

    function registerProvider() external payable {
        // using external instead of public slightly better in gas, both work, but obv wont call this func from inside this contract
        if (((providers[msg.sender].stake + msg.value)) < MIN_PROVIDER_STAKE) {
            revert InsufficientStake();
        }

        Provider storage p = providers[msg.sender];
        p.stake += msg.value;
        p.registered = true;
        p.unstakeReleaseBlock = 0;
        emit ProviderRegistered(msg.sender, p.stake);
    }

    function requestUnstake() external {
        Provider storage p = providers[msg.sender];

        if (!p.registered) revert NotRegistered();

        if (p.unstakeReleaseBlock != 0) {
            revert UnstakeAlreadyRequested();
        }

        p.unstakeReleaseBlock = block.number + UNBONDING_PERIOD;
        emit ProviderUnstakeRequested(msg.sender, p.unstakeReleaseBlock);
    }

    function withdrawStake() external nonReentrant {
        Provider storage p = providers[msg.sender];

        // if there was no unstake req
        if (p.unstakeReleaseBlock == 0) {
            revert UnstakeNotRequested();
        }
        // if the unbonding period not yet complete
        if (p.unstakeReleaseBlock > block.number) {
            revert UnbondingNotComplete();
        }

        // Insepct this block of code again !!
        // Kuch dekha dekha lag rha????
        // THE FAMOUS 2016 DAO HACK !!!!!!
        // Thus always remember CEI (Check, Effects, Interact), the interact step should be at the last, since interact gives the control to the reciever, which might not just be an EOA, it can be any external smart contract, and when sending eth, it will trigger its recieve() function, and it can execute any arbitrary code, so REMEMBER THATTT

        // -----------------------------------------------------------
        uint256 amount = p.stake; // reset krne ke baad log me bhejne ke liye kahi store krke rakhna padega na

        (bool callSuccess, ) = payable(msg.sender).call{value: p.stake}("");
        if (!callSuccess) revert TransferFailed();

        p.stake = 0;
        p.registered = false;
        p.unstakeReleaseBlock = 0;

        // -----------------------------------------------------------

        // even though here the nonreentrancy is done so technically no problem, but a world famous flaw

        // also u might think that in the CEI model, if the interaction step fails, then there might be an acutal issue, but remember in case of any failure we revert ie rollback to the previous state, thus the state becomes exactly identical to before the ransaction

        emit ProviderStakeWithdrawn(msg.sender, amount);
    }

    // channel functions

    // this just creates a key for a channel based on its sender and recipient
    function _channelKey(
        address sender,
        address recipient
    ) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(sender, recipient)); //can use only abi.encode() too, bas "efficiency";
    }

    // abi.encode() pads everything to 32-byte boundaries while abi.encodePacked() just concatenates all the given thingss

    // REMEMBER reciever opens the channel, bcoz he is the one putting in the money in the channel, and the channel will also be closed by the reciever only
    function openChannel(address recipient) external payable {
        // there is no sense of channel without money in it
        if (msg.value == 0) {
            revert ZeroDeposit();
        }

        bytes32 ChannelKey = _channelKey(msg.sender, recipient);
        Channel storage ch = channels[ChannelKey];

        ch.totalDeposited += msg.value;

        // no need to define these, these are uint256, they are default to 0
        // ch.settledAmount = 0;
        // ch.withdrawnAmount = 0;
        // ch.unlockBlock = 0;

        emit ChannelOpened(msg.sender, recipient, msg.value);
    }

    function depositToChannel(address recipient) external payable {
        if (msg.value == 0) {
            revert ZeroDeposit();
        }

        bytes32 ChannelKey = _channelKey(msg.sender, recipient);
        Channel storage ch = channels[ChannelKey];

        ch.totalDeposited += msg.value;
        emit ChannelDeposited(msg.sender, recipient, msg.value);
    }

    function initiateChannelUnlock(address recipient) external {
        bytes32 ChannelKey = _channelKey(msg.sender, recipient);
        Channel storage ch = channels[ChannelKey];

        if (ch.unlockBlock != 0) revert UnlockAlreadyInitiated();
        ch.unlockBlock = block.number + CHANNEL_UNLOCK_PERIOD;
        emit ChannelUnlockInitiated(msg.sender, recipient, ch.unlockBlock);
    }

    // ek baar check krna ki do i need to still keep the withdrawn amt as a var in the channel struct, like do i need to give the reciever an option to withdraw partial amt, or while withdrawing does it have to always be the full amt, and just decide that
    function withdrawChannel(address recipient) external nonReentrant {
        bytes32 ChannelKey = _channelKey(msg.sender, recipient);
        Channel storage ch = channels[ChannelKey];

        // channel unlock not initiated only
        if (ch.unlockBlock == 0) {
            revert UnlockNotInitiated();
        }
        // unlockblock not yet reached
        if (ch.unlockBlock > block.number) {
            revert UnlockNotComplete();
        }

        uint256 remaining = ch.totalDeposited -
            ch.settledAmount -
            ch.withdrawnAmount;

        if (remaining == 0) revert NothingToWithdraw();
        ch.withdrawnAmount += remaining;

        (bool callSuccess, ) = payable(msg.sender).call{value: remaining}("");
        if (!callSuccess) revert TransferFailed();
        emit ChannelWithdrawn(msg.sender, recipient, remaining);
    }

    // ticket functions

    // every security-relevant field should be included inn computing the hash which would then be signed by the sender, since any change in it should change the hash

    // this hash would be signed by the sender also used as the nullifier key

    function ticketHash(Ticket calldata t) public view returns (bytes32) {
        // obv u wont use storage for the ticket struct passed on, also since im not modifying it, so need of memory, thus calldata, which directly just reads from the given Ticket struct
        return
            keccak256(
                abi.encode(
                    address(this), // this ensures that "This ticket is only valid for this specific deployed contract"
                    block.chainid, // so that the same signed tkt cant be used on another EVM based chain ie "cross-chain replay"
                    t.sender,
                    t.recipient,
                    t.faceValue,
                    t.winProb,
                    t.senderNonce, // to diff txn having everything same, but is a diff txn
                    t.creationBlock, // for ensuring that the future block gets fixed
                    t.recipientRandHash // obv so that the recipient gets commited
                )
            );
    }

    // this is the new func for tkt hash, with the newer fields, like roundId, localIndex
    function roundTicketHashFunc(
        RoundTicket calldata ticket
    ) public view returns (bytes32) {
        return
            keccak256(
                abi.encode(
                    address(this),
                    block.chainid,
                    ticket.sender,
                    ticket.recipient,
                    ticket.faceValue,
                    ticket.winProb,
                    ticket.senderNonce,
                    ticket.roundId,
                    ticket.localIndex
                )
            );
    }

    // since this function uses inbuilt EVM function of ecrecover, this is vulnerable to signature malleability, and this will upgrade to openzeppelin's lib in the future, but as of now since using nullifier, its protected

    /// @dev Recovers the signer of an EIP-191-prefixed ticket hash using
    ///      OpenZeppelin's ECDSA library, which additionally rejects
    ///      malleable (high-s) signatures and malformed signature bytes
    ///      with explicit, gas-efficient custom errors.
    function _recoverSigner(
        bytes32 tHash,
        bytes calldata sig
    ) internal pure returns (address) {
        bytes32 ethSignedHash = MessageHashUtils.toEthSignedMessageHash(tHash);
        return ECDSA.recover(ethSignedHash, sig);
    }

    function claimWinningTicket(
        Ticket calldata ticket,
        bytes calldata senderSig,
        bytes32 secret
    ) external nonReentrant {
        if (msg.sender != ticket.recipient) {
            revert NotRecipient();
        }

        bytes32 tkthash = ticketHash(ticket);

        if (usedTickets[tkthash] == true) {
            revert TicketAlreadyUsed();
        }

        // dont check this with the msg.sender bcoz this func would be called on by the provider who aims to get money, not the client, who is gonna sign the tkt
        if (ticket.sender != _recoverSigner(tkthash, senderSig)) {
            revert InvalidSignature();
        }
        if (keccak256(abi.encodePacked(secret)) != ticket.recipientRandHash) {
            revert InvalidSecret();
        }

        // future block thingsss
        uint256 targetBlock = ticket.creationBlock + CONFIRMATION_DELAY;
        if (block.number < targetBlock) {
            revert TooEarly();
        }
        if (block.number > ticket.creationBlock + MAX_BLOCKHASH_WINDOW) {
            revert TicketExpired();
        }

        bytes32 bh = blockhash(targetBlock);

        if (bh == bytes32(0)) {
            revert BlockhashUnavailable();
        }

        // THIS IS THE MAIN THING THAT WAS DECIDED AFTER A BIGGG LONG TIMEE, which has entropy from all the components interacting in here, the sender, the reciever, the future block hash (ethereum itself)
        bytes32 result = keccak256(
            abi.encodePacked(bh, ticket.senderNonce, secret)
        );
        if (uint256(result) >= ticket.winProb) revert NotAWinner();

        usedTickets[tkthash] = true;

        // paisa bhejne se pehle check toh krle channel me utne paise hai bhi ki nhiiii !!

        bytes32 chKey = _channelKey(ticket.sender, ticket.recipient);
        Channel storage ch = channels[chKey];

        // uint256 available = ch.totalDeposited - ch.settledAmount -ch.withdrawnAmount; -- DESIGN CHOICE
        uint256 available = ch.totalDeposited -
            ch.settledAmount -
            ch.withdrawnAmount;
        if (ticket.faceValue > available) revert InsufficientChannelBalance();

        ch.settledAmount += ticket.faceValue;

        (bool callSuccess, ) = payable(msg.sender).call{
            value: ticket.faceValue
        }("");
        if (!callSuccess) revert TransferFailed();

        emit TicketClaimed(
            tkthash,
            ticket.sender,
            ticket.recipient,
            ticket.faceValue
        );
    }

    function _roundKey(
        bytes32 channelKey,
        uint256 roundId
    ) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(channelKey, roundId));
    }

    function commitRoundSecret(
        address sender,
        uint256 roundId,
        uint256 tau,
        uint256 faceValue,
        bytes32 recipientRandHash
    ) external {
        // Provider is msg.sender.
        address recipient = msg.sender;

        if (tau == 0 || tau > MAX_ROUND_SIZE || faceValue == 0) {
            revert InvalidRoundParams();
        }

        // Identify the channel

        bytes32 channelKey = _channelKey(sender, recipient);

        // Make sure this roundId has not already been committed
        // for this particular channel.
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

        emit RoundCommitted(
            sender,
            recipient,
            roundId,
            tau,
            faceValue,
            block.number
        );
    }

    function _computeWinnerIndex(
        bytes32 bh,
        bytes32 secret,
        uint256 roundId,
        uint256 tau
    ) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(bh, secret, roundId))) % tau;
    }

    function previewWinnerIndex(
        address sender,
        uint256 roundId,
        bytes32 secret
    ) external view returns (uint256) {
        address recipient = msg.sender;

        bytes32 channelKey = _channelKey(sender, recipient);

        bytes32 roundKey = _roundKey(channelKey, roundId);

        Round storage round = rounds[roundKey];

        if (round.status == RoundStatus.Uncommitted) {
            revert RoundNotCommitted();
        }

        uint256 targetBlock = round.commitBlock + CONFIRMATION_DELAY;

        // Draw block has not happened yet.
        if (block.number < targetBlock) {
            revert RoundNotYetDrawable();
        }

        // Keep the same historical blockhash window as the existing ticket mechanism.
        if (block.number > round.commitBlock + MAX_BLOCKHASH_WINDOW) {
            revert RoundDrawExpired();
        }

        // Verify that the revealed secret is the one
        // committed by the provider.
        if (keccak256(abi.encodePacked(secret)) != round.recipientRandHash) {
            revert InvalidSecret();
        }

        bytes32 bh = blockhash(targetBlock);

        if (bh == bytes32(0)) {
            revert BlockhashUnavailable();
        }

        return _computeWinnerIndex(bh, secret, roundId, round.tau);
    }

    function settleRound(
        RoundTicket calldata ticket,
        bytes calldata senderSig,
        bytes32 secret
    ) external nonReentrant {
        // 1. Identify round

        bytes32 channelKey = _channelKey(ticket.sender, ticket.recipient);

        bytes32 roundKey = _roundKey(channelKey, ticket.roundId);

        Round storage round = rounds[roundKey];

        // 2. Round must exist and be Open

        if (round.status == RoundStatus.Uncommitted) {
            revert RoundNotCommitted();
        }

        if (round.status != RoundStatus.Open) {
            if (round.status == RoundStatus.Voided) {
                revert RoundIsVoided();
            }

            revert RoundNotOpen();
        }

        // 3. Verify ticket belongs to this round

        if (ticket.faceValue != round.faceValue) {
            revert InvalidRoundParams();
        }

        if (ticket.localIndex >= round.tau) {
            revert InvalidRoundParams();
        }

        // 4. Verify sender signature

        bytes32 roundTicketHash = roundTicketHashFunc(ticket);

        if (ticket.sender != _recoverSigner(roundTicketHash, senderSig)) {
            revert InvalidSignature();
        }

        // 5. Verify committed secret

        if (keccak256(abi.encodePacked(secret)) != round.recipientRandHash) {
            revert InvalidSecret();
        }

        // 6. Check draw timing

        uint256 targetBlock = round.commitBlock + CONFIRMATION_DELAY;

        if (block.number < targetBlock) {
            revert RoundNotYetDrawable();
        }

        if (block.number > round.commitBlock + MAX_BLOCKHASH_WINDOW) {
            revert RoundDrawExpired();
        }

        // 7. Get future blockhash

        bytes32 bh = blockhash(targetBlock);

        if (bh == bytes32(0)) {
            revert BlockhashUnavailable();
        }

        // 8. Compute winner

        uint256 winnerIndex = _computeWinnerIndex(
            bh,
            secret,
            ticket.roundId,
            round.tau
        );

        if (ticket.localIndex != winnerIndex) {
            revert NotWinningIndex();
        }

        // 9. Check channel balance

        Channel storage ch = channels[channelKey];

        uint256 available = ch.totalDeposited -
            ch.settledAmount -
            ch.withdrawnAmount;

        if (ticket.faceValue > available) {
            revert InsufficientChannelBalance();
        }

        // 10. Settle round

        round.status = RoundStatus.Settled;

        ch.settledAmount += ticket.faceValue;

        // ============================================================
        // 11. Pay provider
        // ============================================================

        (bool callSuccess, ) = payable(ticket.recipient).call{
            value: ticket.faceValue
        }("");

        if (!callSuccess) {
            revert TransferFailed();
        }

        emit RoundSettled(
            ticket.sender,
            ticket.recipient,
            ticket.roundId,
            ticket.localIndex,
            ticket.faceValue
        );
    }

    // View helpers

    function channelBalance(
        address sender,
        address recipient
    ) external view returns (uint256) {
        Channel storage ch = channels[_channelKey(sender, recipient)];
        return ch.totalDeposited - ch.settledAmount - ch.withdrawnAmount;
    }

    function isProviderActive(address provider) external view returns (bool) {
        Provider storage p = providers[provider];
        return
            p.registered &&
            p.stake >= MIN_PROVIDER_STAKE &&
            p.unstakeReleaseBlock == 0;
    }
}
