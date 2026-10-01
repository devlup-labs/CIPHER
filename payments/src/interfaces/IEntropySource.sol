// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title IEntropySource
 * @notice Interface for randomness, blockhash retrieval, EIP-712 ticket hashing, and signature verification.
 */
interface IEntropySource {
    struct RoundTicket {
        address sender;
        address recipient;
        uint256 roundId;
        uint256 localIndex;
        uint256 faceValue;
        uint256 winProb;
        uint256 senderNonce;
    }

    error RoundNotYetDrawable();
    error RoundDrawExpired();
    error BlockhashUnavailable();
    error InvalidSecret();
    error InvalidSignature();

    function computeWinnerIndex(
        uint256 commitBlock,
        bytes32 secret,
        uint256 roundId,
        uint256 tau
    ) external view returns (uint256 winnerIndex, bytes32 targetBlockhash);

    function roundTicketHash(RoundTicket calldata ticket) external view returns (bytes32);

    function recoverSigner(bytes32 ticketHash, bytes calldata sig) external view returns (address);

    function verifySecretCommitment(bytes32 secret, bytes32 committedHash) external pure returns (bool);
}
