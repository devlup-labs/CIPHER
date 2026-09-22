// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ProtocolConfig} from "../src/config/ProtocolConfig.sol";

contract ProtocolConfigTest is Test {
    ProtocolConfig public config;
    address public owner = address(0x1);
    address public alice = address(0x2);

    function setUp() public {
        vm.prank(owner);
        config = new ProtocolConfig(owner);
    }

    function test_InitialConstants() public view {
        assertEq(config.CONFIRMATION_DELAY(), 15);
        assertEq(config.MAX_BLOCKHASH_WINDOW(), 256);
        assertEq(config.MIN_PROVIDER_STAKE(), 1 ether);
        assertEq(config.UNBONDING_PERIOD(), 7200);
        assertEq(config.CHANNEL_UNLOCK_PERIOD(), 7200);
        assertEq(config.DISPUTE_RESPONSE_WINDOW(), 1200);
        assertEq(config.ROUND_STALE_TIMEOUT(), 3600);
        assertEq(config.MAX_ROUND_SIZE(), 10000);
        assertEq(config.SLASH_BPS(), 2000);
    }

    function test_InitialDynamicParameters() public view {
        assertEq(config.disputeWindow(), 1200);
        assertEq(config.roundStaleTimeout(), 3600);
        assertEq(config.slashBps(), 2000);
        assertEq(config.owner(), owner);
    }

    function test_SetParameters_Success() public {
        vm.prank(owner);
        config.setParameters(1800, 7200, 3000);

        assertEq(config.disputeWindow(), 1800);
        assertEq(config.roundStaleTimeout(), 7200);
        assertEq(config.slashBps(), 3000);
    }

    function test_SetParameters_RevertIfNotOwner() public {
        vm.prank(alice);
        vm.expectRevert();
        config.setParameters(1800, 7200, 3000);
    }

    function test_SetParameters_RevertInvalidBps() public {
        vm.prank(owner);
        vm.expectRevert("ProtocolConfig: Slash BPS exceeds 100%");
        config.setParameters(1200, 3600, 10001);
    }

    function test_SetParameters_RevertInvalidDisputeWindow() public {
        vm.prank(owner);
        vm.expectRevert("ProtocolConfig: Invalid dispute window");
        config.setParameters(0, 3600, 2000);
    }

    function test_SetParameters_RevertInvalidStaleTimeout() public {
        vm.prank(owner);
        vm.expectRevert("ProtocolConfig: Invalid stale timeout");
        config.setParameters(1200, 0, 2000);
    }
}
