// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title IPaymentChannel
 * @notice Interface for isolated payment escrow channels between clients and providers.
 */
interface IPaymentChannel {
    struct Channel {
        uint256 totalDeposited;
        uint256 settledAmount;
        uint256 withdrawnAmount;
        uint256 unlockBlock;
    }

    event ChannelOpened(address indexed sender, address indexed recipient, uint256 amount);
    event ChannelDeposited(address indexed sender, address indexed recipient, uint256 amount);
    event ChannelUnlockInitiated(address indexed sender, address indexed recipient, uint256 unlockBlock);
    event ChannelWithdrawn(address indexed sender, address indexed recipient, uint256 amount);
    event ChannelPayoutExecuted(address indexed sender, address indexed recipient, uint256 amount);

    error ZeroDeposit();
    error UnlockAlreadyInitiated();
    error UnlockNotInitiated();
    error UnlockNotComplete();
    error NothingToWithdraw();
    error InsufficientChannelBalance();
    error TransferFailed();
    error UnauthorizedCaller();

    function openChannel(address recipient) external payable;
    function depositToChannel(address recipient) external payable;
    function initiateChannelUnlock(address recipient) external;
    function withdrawChannel(address recipient) external;
    function payoutProvider(address sender, address recipient, uint256 amount) external;
    function channelBalance(address sender, address recipient) external view returns (uint256);
    function isUnlockInitiated(address sender, address recipient) external view returns (bool);
    function getChannelKey(address sender, address recipient) external pure returns (bytes32);
}
