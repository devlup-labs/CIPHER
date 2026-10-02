// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title RaffleMath
 * @notice Pure mathematical library for ticket winner indexing, hashing, and Merkle tree verification.
 */
library RaffleMath {
    /**
     * @notice Computes the winner index given a blockhash, revealed secret, roundId, and round size tau.
     */
    function computeWinnerIndex(
        bytes32 bh,
        bytes32 secret,
        uint256 roundId,
        uint256 tau
    ) internal pure returns (uint256) {
        require(tau > 0, "RaffleMath: Tau must be positive");
        return uint256(keccak256(abi.encodePacked(bh, secret, roundId))) % tau;
    }

    /**
     * @notice Verifies a Merkle inclusion proof for a chunk at a given localIndex against a expected Merkle root.
     * @dev Uses double keccak256 leaf hashing (`keccak256(abi.encodePacked(localIndex, chunk))`) to prevent second pre-image attacks.
     */
    function verifyMerkleProof(
        bytes calldata chunk,
        uint256 localIndex,
        bytes32[] calldata proof,
        bytes32 root
    ) internal pure returns (bool) {
        bytes32 leaf = keccak256(abi.encodePacked(keccak256(abi.encodePacked(localIndex, chunk))));
        bytes32 computedHash = leaf;

        for (uint256 i = 0; i < proof.length; i++) {
            bytes32 proofElement = proof[i];
            if (computedHash <= proofElement) {
                computedHash = keccak256(abi.encodePacked(computedHash, proofElement));
            } else {
                computedHash = keccak256(abi.encodePacked(proofElement, computedHash));
            }
        }

        return computedHash == root;
    }

    /**
     * @notice Computes fallback ticket winning condition based on ticket win probability.
     */
    function isFallbackWinner(
        bytes32 bh,
        uint256 senderNonce,
        bytes32 secret,
        uint256 winProb
    ) internal pure returns (bool) {
        bytes32 result = keccak256(abi.encodePacked(bh, senderNonce, secret));
        return uint256(result) < winProb;
    }
}
