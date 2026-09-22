// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ProtocolConfig} from "../src/config/ProtocolConfig.sol";
import {ProviderRegistry} from "../src/core/ProviderRegistry.sol";
import {PaymentChannel} from "../src/core/PaymentChannel.sol";
import {CommitRevealEntropy} from "../src/randomness/CommitRevealEntropy.sol";
import {SettlementEngine} from "../src/core/SettlementEngine.sol";
import {ISettlementEngine} from "../src/interfaces/ISettlementEngine.sol";
import {IEntropySource} from "../src/interfaces/IEntropySource.sol";
import {RaffleMath} from "../src/randomness/RaffleMath.sol";

contract SettlementEngineTest is Test {
    ProtocolConfig public config;
    ProviderRegistry public registry;
    PaymentChannel public channel;
    CommitRevealEntropy public entropy;
    SettlementEngine public engine;

    address public owner = address(this);
    uint256 internal clientPk = 0xC11E7;
    address public client;
    address public provider = address(0x200);

    bytes32 secret = keccak256("PROVIDER_SECRET");
    bytes32 secretHash;

    receive() external payable {}

    function setUp() public {
        client = vm.addr(clientPk);
        secretHash = keccak256(abi.encodePacked(secret));

        config = new ProtocolConfig(owner);
        registry = new ProviderRegistry(address(config));
        channel = new PaymentChannel(address(config));
        entropy = new CommitRevealEntropy(address(config));

        engine = new SettlementEngine(
            address(config),
            address(entropy),
            address(channel),
            address(registry)
        );

        channel.setPayoutEngine(address(engine), true);

        vm.deal(provider, 10 ether);
        vm.deal(client, 10 ether);

        vm.prank(provider);
        registry.registerProvider{value: 2 ether}();

        vm.prank(client);
        channel.openChannel{value: 5 ether}(provider);
    }

    function test_CommitRoundSecret_Success() public {
        vm.prank(provider);
        engine.commitRoundSecret(client, 1, 100, 1 ether, secretHash);

        assertEq(
            uint256(engine.getRoundStatus(client, provider, 1)),
            uint256(ISettlementEngine.RoundStatus.Open)
        );
    }

    function test_CommitRoundSecret_RevertInvalidParams() public {
        vm.prank(provider);
        vm.expectRevert(ISettlementEngine.InvalidRoundParams.selector);
        engine.commitRoundSecret(client, 1, 0, 1 ether, secretHash);

        vm.prank(provider);
        vm.expectRevert(ISettlementEngine.InvalidRoundParams.selector);
        engine.commitRoundSecret(client, 1, 100, 0, secretHash);
    }

    function test_CommitRoundSecret_RevertAlreadyCommitted() public {
        vm.startPrank(provider);
        engine.commitRoundSecret(client, 1, 100, 1 ether, secretHash);
        vm.expectRevert(ISettlementEngine.RoundAlreadyCommitted.selector);
        engine.commitRoundSecret(client, 1, 100, 1 ether, secretHash);
        vm.stopPrank();
    }

    function test_SettleRound_Success() public {
        uint256 roundId = 1;
        uint256 tau = 10;
        uint256 faceValue = 1 ether;

        vm.prank(provider);
        engine.commitRoundSecret(client, roundId, tau, faceValue, secretHash);

        uint256 commitBlock = block.number;
        uint256 targetBlock = commitBlock + config.CONFIRMATION_DELAY();
        vm.roll(targetBlock + 1);

        bytes32 bh = blockhash(targetBlock);
        uint256 winnerIdx = RaffleMath.computeWinnerIndex(bh, secret, roundId, tau);

        IEntropySource.RoundTicket memory ticket = IEntropySource.RoundTicket({
            sender: client,
            recipient: provider,
            roundId: roundId,
            localIndex: winnerIdx,
            faceValue: faceValue,
            winProb: type(uint256).max,
            senderNonce: 1
        });

        bytes32 tHash = entropy.roundTicketHash(ticket);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(clientPk, tHash);
        bytes memory sig = abi.encodePacked(r, s, v);

        uint256 providerBalBefore = provider.balance;
        vm.prank(provider);
        engine.settleRound(ticket, sig, secret);

        assertEq(provider.balance, providerBalBefore + faceValue);
        assertEq(
            uint256(engine.getRoundStatus(client, provider, roundId)),
            uint256(ISettlementEngine.RoundStatus.Settled)
        );
    }

    function test_MarkRoundPartial_Success() public {
        uint256 roundId = 1;
        vm.prank(provider);
        engine.commitRoundSecret(client, roundId, 10, 1 ether, secretHash);

        // Roll past ROUND_STALE_TIMEOUT (3600)
        vm.roll(block.number + 3601);

        vm.prank(provider);
        engine.markRoundPartial(client, roundId);

        assertEq(
            uint256(engine.getRoundStatus(client, provider, roundId)),
            uint256(ISettlementEngine.RoundStatus.PartialFallback)
        );
    }
}
