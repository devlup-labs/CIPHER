// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
contract ProtocolConfig is Ownable {
   uint256 public constant CONFIRMATION_DELAY = 15;
   uint256 public constant MAX_BLOCKHASH_WINDOW = 256;
   uint256 public constant MIN_PROVIDER_STAKE = 1 ether;
   uint256 public constant UNBONDING_PERIOD = 7200; // ~1 day
   uint256 public constant CHANNEL_UNLOCK_PERIOD = 7200; // ~1 day
   uint256 public constant DISPUTE_RESPONSE_WINDOW = 1200; // ~4h
   uint256 public constant ROUND_STALE_TIMEOUT = 3600; // ~12h
   uint256 public constant MAX_ROUND_SIZE = 10_000;
   uint256 public constant SLASH_BPS = 2000; // 20%
   uint256 public disputeWindow;
   uint256 public roundStaleTimeout;
   uint256 public slashBps;
   constructor(address initialOwner) Ownable(initialOwner) {
       disputeWindow = DISPUTE_RESPONSE_WINDOW;
       roundStaleTimeout = ROUND_STALE_TIMEOUT;
       slashBps = SLASH_BPS;
   }
   function setParameters(uint256 _disputeWindow, uint256 _roundStaleTimeout, uint256 _slashBps) external onlyOwner {
       require(_slashBps <= 10_000, "Invalid BPS");
       disputeWindow = _disputeWindow;
       roundStaleTimeout = _roundStaleTimeout;
       slashBps = _slashBps;
   }
}
