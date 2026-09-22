// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title ProtocolConfig
 * @notice Centralizes all configurable protocol parameters for the CIPHER Payment & Escrow Layer.
 * @dev Enforces domain rules and parameter boundary conditions.
 */
contract ProtocolConfig is Ownable {
    uint256 public constant CONFIRMATION_DELAY = 15;
    uint256 public constant MAX_BLOCKHASH_WINDOW = 256;
    uint256 public constant MIN_PROVIDER_STAKE = 1 ether;
    uint256 public constant UNBONDING_PERIOD = 7200; // ~1 day @ 12s blocks
    uint256 public constant CHANNEL_UNLOCK_PERIOD = 7200; // ~1 day
    uint256 public constant DISPUTE_RESPONSE_WINDOW = 1200; // ~4h @ 12s blocks
    uint256 public constant ROUND_STALE_TIMEOUT = 3600; // ~12h @ 12s blocks
    uint256 public constant MAX_ROUND_SIZE = 10_000;
    uint256 public constant SLASH_BPS = 2000; // 20% (out of 10,000)

    // Dynamic parameters overrideable by protocol owner if necessary
    uint256 public disputeWindow;
    uint256 public roundStaleTimeout;
    uint256 public slashBps;

    event ParametersUpdated(uint256 disputeWindow, uint256 roundStaleTimeout, uint256 slashBps);

    constructor(address initialOwner) Ownable(initialOwner) {
        disputeWindow = DISPUTE_RESPONSE_WINDOW;
        roundStaleTimeout = ROUND_STALE_TIMEOUT;
        slashBps = SLASH_BPS;
    }

    /**
     * @notice Updates tunable dispute window and stale timeouts.
     */
    function setParameters(
        uint256 _disputeWindow,
        uint256 _roundStaleTimeout,
        uint256 _slashBps
    ) external onlyOwner {
        require(_slashBps <= 10_000, "ProtocolConfig: Slash BPS exceeds 100%");
        require(_disputeWindow > 0, "ProtocolConfig: Invalid dispute window");
        require(_roundStaleTimeout > 0, "ProtocolConfig: Invalid stale timeout");

        disputeWindow = _disputeWindow;
        roundStaleTimeout = _roundStaleTimeout;
        slashBps = _slashBps;

        emit ParametersUpdated(_disputeWindow, _roundStaleTimeout, _slashBps);
    }
}
