# CIPHER Testing Strategy & Quick Start Guide

This document provides a comprehensive technical guide to testing the CIPHER decentralized content delivery network. It covers unit testing, adversarial wire testing, automated end-to-end role validation, multi-provider replication, physical multi-laptop setups, and live EVM micropayment settlement on Anvil.

> [!IMPORTANT]
> **Notice for Local Testing:**
> All local multi-terminal testing, role simulations, and verification suites must be executed against the **`local`** branch of the repository:
> 👉 **[https://github.com/devlup-labs/CIPHER/tree/local](https://github.com/devlup-labs/CIPHER/tree/local)**
>
> To clone and prepare the test environment locally (macOS, Linux, Windows WSL2 / Git Bash):
> ```bash
> git clone https://github.com/devlup-labs/CIPHER.git
> cd CIPHER
> git checkout local
> ```

---

## 📖 Dedicated Platform Guides

Detailed installation instructions, system permissions, external dependencies, and OS-specific run commands are available in dedicated platform guides:

* 🌐 **[General Architecture & Verification Strategy](testing/general.md)**
* 🪟 **[Windows Testing Guide (WSL2 / Git Bash)](testing/windows.md)**
* 🐧 **[Linux Testing Guide (Ubuntu, Debian, Fedora, Arch)](testing/linux.md)**
* 🍏 **[macOS Testing Guide (Apple Silicon & Intel)](testing/mac.md)**

---

## ⚡ Quick Start: Fast Automated Testing

All test suites and orchestrators run locally across macOS, Linux, and Windows (WSL2 / Git Bash):

```bash
# 1. [COMPLETE 13-STEP ORCHESTRATOR] 10-Role Desktop Tiling & Full Pipeline Test
./local_multiple_terminal_test.sh --auto            # Option 1: macOS 10-window desktop tiled layout (Automated non-stop)
./local_multiple_terminal_test.sh                  # Option 2: macOS 10-window desktop tiled layout (Interactive step-by-step)
./local_multiple_terminal_test.sh --single --auto   # Option 3: Universal single-terminal mode (Automated non-stop CI - Linux, macOS, WSL2)
./local_multiple_terminal_test.sh --single          # Option 4: Universal single-terminal mode (Interactive step-by-step)

# 2. [MASTER WORKFLOW] Run entire verification pipeline (Solidity + Go Unit + Adversarial Wire + Live Anvil Settlement)
./test_workflow.sh

# 3. Run all Go unit, identity, payment, and wire protocol tests
go test ./network/...

# 4. Run Foundry Solidity smart contract test suites (46 tests)
(cd payments && forge test)

# 5. Run Live Anvil P2P Transfer & On-Chain Settlement E2E Test
./tests/e2e/payment_transfer_anvil_e2e.sh

# 6. Run Remote Ingestion & Multi-Provider Replication Test (3 Providers, R=2, Fault Kill)
./tests/e2e/remote_push_e2e.sh

# 7. Run Role-Based DHT Swarming Test (Publisher, Provider, Consumer, Bootstrap)
./tests/e2e/roles_e2e.sh

# 8. Run Provider Persistence & Independence Test
./tests/e2e/provider_lifecycle_e2e.sh
```

> [!NOTE]
> For the complete technical breakdown of the 10-node layout and all 13 checkpoints, see the dedicated [13-Step Local Testing Execution Guide](13_step_local_testing_execution.md).

---

## 🏗️ The CIPHER Roles Under Test

| Role | Source Path | Responsibilities |
| :--- | :--- | :--- |
| **Publisher** | [`nodes/publisher/main.go`](../nodes/publisher/main.go) | Ingests source file, chunks & encrypts via XChaCha20, creates manifest, plans placement, pushes to remote providers via `/cipher/push/1.0.0` with replication $R$, and exits. |
| **Provider** | [`nodes/provider/main.go`](../nodes/provider/main.go) | Standalone daemon hosting Content-Addressed Storage (CAS). Accepts uploads, announces to DHT, serves chunks (`/cipher/chunk/1.0.0`), cryptographically verifies EIP-712 payment tickets, and settles rounds on-chain. |
| **Consumer** | [`nodes/consumer/main.go`](../nodes/consumer/main.go) | Discovers candidate providers via DHT (or direct dial `-d`), resolves manifest, swarms chunks concurrently via worker pool, verifies ciphertext hashes, signs & streams EIP-712 payment tickets, decrypts, and reassembles payload. |
| **Bootstrap** | [`network/cmd/bootstrap/main.go`](../network/cmd/bootstrap/main.go) | Kademlia DHT bootstrap routing node for decentralized provider discovery. |
| **Relay** | [`network/cmd/relay/main.go`](../network/cmd/relay/main.go) | Circuit v2 Relay node for NAT traversal and DCUtR hole punching coordination. |

---

## 🧪 Core Test Scenarios

### Scenario 1: Live Anvil EVM Settlement & Micro-Payment Ticket Verification
**Objective**: Prove that a consumer can stream EIP-712 signed lottery tickets (`MsgTicket 0x07`) over `/cipher/chunk/1.0.0` as chunks are downloaded from a provider, that the provider verifies them in real-time, and that accumulated winning tickets can be settled on-chain against `SettlementEngine.sol` on a live EVM node (Anvil).

```text
[Step 1] Spawn local Anvil node (JSON-RPC on 127.0.0.1:8545)
[Step 2] Run payments/script/Step1_Setup.s.sol:
         - Deploys CommitRevealEntropy, PaymentChannel, SettlementEngine
         - Registers Provider in ProviderRegistry
         - Opens consumer Payment Channel with 1 ETH deposit
[Step 3] Start Bootstrap node (:4001)
[Step 4] Provider daemon starts with --eth-rpc and --entropy-addr flags
[Step 5] Publisher ingests 2 MB test file and pushes to Provider
[Step 6] Consumer fetches content with:
         --eth-rpc "http://127.0.0.1:8545"
         --eth-key <CONSUMER_KEY>
         --entropy-addr <ENTROPY_CONTRACT_ADDR>
         --provider-eth-addr <PROVIDER_ETH_ADDR>
         - For each chunk received, signs and streams EIP-712 RoundTicket (MsgTicket 0x07)
         - Provider recovers consumer address and verifies channel balance
[Step 7] Consumer decrypts and verifies SHA-256 hash matches 100%
[Step 8] Run payments/script/Step2_Settle.s.sol:
         - Settle round on-chain: verifies winning ticket payout to Provider
```

**Automated Command**:
```bash
./tests/e2e/payment_transfer_anvil_e2e.sh
```

---

### Scenario 2: Adversarial Ticket & Security Validation
**Objective**: Prove that the wire protocol strictly rejects forged tickets, corrupted signatures, mismatched channel IDs, and invalid nonces without panicking or serving unauthorized data.

* **Test Suite**: `network/protocol/chunk/ticket_test.go`
* **Test Cases**:
  1. `TestChunkWithValidTicket`: Client signs genuine EIP-712 ticket with Secp256k1 key; provider recovers signer address, verifies ticket matches expected channel, and accepts transfer.
  2. `TestChunkWithForgedTicket`: Client tampers with the signature bytes; provider's `ecrecover` detects invalid signature, returns verification failure, and logs rejection.
  3. `TestChunkWithTamperedNonce`: Client sends a decremented/replayed nonce; provider detects replay attempt and rejects ticket.
  4. `TestChunkWithExceededLimit`: Client attempts to send oversized payload exceeding `MaxTicketSize` (4096 bytes); stream reader enforces protocol limit and rejects packet.

**Automated Command**:
```bash
go test -v ./network/protocol/chunk -run TestChunkWith
```

---

### Scenario 3: Remote Ingestion & Multi-Provider Replication
**Objective**: Prove that the Publisher can push content to $N$ independent remote providers across the network with Replication Factor $R$, verify the replication invariant, and shut down. Then prove that consumers can swarm-download the file even if a provider node is terminated midway.

```text
[Step 1] Bootstrap Node starts on port 48001
[Step 2] 3 Providers start with isolated directories (store_p1, store_p2, store_p3)
[Step 3] Publisher ingests 2 MB test file into an isolated store_pub
[Step 4] Publisher pushes chunks to Provider 1, 2, 3 with Replication R=2:
         - Circular placement assigns each chunk to 2 distinct providers
         - GlobalReplicaTracker confirms all 64 chunks reached >= 2 replicas
         - Providers announce ContentID to DHT on batch completion
         - Publisher exits cleanly (seed=false)
[Step 5] Consumer 1 discovers providers via DHT, swarms chunks from all 3 providers,
         decrypts, and verifies SHA-256 hash matches 100%
[Step 6] Fault Tolerance: Provider 1 is killed (kill $PROV1_PID)
[Step 7] Consumer 2 downloads file from surviving Providers 2 & 3 via DHT discovery
         and reassembles with 100% hash equality
```

**Automated Command**:
```bash
./tests/e2e/remote_push_e2e.sh
```

---

### Scenario 4: Independent Provider Lifecycle (Publisher Offline)
**Objective**: Prove that once content is published to a Provider, the **Publisher is completely removed from the retrieval path**.

```text
[Step 1] Bootstrap Node starts (DHT mesh coordinator)
[Step 2] Publisher ingests 2 MB test file into Provider CAS store
[Step 3] Publisher process is completely KILLED (verified 100% offline)
[Step 4] Provider starts independently with existing CAS store & registers with DHT
[Step 5] Consumer queries DHT solely with ContentID + Key + Bootstrap address
         - Discovers Provider on DHT (no direct peer address passed)
         - Resolves Manifest from Provider over /cipher/chunk/1.0.0
         - Downloads all chunks in parallel & decrypts payload
         - Validates SHA-256 checksum equality (100% match)
```

**Automated Command**:
```bash
./tests/e2e/provider_lifecycle_e2e.sh
```

---

### Scenario 5: Multi-Provider Swarming (Candidate Provider Model)
**Objective**: Prove that the Consumer's `TransferManager` and `Scheduler` parallelize chunk downloads across multiple independent providers, gracefully handling partial chunk assignments.

```text
                     DHT (Control Plane)
                              │
                    FindProviders(ContentID)
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
               Provider A          Provider B
            (Chunks 0, 1, 4)    (Chunks 2, 3, 5)
                    │                   │
                    │   /cipher/chunk   │
                    └───►  Consumer ◄───┘
```

#### Candidate Miss Semantics:
When a consumer requests a chunk from Provider A that Provider A does not host, Provider A returns `ErrChunkNotFound (0x02)`. The consumer scheduler identifies this as a **candidate miss**, marks Provider A in `MissedPeers` for that task, and immediately routes the task to Provider B without incrementing failure penalties.

---

### Scenario 6: Public Relay & DCUtR Hole Punching (NAT Traversal)
**Objective**: Connect across NATs/firewalls via a public `circuitv2` relay and verify automatic, seamless upgrade to direct TCP/UDP sockets via DCUtR.

1. **Start Public Relay (on cloud VM / Azure / Ubuntu)**:
   ```bash
   go run network/cmd/relay/main.go
   ```
   *Copy relay multiaddress:* `/ip4/<PUBLIC_IP>/tcp/4001/p2p/<RELAY_PEER_ID>`
2. **Start Provider behind NAT**:
   ```bash
   go run nodes/provider/main.go -p 4001 -store ./provider_store \
     -relay "/ip4/<PUBLIC_IP>/tcp/4001/p2p/<RELAY_ID>" -bootstrap "<BOOTSTRAP_ADDR>"
   ```
3. **Start Consumer on separate machine / network**:
   ```bash
   go run nodes/consumer/main.go -relay "/ip4/<PUBLIC_IP>/tcp/4001/p2p/<RELAY_ID>" \
     -bootstrap "<BOOTSTRAP_ADDR>" -fetch "<CID>" -key "<KEY>" -out output.mp4
   ```
4. **Verify DCUtR Upgrade**:
   - Initial connection is established over the limited relay circuit.
   - DCUtR executes simultaneous hole punch in the background (`[DCUtR] Hole Punch Event: StartHolePunch` $\rightarrow$ `EndHolePunch`).
   - Transfer shifts to high-throughput direct socket (`Path: Direct`).

---

## 🏆 The Master Workflow Test Runner (`test_workflow.sh`)

To execute the entire verification pipeline with a single command, run:

```bash
./test_workflow.sh
```

### Pipeline Execution Stages:
1. **Prerequisite & Environment Check**:
   - Validates `go`, `forge`, `anvil`, and `cast` binaries.
   - Ensures no conflicting processes occupy Anvil port `8545` or P2P ports.
2. **Phase 1: Solidity Smart Contract Suites**:
   - Runs `(cd payments && forge test)` across all 46 test cases.
3. **Phase 2: Go Unit, Identity & Payments Cryptography**:
   - Runs `go test ./network/identity/...` (Secp256k1 key generation, 0600 permissions, address derivation).
   - Runs `go test ./network/payments/...` (EIP-712 typed hashing and OpenZeppelin signature recovery).
4. **Phase 3: Adversarial Wire Protocol Tests**:
   - Runs `go test ./network/protocol/chunk/...` (EIP-712 `MsgTicket 0x07` streaming, signature verification, and forge rejection).
5. **Phase 4: Live Anvil End-to-End P2P Settlement**:
   - Launches local Anvil node in background.
   - Executes `Step1_Setup.s.sol` to deploy contracts and fund channel.
   - Launches Provider daemon with `--eth-rpc` and `--entropy-addr`.
   - Runs Consumer swarming download with streaming payment tickets.
   - Decrypts and verifies 100% SHA-256 binary hash match.
   - Executes `Step2_Settle.s.sol` on-chain settlement.
   - Cleanly terminates all child processes and deletes transient test stores.

---

## 💻 Hardware & Multi-Laptop Testing Guide

For genuine real-world validation with physical machines:

### Option A: The 2-Laptop Setup (Minimum Hardware)
* **Laptop 1** (e.g. `192.168.1.50` on Home Wi-Fi):
  - Terminal 1: `go run network/cmd/bootstrap/main.go -p 48001 -ws-port 0`
  - Terminal 2: `go run nodes/provider/main.go -p 48010 -ws-port 0 -store ./store_p1 -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>"`
  - Terminal 3: `go run nodes/provider/main.go -p 48020 -ws-port 0 -store ./store_p2 -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>"`
* **Laptop 2** (e.g. `192.168.1.80` or Mobile Hotspot):
  - Push content:
    ```bash
    go run nodes/publisher/main.go -file test.mp4 -push -replication 2 \
      -providers "/ip4/192.168.1.50/tcp/48010/p2p/<P1_ID>,/ip4/192.168.1.50/tcp/48020/p2p/<P2_ID>" \
      -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>" -seed=false
    ```
  - Retrieve content via DHT:
    ```bash
    go run nodes/consumer/main.go -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>" \
      -fetch "<CONTENT_ID>" -key "<KEY>" -out downloaded.mp4
    ```

### Option B: The 3-Machine Setup (VPS + 2 NAT Laptops)
* **Cloud VM / VPS** (Public IP `20.x.x.x`): Runs Bootstrap (`:4003`) and Relay (`:4001`).
* **Laptop A** (Home Wi-Fi): Runs Providers 1 & 2 pointing to VPS Relay and Bootstrap.
* **Laptop B** (Mobile Hotspot): Runs Publisher & Consumer. Traverses NAT via DCUtR hole punching to download directly from Laptop A.

---

## 📊 Unit & Robustness Test Suites

```bash
# Push wire protocol suite (framing, limits, serializers)
go test -v ./network/protocol/push/...

# Distribution & Replication suite (circular planner, replica tracker)
go test -v ./network/distribution/...

# Content Engine suite (chunking, ChaCha20, digests, CAS)
go test -v ./network/content/...

# Chunk protocol suite (wire framing, handlers, ACKs, tickets)
go test -v ./network/protocol/chunk/...

# Identity & Payment Cryptography suites
go test -v ./network/identity/...
go test -v ./network/payments/...
```

---

## 🛠️ CLI Quick Reference

```bash
# Publisher: Push file to remote providers with R=2 replication and exit:
go run nodes/publisher/main.go -file <file> -push -replication 2 \
  -providers "<P1_ADDR>,<P2_ADDR>" -bootstrap <boot_addr> -seed=false

# Provider: Start storage node with push enabled and EVM ticket verification:
go run nodes/provider/main.go -p 4001 -ws-port 4002 -store <store_dir> \
  -bootstrap <boot_addr> -allow-push=true \
  --eth-rpc "http://127.0.0.1:8545" --entropy-addr "<ENTROPY_ADDR>"

# Consumer: Fetch content via DHT swarming with streaming payment tickets:
go run nodes/consumer/main.go -bootstrap <boot_addr> -fetch <CID> -key <KEY_HEX> -out <output_path> \
  --eth-rpc "http://127.0.0.1:8545" --eth-key "<ETH_PRIVKEY>" \
  --entropy-addr "<ENTROPY_ADDR>" --provider-eth-addr "<PROV_ETH_ADDR>"

# Consumer: Query active download sessions:
go run nodes/consumer/main.go -store <store_dir> -status

# Consumer: Cancel a download session:
go run nodes/consumer/main.go -store <store_dir> -cancel <CID>
```
