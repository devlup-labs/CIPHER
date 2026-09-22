// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ProtocolConfig} from "../src/config/ProtocolConfig.sol";
import {CommitRevealEntropy} from "../src/randomness/CommitRevealEntropy.sol";
import {IEntropySource} from "../src/interfaces/IEntropySource.sol";

contract CommitRevealEntropyTest is Test {
    ProtocolConfig public config;
    CommitRevealEntropy public entropy;
    address public owner = address(this);

    uint256 internal clientPrivateKey = 0xA11CE;
    address internal client;

    function setUp() public {
        config = new ProtocolConfig(owner);
        entropy = new CommitRevealEntropy(address(config));
        client = vm.addr(clientPrivateKey);
    }

    function test_VerifySecretCommitment() public view {
        bytes32 secret = keccak256("SECRET_KEY");
        bytes32 secretHash = keccak256(abi.encodePacked(secret));

        assertTrue(entropy.verifySecretCommitment(secret, secretHash));
        assertFalse(entropy.verifySecretCommitment(keccak256("WRONG"), secretHash));
    }

    function test_TicketHashAndSignerRecovery() public view {
        IEntropySource.RoundTicket memory ticket = IEntropySource.RoundTicket({
            sender: client,
            recipient: address(0x200),
            roundId: 1,
            localIndex: 0,
            faceValue: 1 ether,
            winProb: type(uint256).max,
            senderNonce: 42
        });

        bytes32 tHash = entropy.roundTicketHash(ticket);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(clientPrivateKey, tHash);
        bytes memory signature = abi.encodePacked(r, s, v);

        address recovered = entropy.recoverSigner(tHash, signature);
        assertEq(recovered, client);
    }

    function test_ComputeWinnerIndex_RevertNotYetDrawable() public {
        bytes32 secret = keccak256("MY_SECRET");
        uint256 commitBlock = 100;
        vm.roll(105);

        // Target block is 115 (100 + 15). At block 105, it's not yet drawable.
        vm.expectRevert(IEntropySource.RoundNotYetDrawable.selector);
        entropy.computeWinnerIndex(commitBlock, secret, 1, 100);
    }

    function test_ComputeWinnerIndex_RevertDrawExpired() public {
        bytes32 secret = keccak256("MY_SECRET");
        uint256 commitBlock = 100;
        vm.roll(400);

        // Max window is 256 (100 + 256 = 356). At block 400, draw is expired.
        vm.expectRevert(IEntropySource.RoundDrawExpired.selector);
        entropy.computeWinnerIndex(commitBlock, secret, 1, 100);
    }
}
