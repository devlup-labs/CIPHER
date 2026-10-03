# CIPHER General Testing Architecture & Unified Verification Strategy

This document outlines the testing methodology, system roles, cryptographic verification layers, prerequisite toolchains, and repository-relative execution conventions for the CIPHER decentralized content delivery network (CDN) and probabilistic micropayment architecture.

> [!IMPORTANT]
> **Notice for Local Testing:**
> All local multi-terminal testing, role simulations, and verification suites must be executed against the **`local`** branch of the repository:
> 👉 **[https://github.com/devlup-labs/CIPHER/tree/local](https://github.com/devlup-labs/CIPHER/tree/local)**
>
> ```bash
> git clone https://github.com/devlup-labs/CIPHER.git
> cd CIPHER
> git checkout local
> ```

---

## 1. Operating System Platform Index

Detailed setup guides, external resource installations, permission requirements, and troubleshooting for each specific operating system are documented in dedicated guides:

* 🪟 **[Windows Testing Guide (WSL2 / Git Bash)](windows.md)**: Environment setup, WSL2 installation, Foundry toolchain, firewall permissions, and execution commands.
* 🐧 **[Linux Testing Guide (Ubuntu, Debian, Fedora, Arch)](linux.md)**: Package manager dependencies, non-root port allocation, ulimit tuning, and headless CI automation.
* 🍏 **[macOS Testing Guide (Apple Silicon & Intel)](mac.md)**: Homebrew packages, Terminal.app AppleScript desktop perimeter tiling, Automation permissions, and single-terminal fallback.

---

## 2. Core Architecture & Verification Matrix

CIPHER operates as a trust-minimized, decentralized CDN with cryptographic content authentication, distributed Kademlia DHT routing, libp2p Circuit Relay v2 NAT traversal, EIP-712 probabilistic micropayment ticketing, and on-chain EVM settlement.

```text
+---------------------------------------------------------------------------------------------------+
|                                      CIPHER TESTBED TOPOLOGY                                      |
+---------------------------------------------------------------------------------------------------+
|                                                                                                   |
|  [ LAYER 1 EVM BLOCKCHAIN (Anvil :8545) ] <-----------------------------------+                   |
|    - CommitRevealEntropy.sol (Raffle randomness)                              | (On-Chain Dispute)|
|    - PaymentChannel.sol (Escrow & Staking)                                    |                   |
|    - SettlementEngine.sol (Winning ticket payouts)                            |                   |
|    - ProviderRegistry.sol (Node registration)                                 |                   |
|                                                                               |                   |
|  [ CONTROL PLANE: Kademlia DHT Bootstrap (:4003) ] <-----------------------+  |                   |
|    - Peer Routing & Provider Indexing                                      |  |                   |
|    - Multi-Client Demand Aggregation & Surge Alerts                        |  |                   |
|                                                                            |  |                   |
|  [ NAT TRAVERSAL GATEWAY: Circuit Relay v2 (:4001) ]                       |  |                   |
|    - HOP Protocol Reservations & Hole-Punching Proxy                       |  |                   |
|                                                                            |  |                   |
|  [ 4 DIFFERENTIATED STORAGE PROVIDER TIERS ]                               |  |                   |
|    - Tier-1 Core Storage      (:4101) -> Direct Swarm, Ticket Verifier     |  |                   |
|    - Tier-2 Edge Cache        (:4102) -> Relay-Forced NAT Cache            |  | (Off-Chain        |
|    - Tier-3 Audit Guardian    (:4103) -> 5s Availability Prover            |  |  EIP-712 Tickets) |
|    - Tier-4 Hot Standby       (:4104) -> Disaster Recovery Parity Store    |  |                   |
|                                                                            |  |                   |
|  [ PUBLISHER NODE (:4201) ] -----------------------------------------------+--+                   |
|    - Ingestion, ChaCha20/AES-256 Encryption, Merkle Rooting                |                      |
|    - Push Protocol (/cipher/push/1.0.0) with Replication R=2               |                      |
|    - Periodic Availability Challenges & Closing Daemon Repayments          |                      |
|                                                                            |                      |
|  [ CONSUMER CLIENTS (:4301, :4302, :4303) ]                                |                      |
|    - Honest Swarm Client: Concurrent chunk download + EIP-712 ticket stream|                      |
|    - Malicious Fraud Client: Forged ticket simulation (tested & rejected)  |                      |
|    - Failover Recovery Client: Reconstructs payload during dead-node outage |                      |
+----------------------------------------------------------------------------+----------------------+
```

---

## 3. Comprehensive Test Suites & Tools Inventory

The repository includes five layers of verification suites:

### A. The 13-Checkpoint Local Architecture Orchestrator ([`local_multiple_terminal_test.sh`](../../local_multiple_terminal_test.sh))
Executes the full 10-node distributed network end-to-end:
1. Cryptographic binary compilation into `./bin/`
2. Cross-platform process isolation and port remediation
3. Local Anvil EVM deployment and 5.0 ETH escrow funding
4. Circuit Relay v2 initialization and HOP reservation
5. Kademlia DHT bootstrap root initialization
6. 4 Differentiated Storage Provider tiers startup & DHT registration
7. Pre-flight financial wallet balance audit
8. Ingestion, replication ($R=2$), and live availability challenge
9. Honest consumer swarm download with streaming EIP-712 tickets & SHA-256 validation
10. Malicious consumer cheating defense (forged signature detection & stake preservation)
11. Dead-node hardware crash simulation, Kademlia DHT demand alert & survivor recovery
12. On-chain EVM raffle winner calculation and 1.0 ETH prize transfer
13. Daemon continuous Proof of Storage audit and final 0.5 ETH publisher repayment

*Detailed reference*: [`docs/13_step_local_testing_execution.md`](../13_step_local_testing_execution.md).

### B. Master Verification Pipeline ([`test_workflow.sh`](../../test_workflow.sh))
Runs the complete CI verification pipeline:
- **Phase 1**: Solidity Smart Contracts (53 tests: 46 payments + 7 escrow)
- **Phase 2**: Go Cryptographic Signers, Identity & EIP-712 Verification
- **Phase 3**: Availability Cross-Domain Integration & Merkle Proof Engine
- **Phase 4**: Local Anvil EVM Blockchain Deployment
- **Phase 5**: Live P2P Chunk-for-Ticket Transfer & On-Chain Settlement

### C. Foundry Solidity Contract Test Suites
```bash
# Payment channels, lotteries, and dispute arbitration
(cd payments && forge test -v)

# Availability escrow, epochs, and collateral slashing
(cd availability/escrow-payment/escrow && forge test -v)
```

### D. Go Unit, Cryptography & Adversarial Wire Tests
```bash
# Push wire protocol framing, serializers, and limits
go test -v ./network/protocol/push/...

# Chunk protocol wire streaming, ticket payment, and forge rejection
go test -v ./network/protocol/chunk/...

# Secp256k1 Ethereum & Ed25519 P2P identity key management
go test -v ./network/identity/...

# EIP-712 typed data hashing and signature recovery
go test -v ./network/payments/...

# Content-Addressed Storage (CAS), chunking, and encryption
go test -v ./network/content/...

# Circular replication planner and replica placement
go test -v ./network/distribution/...

# Merkle proof trees and provider availability proof engines
go test -v ./integration/availability/...
```

### E. End-to-End System Integration Suites
* [`tests/e2e/payment_transfer_anvil_e2e.sh`](../../tests/e2e/payment_transfer_anvil_e2e.sh): P2P chunk transfer with live streaming EIP-712 tickets and Anvil EVM settlement.
* [`tests/e2e/remote_push_e2e.sh`](../../tests/e2e/remote_push_e2e.sh): Remote ingestion across 3 providers with replication $R=2$ and dead-node recovery.
* [`tests/e2e/roles_e2e.sh`](../../tests/e2e/roles_e2e.sh): Publisher, standalone provider, and consumer swarming via DHT discovery.
* [`tests/e2e/provider_lifecycle_e2e.sh`](../../tests/e2e/provider_lifecycle_e2e.sh): Provider crash, restart, and persistent content retrieval.
* [`tests/e2e/legacy_peer_transfer_e2e.sh`](../../tests/e2e/legacy_peer_transfer_e2e.sh): Direct peer-to-peer self-test.
* [`scripts/run_simulation.sh`](../../scripts/run_simulation.sh): Dynamic multi-epoch availability, demand spike, and on-chain collateral slashing simulator.

---

## 4. Required External Toolchains & Global Prerequisites

All testing environments require the following external tools installed and available in `PATH`:

| Tool | Minimum Version | Purpose |
| :--- | :--- | :--- |
| **Go** | `1.22+` | Compiling nodes, network protocols, and cryptographic modules |
| **Foundry (`forge`, `anvil`, `cast`)** | Latest Stable | Compiling Solidity smart contracts, local EVM ledger, and RPC transactions |
| **Bash** | `4.0+` | Executing test harnesses and automated multi-process orchestrators |
| **Git** | `2.30+` | Repository version control and branch checkout |
| **cURL** | Any modern | HTTP RPC status polling and health check verification |
| **SHA-256 Tool** | `sha256sum`, `shasum`, or `openssl` | Bit-for-bit payload integrity validation |
| **Port Utility** | `lsof` or `fuser` | Process isolation and TCP port conflict remediation |

---

## 5. Repository-Relative Path Standards

To guarantee that cloned repositories operate immediately on any user workstation without path configuration:

1. **Working Directory Anchoring**: All shell scripts resolve their root relative to the script location:
   ```bash
   SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
   if [ -f "$SCRIPT_DIR/go.mod" ]; then
       ROOT="$SCRIPT_DIR"
   else
       ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
   fi
   cd "$ROOT"
   ```
2. **Binary Output**: All compiled executables are placed in `./bin/` relative to the repository root.
3. **Data Stores**: Ephemeral storage directories (`./store_p1`, `./store_p2`, `./store_pub`, `./store_client`) are created and cleaned up within the repository root.
4. **Documentation Links**: All Markdown references use clean relative paths (e.g. `[testing.md](../testing.md)`).
