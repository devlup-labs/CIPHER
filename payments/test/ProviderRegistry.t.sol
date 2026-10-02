// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ProtocolConfig} from "../src/config/ProtocolConfig.sol";
import {ProviderRegistry} from "../src/core/ProviderRegistry.sol";
import {IProviderRegistry} from "../src/interfaces/IProviderRegistry.sol";

contract ProviderRegistryTest is Test {
    ProtocolConfig public config;
    ProviderRegistry public registry;
    address public owner = address(this);
    address public provider1 = makeAddr("provider1");
    address public slasher = makeAddr("slasher");

    receive() external payable {}

    function setUp() public {
        config = new ProtocolConfig(owner);
        registry = new ProviderRegistry(address(config));
        vm.deal(provider1, 10 ether);
    }

    function test_RegisterProvider_Success() public {
        vm.prank(provider1);
        registry.registerProvider{value: 2 ether}();

        assertEq(registry.getProviderStake(provider1), 2 ether);
        assertTrue(registry.isProviderActive(provider1));
    }

    function test_RegisterProvider_RevertInsufficientStake() public {
        vm.prank(provider1);
        vm.expectRevert(IProviderRegistry.InsufficientStake.selector);
        registry.registerProvider{value: 0.5 ether}();
    }

    function test_TopUpProviderStake() public {
        vm.startPrank(provider1);
        registry.registerProvider{value: 1 ether}();
        registry.registerProvider{value: 1 ether}();
        vm.stopPrank();

        assertEq(registry.getProviderStake(provider1), 2 ether);
    }

    function test_RequestUnstake_Success() public {
        vm.startPrank(provider1);
        registry.registerProvider{value: 2 ether}();
        registry.requestUnstake();
        vm.stopPrank();

        assertFalse(registry.isProviderActive(provider1));
    }

    function test_RequestUnstake_RevertIfNotRegistered() public {
        vm.prank(provider1);
        vm.expectRevert(IProviderRegistry.NotRegistered.selector);
        registry.requestUnstake();
    }

    function test_RequestUnstake_RevertAlreadyRequested() public {
        vm.startPrank(provider1);
        registry.registerProvider{value: 2 ether}();
        registry.requestUnstake();
        vm.expectRevert(IProviderRegistry.UnstakeAlreadyRequested.selector);
        registry.requestUnstake();
        vm.stopPrank();
    }

    function test_RequestPartialUnstake_Success() public {
        vm.startPrank(provider1);
        registry.registerProvider{value: 3 ether}();
        registry.requestPartialUnstake(1 ether);
        vm.stopPrank();

        assertFalse(registry.isProviderActive(provider1));
    }

    function test_RequestPartialUnstake_RevertRemainingTooLow() public {
        vm.startPrank(provider1);
        registry.registerProvider{value: 1.5 ether}();
        vm.expectRevert(IProviderRegistry.InsufficientStake.selector);
        registry.requestPartialUnstake(1 ether);
        vm.stopPrank();
    }

    function test_WithdrawStake_Success() public {
        vm.startPrank(provider1);
        registry.registerProvider{value: 2 ether}();
        registry.requestUnstake();
        vm.stopPrank();

        uint256 unbondBlocks = config.UNBONDING_PERIOD();
        vm.roll(block.number + unbondBlocks + 1);

        uint256 balBefore = provider1.balance;
        vm.prank(provider1);
        registry.withdrawStake();

        assertEq(provider1.balance, balBefore + 2 ether);
        assertEq(registry.getProviderStake(provider1), 0);
    }

    function test_WithdrawStake_RevertEarlyWithdraw() public {
        vm.startPrank(provider1);
        registry.registerProvider{value: 2 ether}();
        registry.requestUnstake();

        vm.expectRevert(IProviderRegistry.UnbondingNotComplete.selector);
        registry.withdrawStake();
        vm.stopPrank();
    }

    function test_SlashProvider_Success() public {
        vm.prank(provider1);
        registry.registerProvider{value: 2 ether}();

        registry.setSlasher(slasher, true);

        vm.prank(slasher);
        uint256 slashed = registry.slashProvider(provider1, 0.5 ether, keccak256("REASON"));

        assertEq(slashed, 0.5 ether);
        assertEq(registry.getProviderStake(provider1), 1.5 ether);
    }

    function test_SlashProvider_RevertUnauthorized() public {
        vm.prank(provider1);
        registry.registerProvider{value: 2 ether}();

        vm.prank(provider1);
        vm.expectRevert(IProviderRegistry.UnauthorizedCaller.selector);
        registry.slashProvider(provider1, 0.5 ether, keccak256("REASON"));
    }
}
