// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ProtocolConfig} from "../src/config/ProtocolConfig.sol";
import {ProviderRegistry} from "../src/core/ProviderRegistry.sol";
import {PaymentChannel} from "../src/core/PaymentChannel.sol";
import {CommitRevealEntropy} from "../src/randomness/CommitRevealEntropy.sol";
import {SettlementEngine} from "../src/core/SettlementEngine.sol";
import {ChunkDisputeResolver} from "../src/core/ChunkDisputeResolver.sol";
import {IDisputeResolver} from "../src/interfaces/IDisputeResolver.sol";
import {ISettlementEngine} from "../src/interfaces/ISettlementEngine.sol";

contract ChunkDisputeResolverTest is Test {
    ProtocolConfig public config;
    ProviderRegistry public registry;
    PaymentChannel public channel;
    CommitRevealEntropy public entropy;
    SettlementEngine public engine;
    ChunkDisputeResolver public resolver;

    address public owner = address(this);
    address public client = address(0x100);
    address public provider = address(0x200);

    bytes32 secretHash = keccak256("SECRET");

    function setUp() public {
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
        resolver = new ChunkDisputeResolver(
            address(config),
            address(registry),
            address(engine),
            address(channel)
        );

        engine.setDisputeResolver(address(resolver));
        registry.setSlasher(address(resolver), true);

        vm.deal(provider, 10 ether);
        vm.deal(client, 10 ether);

        vm.prank(provider);
        registry.registerProvider{value: 2 ether}();

        vm.prank(client);
        channel.openChannel{value: 5 ether}(provider);
    }

    function test_RaiseChunkDispute_Success() public {
        vm.prank(provider);
        engine.commitRoundSecret(client, 1, 100, 1 ether, secretHash);

        vm.prank(client);
        resolver.raiseChunkDispute(provider, 1, 0, keccak256("EVIDENCE"));

        assertEq(
            uint256(resolver.getDisputeStatus(client, provider, 1, 0)),
            uint256(IDisputeResolver.DisputeStatus.Open)
        );
        assertTrue(resolver.isDisputePending(client, provider, 1));
    }

    function test_RaiseChunkDispute_RevertAlreadyOpen() public {
        vm.prank(provider);
        engine.commitRoundSecret(client, 1, 100, 1 ether, secretHash);

        vm.startPrank(client);
        resolver.raiseChunkDispute(provider, 1, 0, keccak256("EVIDENCE"));
        vm.expectRevert(IDisputeResolver.DisputeAlreadyOpen.selector);
        resolver.raiseChunkDispute(provider, 1, 0, keccak256("EVIDENCE"));
        vm.stopPrank();
    }

    function test_ResolveExpiredDispute_Success() public {
        vm.prank(provider);
        engine.commitRoundSecret(client, 1, 100, 1 ether, secretHash);

        vm.prank(client);
        resolver.raiseChunkDispute(provider, 1, 0, keccak256("EVIDENCE"));

        // Roll past DISPUTE_RESPONSE_WINDOW (1200)
        vm.roll(block.number + 1201);

        uint256 stakeBefore = registry.getProviderStake(provider);
        resolver.resolveExpiredDispute(client, provider, 1, 0);

        uint256 expectedSlash = (stakeBefore * config.SLASH_BPS()) / 10_000;
        assertEq(registry.getProviderStake(provider), stakeBefore - expectedSlash);

        assertEq(
            uint256(resolver.getDisputeStatus(client, provider, 1, 0)),
            uint256(IDisputeResolver.DisputeStatus.ResolvedSlashed)
        );
        assertEq(
            uint256(engine.getRoundStatus(client, provider, 1)),
            uint256(ISettlementEngine.RoundStatus.Voided)
        );
    }
}
