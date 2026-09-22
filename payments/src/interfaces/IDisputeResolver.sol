// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title IDisputeResolver
 * @notice Interface for on-chain chunk disputes, Merkle proof verification, and dispute timing resolution.
 */
interface IDisputeResolver {
    enum DisputeStatus {
        None,
        Open,
        ResolvedHonest,
        ResolvedSlashed
    }

    struct Dispute {
        DisputeStatus status;
        bytes32 evidenceHash;
        uint256 raisedBlock;
        uint256 deadlineBlock;
    }

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

    event DisputeResolvedSlashed(
        address indexed sender,
        address indexed recipient,
        uint256 indexed roundId,
        uint256 localIndex,
        uint256 slashedStake
    );

    error DisputeAlreadyOpen();
    error DisputeNotOpen();
    error DisputeWindowElapsed();
    error DisputeWindowNotElapsed();
    error InvalidMerkleProof();
    error NotChannelParty();
    error UnauthorizedCaller();

    function raiseChunkDispute(
        address recipient,
        uint256 roundId,
        uint256 localIndex,
        bytes32 evidenceHash
    ) external;

    function respondToDispute(
        address sender,
        uint256 roundId,
        uint256 localIndex,
        bytes calldata chunk,
        bytes32[] calldata merkleProof,
        bytes32 merkleRoot
    ) external;

    function resolveExpiredDispute(
        address sender,
        address recipient,
        uint256 roundId,
        uint256 localIndex
    ) external;

    function getDisputeStatus(
        address sender,
        address recipient,
        uint256 roundId,
        uint256 localIndex
    ) external view returns (DisputeStatus);

    function isDisputePending(
        address sender,
        address recipient,
        uint256 roundId
    ) external view returns (bool);
}
