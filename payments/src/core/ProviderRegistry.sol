// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {ProtocolConfig} from "../config/ProtocolConfig.sol";

/**
 * @title ProviderRegistry
 * @notice Manages provider registration, collateral staking, flexible partial unbonding timers, and authorized slashing.
 * @dev Employs strict Checks-Effects-Interactions (CEI) to eliminate reentrancy risks.
 */
contract ProviderRegistry is IProviderRegistry, ReentrancyGuard {
    ProtocolConfig public immutable config;

    mapping(address => Provider) public providers;
    mapping(address => bool) public authorizedSlashers;
    address public owner;

    modifier onlyOwner() {
        if (msg.sender != owner) revert UnauthorizedCaller();
        _;
    }

    modifier onlySlasher() {
        if (!authorizedSlashers[msg.sender]) revert UnauthorizedCaller();
        _;
    }

    constructor(address _config) {
        require(_config != address(0), "ProviderRegistry: Zero config address");
        config = ProtocolConfig(_config);
        owner = msg.sender;
    }

    function setSlasher(address slasher, bool authorized) external onlyOwner {
        authorizedSlashers[slasher] = authorized;
    }

    /**
     * @notice Registers provider or tops up existing provider stake.
     */
    function registerProvider() external payable override {
        Provider storage p = providers[msg.sender];
        if (p.stake + msg.value < config.MIN_PROVIDER_STAKE()) {
            revert InsufficientStake();
        }

        p.stake += msg.value;
        p.registered = true;

        emit ProviderRegistered(msg.sender, p.stake);
    }

    /**
     * @notice Signals intent to withdraw full provider stake and starts unbonding timer.
     */
    function requestUnstake() external override {
        Provider storage p = providers[msg.sender];
        if (!p.registered) revert NotRegistered();
        if (p.unstakeReleaseBlock != 0) revert UnstakeAlreadyRequested();

        p.pendingUnstakeAmount = p.stake;
        p.unstakeReleaseBlock = block.number + config.UNBONDING_PERIOD();
        emit ProviderUnstakeRequested(msg.sender, p.unstakeReleaseBlock, p.pendingUnstakeAmount);
    }

    /**
     * @notice Signals intent to withdraw a partial amount of provider stake without deactivating node.
     */
    function requestPartialUnstake(uint256 amount) external override {
        Provider storage p = providers[msg.sender];
        if (!p.registered) revert NotRegistered();
        if (p.unstakeReleaseBlock != 0) revert UnstakeAlreadyRequested();
        if (amount == 0 || p.stake < amount) revert InvalidAmount();
        if (p.stake - amount < config.MIN_PROVIDER_STAKE()) revert InsufficientStake();

        p.pendingUnstakeAmount = amount;
        p.unstakeReleaseBlock = block.number + config.UNBONDING_PERIOD();
        emit ProviderUnstakeRequested(msg.sender, p.unstakeReleaseBlock, amount);
    }

    /**
     * @notice Reclaims pending unbonded provider stake after unbonding period completes.
     * @dev Strictly applies CEI prior to ETH transfer.
     */
    function withdrawStake() external override nonReentrant {
        Provider storage p = providers[msg.sender];
        if (p.unstakeReleaseBlock == 0) revert UnstakeNotRequested();
        if (p.unstakeReleaseBlock > block.number) revert UnbondingNotComplete();

        uint256 amount = p.pendingUnstakeAmount;
        if (amount == 0 || p.stake < amount) revert InsufficientStake();

        // Effects
        p.stake -= amount;
        p.pendingUnstakeAmount = 0;
        p.unstakeReleaseBlock = 0;

        if (p.stake < config.MIN_PROVIDER_STAKE()) {
            p.registered = false;
        }

        // Interactions
        (bool callSuccess, ) = payable(msg.sender).call{value: amount}("");
        if (!callSuccess) revert TransferFailed();

        emit ProviderStakeWithdrawn(msg.sender, amount);
    }

    /**
     * @notice Deducts penalty from provider's stake upon verified dispute timeout or failure.
     */
    function slashProvider(
        address provider,
        uint256 amount,
        bytes32 reasonId
    ) external override onlySlasher returns (uint256 slashedAmount) {
        Provider storage p = providers[provider];
        if (p.stake == 0) return 0;

        slashedAmount = amount > p.stake ? p.stake : amount;
        p.stake -= slashedAmount;

        if (p.stake < config.MIN_PROVIDER_STAKE()) {
            p.registered = false;
        }

        emit ProviderSlashed(provider, slashedAmount, reasonId);
    }

    function isProviderActive(address provider) external view override returns (bool) {
        Provider storage p = providers[provider];
        return p.registered && p.stake >= config.MIN_PROVIDER_STAKE() && p.unstakeReleaseBlock == 0;
    }

    function getProviderStake(address provider) external view override returns (uint256) {
        return providers[provider].stake;
    }
}
