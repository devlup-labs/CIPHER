// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";
import {IPaymentChannel} from "../interfaces/IPaymentChannel.sol";
import {ProtocolConfig} from "../config/ProtocolConfig.sol";

/**
 * @title PaymentChannel
 * @notice Isolated escrow channel contract managing fund locking, deposits, unlock timeouts, and provider payouts.
 * @dev Includes OpenZeppelin Pausable circuit breaker for emergency protocol protection.
 */
contract PaymentChannel is IPaymentChannel, ReentrancyGuard, Pausable {
    ProtocolConfig public immutable config;

    // key = keccak256(abi.encodePacked(sender, recipient))
    mapping(bytes32 => Channel) public channels;

    // Authorized callers allowed to trigger payouts (e.g., SettlementEngine)
    mapping(address => bool) public authorizedPayoutEngine;

    address public owner;

    modifier onlyOwner() {
        if (msg.sender != owner) revert UnauthorizedCaller();
        _;
    }

    modifier onlyPayoutEngine() {
        if (!authorizedPayoutEngine[msg.sender]) revert UnauthorizedCaller();
        _;
    }

    constructor(address _config) {
        require(_config != address(0), "PaymentChannel: Zero config address");
        config = ProtocolConfig(_config);
        owner = msg.sender;
    }

    function setPayoutEngine(address engine, bool authorized) external onlyOwner {
        authorizedPayoutEngine[engine] = authorized;
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    function getChannelKey(address sender, address recipient) public pure override returns (bytes32) {
        return keccak256(abi.encodePacked(sender, recipient));
    }

    /**
     * @notice Client opens a payment escrow channel for a provider.
     */
    function openChannel(address recipient) external payable override whenNotPaused {
        if (msg.value == 0) revert ZeroDeposit();

        bytes32 key = getChannelKey(msg.sender, recipient);
        Channel storage ch = channels[key];
        ch.totalDeposited += msg.value;

        emit ChannelOpened(msg.sender, recipient, msg.value);
    }

    /**
     * @notice Top up an existing payment channel.
     */
    function depositToChannel(address recipient) external payable override whenNotPaused {
        if (msg.value == 0) revert ZeroDeposit();

        bytes32 key = getChannelKey(msg.sender, recipient);
        Channel storage ch = channels[key];
        ch.totalDeposited += msg.value;

        emit ChannelDeposited(msg.sender, recipient, msg.value);
    }

    /**
     * @notice Client initiates unlock/closure of payment channel.
     */
    function initiateChannelUnlock(address recipient) external override whenNotPaused {
        bytes32 key = getChannelKey(msg.sender, recipient);
        Channel storage ch = channels[key];

        if (ch.unlockBlock != 0) revert UnlockAlreadyInitiated();

        ch.unlockBlock = block.number + config.CHANNEL_UNLOCK_PERIOD();
        emit ChannelUnlockInitiated(msg.sender, recipient, ch.unlockBlock);
    }

    /**
     * @notice Client reclaims unspent funds after unlock window expires.
     */
    function withdrawChannel(address recipient) external override nonReentrant whenNotPaused {
        bytes32 key = getChannelKey(msg.sender, recipient);
        Channel storage ch = channels[key];

        if (ch.unlockBlock == 0) revert UnlockNotInitiated();
        if (ch.unlockBlock > block.number) revert UnlockNotComplete();

        uint256 remaining = ch.totalDeposited - ch.settledAmount - ch.withdrawnAmount;
        if (remaining == 0) revert NothingToWithdraw();

        ch.withdrawnAmount += remaining;

        (bool callSuccess, ) = payable(msg.sender).call{value: remaining}("");
        if (!callSuccess) revert TransferFailed();

        emit ChannelWithdrawn(msg.sender, recipient, remaining);
    }

    /**
     * @notice Transfers earned payout to provider from channel escrow upon winning ticket settlement.
     */
    function payoutProvider(
        address sender,
        address recipient,
        uint256 amount
    ) external override onlyPayoutEngine nonReentrant whenNotPaused {
        bytes32 key = getChannelKey(sender, recipient);
        Channel storage ch = channels[key];

        uint256 available = ch.totalDeposited - ch.settledAmount - ch.withdrawnAmount;
        if (amount > available) revert InsufficientChannelBalance();

        ch.settledAmount += amount;

        (bool callSuccess, ) = payable(recipient).call{value: amount}("");
        if (!callSuccess) revert TransferFailed();

        emit ChannelPayoutExecuted(sender, recipient, amount);
    }

    function channelBalance(address sender, address recipient) external view override returns (uint256) {
        bytes32 key = getChannelKey(sender, recipient);
        Channel storage ch = channels[key];
        return ch.totalDeposited - ch.settledAmount - ch.withdrawnAmount;
    }

    function isUnlockInitiated(address sender, address recipient) external view override returns (bool) {
        bytes32 key = getChannelKey(sender, recipient);
        return channels[key].unlockBlock != 0;
    }
}
