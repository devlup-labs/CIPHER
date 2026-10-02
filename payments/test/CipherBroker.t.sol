// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import {CipherBroker} from "../src/Protocol_Raffle.sol";

contract CipherBrokerTest is Test {
    CipherBroker public broker;

    uint256 internal clientPk = 0xDEAF;
    address public client;
    address public provider = address(0x200);

    bytes32 secret = keccak256("BROKER_SECRET");
    bytes32 secretHash;

    receive() external payable {}

    function setUp() public {
        client = vm.addr(clientPk);
        secretHash = keccak256(abi.encodePacked(secret));

        broker = new CipherBroker();

        vm.deal(provider, 10 ether);
        vm.deal(client, 10 ether);
    }

    function test_Broker_ProviderRegistrationAndUnstake() public {
        vm.prank(provider);
        broker.registerProvider{value: 2 ether}();

        (uint256 stake, uint256 releaseBlock, bool registered) = broker.providers(provider);
        assertEq(stake, 2 ether);
        assertTrue(registered);
        assertEq(releaseBlock, 0);

        vm.prank(provider);
        broker.requestUnstake();

        (, releaseBlock, ) = broker.providers(provider);
        assertEq(releaseBlock, block.number + broker.UNBONDING_PERIOD());
    }

    function test_Broker_ChannelLifecycle() public {
        vm.prank(client);
        broker.openChannel{value: 5 ether}(provider);

        assertEq(broker.channelBalance(client, provider), 5 ether);

        vm.prank(client);
        broker.initiateChannelUnlock(provider);

        bytes32 chKey = keccak256(abi.encodePacked(client, provider));
        (,,, uint256 unlockBlock) = broker.channels(chKey);
        assertEq(unlockBlock, block.number + broker.CHANNEL_UNLOCK_PERIOD());
    }

    function test_Broker_RoundCommitAndSettlement() public {
        vm.prank(provider);
        broker.registerProvider{value: 2 ether}();

        vm.prank(client);
        broker.openChannel{value: 5 ether}(provider);

        uint256 roundId = 1;
        uint256 tau = 10;
        uint256 faceValue = 1 ether;

        vm.prank(provider);
        broker.commitRoundSecret(client, roundId, tau, faceValue, secretHash);

        uint256 commitBlock = block.number;
        uint256 targetBlock = commitBlock + broker.CONFIRMATION_DELAY();
        vm.roll(targetBlock + 1);

        bytes32 bh = blockhash(targetBlock);
        uint256 winnerIdx = uint256(keccak256(abi.encodePacked(bh, secret, roundId))) % tau;

        CipherBroker.RoundTicket memory ticket = CipherBroker.RoundTicket({
            sender: client,
            recipient: provider,
            roundId: roundId,
            localIndex: winnerIdx,
            faceValue: faceValue,
            winProb: type(uint256).max,
            senderNonce: 1
        });

        bytes32 rawHash = broker.roundTicketHashFunc(ticket);
        bytes32 ethSignedHash = MessageHashUtils.toEthSignedMessageHash(rawHash);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(clientPk, ethSignedHash);
        bytes memory sig = abi.encodePacked(r, s, v);

        uint256 providerBalBefore = provider.balance;
        vm.prank(provider);
        broker.settleRound(ticket, sig, secret);

        assertEq(provider.balance, providerBalBefore + faceValue);
    }
}
