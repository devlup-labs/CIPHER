// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IEntropySource} from "./IEntropySource.sol";

/**
 * @title ISettlementEngine
 * @notice Interface for round secret commitments, probabilistic ticket settlement, batch claims, and fallback claims.
 */
interface ISettlementEngine {
    enum RoundStatus {
        Uncommitted,
        Open,
        Settled,
        PartialFallback,
        Voided
    }

    struct Round {
        uint256 tau;
        uint256 faceValue;
        bytes32 recipientRandHash;
        uint256 commitBlock;
        RoundStatus status;
    }

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

    error RoundAlreadyCommitted();
    error RoundNotCommitted();
    error InvalidRoundParams();
    error RoundNotOpen();
    error RoundIsVoided();
    error NotWinningIndex();
    error TicketAlreadyUsed();
    error RoundNotPartial();
    error RoundNotStaleEnough();
    error RoundDisputePending();
    error UnauthorizedCaller();
    error ArrayLengthMismatch();

    function commitRoundSecret(
        address sender,
        uint256 roundId,
        uint256 tau,
        uint256 faceValue,
        bytes32 recipientRandHash
    ) external;

    function settleRound(
        IEntropySource.RoundTicket calldata ticket,
        bytes calldata senderSig,
        bytes32 secret
    ) external;

    function settleRoundBatch(
        IEntropySource.RoundTicket[] calldata tickets,
        bytes[] calldata senderSigs,
        bytes32[] calldata secrets
    ) external;

    function markRoundPartial(address sender, uint256 roundId) external;

    function claimFallbackTicket(
        IEntropySource.RoundTicket calldata ticket,
        bytes calldata senderSig,
        bytes32 secret
    ) external;

    function voidRound(address sender, address recipient, uint256 roundId, bytes32 reason) external;

    function getRoundStatus(address sender, address recipient, uint256 roundId) external view returns (RoundStatus);

    function previewWinnerIndex(
        address sender,
        address recipient,
        uint256 roundId,
        bytes32 secret
    ) external view returns (uint256);
}
