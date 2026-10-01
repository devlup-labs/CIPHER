// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title IProviderRegistry
 * @notice Interface for provider registration, stake management, partial unbonding, and slashing.
 */
interface IProviderRegistry {
    struct Provider {
        uint256 stake;
        uint256 unstakeReleaseBlock;
        uint256 pendingUnstakeAmount;
        bool registered;
    }

    event ProviderRegistered(address indexed provider, uint256 totalStake);
    event ProviderUnstakeRequested(address indexed provider, uint256 releaseBlock, uint256 amount);
    event ProviderStakeWithdrawn(address indexed provider, uint256 amount);
    event ProviderSlashed(address indexed provider, uint256 amount, bytes32 indexed reasonId);

    error InsufficientStake();
    error NotRegistered();
    error UnstakeAlreadyRequested();
    error UnstakeNotRequested();
    error UnbondingNotComplete();
    error TransferFailed();
    error UnauthorizedCaller();
    error InvalidAmount();

    function registerProvider() external payable;
    function requestUnstake() external;
    function requestPartialUnstake(uint256 amount) external;
    function withdrawStake() external;
    function slashProvider(address provider, uint256 amount, bytes32 reasonId) external returns (uint256 slashedAmount);
    function isProviderActive(address provider) external view returns (bool);
    function getProviderStake(address provider) external view returns (uint256);
}
