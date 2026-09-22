// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ProtocolConfig} from "../src/config/ProtocolConfig.sol";
import {PaymentChannel} from "../src/core/PaymentChannel.sol";
import {IPaymentChannel} from "../src/interfaces/IPaymentChannel.sol";

contract PaymentChannelTest is Test {
    ProtocolConfig public config;
    PaymentChannel public channel;
    address public owner = address(this);
    address public client = address(0x100);
    address public provider = address(0x200);
    address public payoutEngine = address(0x300);

    receive() external payable {}

    function setUp() public {
        config = new ProtocolConfig(owner);
        channel = new PaymentChannel(address(config));
        vm.deal(client, 10 ether);
    }

    function test_OpenChannel_Success() public {
        vm.prank(client);
        channel.openChannel{value: 5 ether}(provider);

        assertEq(channel.channelBalance(client, provider), 5 ether);
    }

    function test_OpenChannel_RevertZeroDeposit() public {
        vm.prank(client);
        vm.expectRevert(IPaymentChannel.ZeroDeposit.selector);
        channel.openChannel{value: 0}(provider);
    }

    function test_DepositToChannel_Success() public {
        vm.startPrank(client);
        channel.openChannel{value: 2 ether}(provider);
        channel.depositToChannel{value: 3 ether}(provider);
        vm.stopPrank();

        assertEq(channel.channelBalance(client, provider), 5 ether);
    }

    function test_InitiateChannelUnlock_Success() public {
        vm.startPrank(client);
        channel.openChannel{value: 5 ether}(provider);
        channel.initiateChannelUnlock(provider);
        vm.stopPrank();

        assertTrue(channel.isUnlockInitiated(client, provider));
    }

    function test_InitiateChannelUnlock_RevertAlreadyInitiated() public {
        vm.startPrank(client);
        channel.openChannel{value: 5 ether}(provider);
        channel.initiateChannelUnlock(provider);
        vm.expectRevert(IPaymentChannel.UnlockAlreadyInitiated.selector);
        channel.initiateChannelUnlock(provider);
        vm.stopPrank();
    }

    function test_WithdrawChannel_Success() public {
        vm.startPrank(client);
        channel.openChannel{value: 5 ether}(provider);
        channel.initiateChannelUnlock(provider);
        vm.stopPrank();

        uint256 unlockPeriod = config.CHANNEL_UNLOCK_PERIOD();
        vm.roll(block.number + unlockPeriod + 1);

        uint256 balBefore = client.balance;
        vm.prank(client);
        channel.withdrawChannel(provider);

        assertEq(client.balance, balBefore + 5 ether);
        assertEq(channel.channelBalance(client, provider), 0);
    }

    function test_PayoutProvider_Success() public {
        vm.prank(client);
        channel.openChannel{value: 5 ether}(provider);

        channel.setPayoutEngine(payoutEngine, true);

        uint256 providerBalBefore = provider.balance;
        vm.prank(payoutEngine);
        channel.payoutProvider(client, provider, 2 ether);

        assertEq(provider.balance, providerBalBefore + 2 ether);
        assertEq(channel.channelBalance(client, provider), 3 ether);
    }

    function test_PayoutProvider_RevertUnauthorized() public {
        vm.prank(client);
        channel.openChannel{value: 5 ether}(provider);

        vm.prank(client);
        vm.expectRevert(IPaymentChannel.UnauthorizedCaller.selector);
        channel.payoutProvider(client, provider, 2 ether);
    }

    function test_PayoutProvider_RevertExceedsBalance() public {
        vm.prank(client);
        channel.openChannel{value: 1 ether}(provider);

        channel.setPayoutEngine(payoutEngine, true);

        vm.prank(payoutEngine);
        vm.expectRevert(IPaymentChannel.InsufficientChannelBalance.selector);
        channel.payoutProvider(client, provider, 2 ether);
    }

    function test_PauseUnpause() public {
        channel.pause();
        vm.prank(client);
        vm.expectRevert();
        channel.openChannel{value: 1 ether}(provider);

        channel.unpause();
        vm.prank(client);
        channel.openChannel{value: 1 ether}(provider);
        assertEq(channel.channelBalance(client, provider), 1 ether);
    }
}
