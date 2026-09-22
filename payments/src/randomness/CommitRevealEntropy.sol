// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {EIP712} from "@openzeppelin/contracts/utils/cryptography/EIP712.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {IEntropySource} from "../interfaces/IEntropySource.sol";
import {RaffleMath} from "./RaffleMath.sol";
import {ProtocolConfig} from "../config/ProtocolConfig.sol";

/**
 * @title CommitRevealEntropy
 * @notice Production randomness module isolating commit-reveal blockhash entropy and EIP-712 structured data ticket signature verification.
 */
contract CommitRevealEntropy is IEntropySource, EIP712 {
    using ECDSA for bytes32;

    ProtocolConfig public immutable config;

    bytes32 public constant ROUND_TICKET_TYPEHASH = keccak256(
        "RoundTicket(address sender,address recipient,uint256 roundId,uint256 localIndex,uint256 faceValue,uint256 winProb,uint256 senderNonce)"
    );

    constructor(address _config) EIP712("CIPHER Payment Protocol", "1.0.0") {
        require(_config != address(0), "CommitRevealEntropy: Zero config address");
        config = ProtocolConfig(_config);
    }

    /**
     * @notice Computes winner index for a given round based on target blockhash and revealed secret.
     */
    function computeWinnerIndex(
        uint256 commitBlock,
        bytes32 secret,
        uint256 roundId,
        uint256 tau
    ) external view override returns (uint256 winnerIndex, bytes32 targetBlockhash) {
        uint256 targetBlock = commitBlock + config.CONFIRMATION_DELAY();

        if (block.number < targetBlock) {
            revert RoundNotYetDrawable();
        }

        if (block.number > commitBlock + config.MAX_BLOCKHASH_WINDOW()) {
            revert RoundDrawExpired();
        }

        targetBlockhash = blockhash(targetBlock);
        if (targetBlockhash == bytes32(0)) {
            revert BlockhashUnavailable();
        }

        winnerIndex = RaffleMath.computeWinnerIndex(targetBlockhash, secret, roundId, tau);
    }

    /**
     * @notice Generates EIP-712 domain-separated ticket hash for a RoundTicket.
     */
    function roundTicketHash(RoundTicket calldata ticket) public view override returns (bytes32) {
        bytes32 structHash = keccak256(
            abi.encode(
                ROUND_TICKET_TYPEHASH,
                ticket.sender,
                ticket.recipient,
                ticket.roundId,
                ticket.localIndex,
                ticket.faceValue,
                ticket.winProb,
                ticket.senderNonce
            )
        );
        return _hashTypedDataV4(structHash);
    }

    /**
     * @notice Recovers the signer of an EIP-712 ticket hash using OpenZeppelin ECDSA library.
     */
    function recoverSigner(bytes32 ticketHash, bytes calldata sig) public pure override returns (address) {
        return ticketHash.recover(sig);
    }

    /**
     * @notice Verifies that keccak256(secret) matches the committed hash.
     */
    function verifySecretCommitment(bytes32 secret, bytes32 committedHash) public pure override returns (bool) {
        return keccak256(abi.encodePacked(secret)) == committedHash;
    }
}
