# CIPHER: 13-Step Local Testing & Architecture Execution Guide

This document specifies the end-to-end local test orchestration for the CIPHER decentralized content delivery network (CDN) and probabilistic micropayment architecture. It details the 10-node perimeter tiling matrix, role differentiation, cryptographic verification phases, and live terminal outputs across all 13 checkpoints executed by [`local_multiple_terminal_test.sh`](../local_multiple_terminal_test.sh).

> [!IMPORTANT]
> **Notice for Local Testing:**
> All local multi-terminal testing, role simulations, and verification suites must be executed against the **`local`** branch of the repository:
> 👉 **[https://github.com/devlup-labs/CIPHER/tree/local](https://github.com/devlup-labs/CIPHER/tree/local)**
>
> To clone and prepare the test environment locally:
> ```bash
> git clone https://github.com/devlup-labs/CIPHER.git
> cd CIPHER
> git checkout local
> ```

---

## 1. Desktop Perimeter Tiling Matrix & Cross-Platform Execution

The test runner [`local_multiple_terminal_test.sh`](../local_multiple_terminal_test.sh) is engineered to run seamlessly across **macOS**, **Linux** (Ubuntu/Debian/Fedora/Arch), and **Windows** (WSL2 / Git Bash / MSYS2).

### Operating System Support
* **macOS**: Spawns 10 dedicated `Terminal.app` windows tiled around the screen perimeter (see [macOS Testing Guide](testing/mac.md)).
* **Linux**: Automatically runs in unified single-terminal background orchestrator mode (see [Linux Testing Guide](testing/linux.md)).
* **Windows (WSL2 / Git Bash)**: Automatically runs in unified single-terminal background orchestrator mode (see [Windows Testing Guide](testing/windows.md)).
* **General Architecture & Verification Strategy**: See [General Testing Guide](testing/general.md).

```text
+---------------------+---------------------+---------------------+---------------------+
| [Slot 1] ANVIL EVM  | [Slot 2] RELAY V2   | [Slot 3] BOOTSTRAP  | [Slot 4] TIER-1 CORE|
| Port: 8545          | Port: 4001 (HOP)    | Port: 4003 (DHT)    | Port: 4101 (Direct) |
| L1 Settlement & Escrow NAT Traversal Proxy  DHT Root Routing Table Ticket Verifier    |
+---------------------+---------------------+---------------------+---------------------+
| [Slot 5] TIER-2 EDGE|                                           | [Slot 6] TIER-3 AUDIT
| Port: 4102          |           [ CENTER DESKTOP AREA ]         | Port: 4103          |
| NAT-Relayed Cache   |        (Main IDE / Controller Workspace)  | Availability Prover |
+---------------------+                                           +---------------------+
| [Slot 7] TIER-4 STBY| [Slot 8] PUBLISHER  | [Slot 9] CONSUMER 1 | [Slot 10] CONSUMER 2|
| Port: 4104          | Port: 4201          | Port: 4301          | Port: 4302 / 4303   |
| Failover Standby    | Ingest, Disperse, 5s| Swarm Downloader &  | Fraud Defense &     |
| Replica Store       | Availability Audit  | EIP-712 Lottery Pay | Failover Survivor   |
+---------------------+---------------------+---------------------+---------------------+
```

### Execution Modes
```bash
# Option 1: macOS 10-Window Desktop Tiled Layout (Automated non-stop execution)
./local_multiple_terminal_test.sh --auto

# Option 2: macOS 10-Window Desktop Tiled Layout (Interactive step-by-step with checkpoints)
./local_multiple_terminal_test.sh

# Option 3: Universal Single-Terminal Mode (Automated non-stop CI - Linux, macOS, WSL2)
./local_multiple_terminal_test.sh --single --auto

# Option 4: Universal Single-Terminal Mode (Interactive step-by-step with checkpoints)
./local_multiple_terminal_test.sh --single
```

---

## 2. The 13 Checkpoints: Deep Technical Breakdown

### Checkpoint 1/13: Cryptographic Integrity & Binary Compilation
* **Purpose**: Compiles all specialized nodes and verification tools into `./bin/`, ensuring zero compilation errors across libp2p, Foundry bindings, Merkle trees, and cryptographic modules.
* **Binaries Built**:
  - `./bin/bootstrap`: Kademlia DHT bootstrap root node.
  - `./bin/relay`: libp2p Circuit Relay v2 proxy.
  - `./bin/provider`: Storage provider daemon with EIP-712 ticket validation & Proof of Storage.
  - `./bin/publisher`: Ingestion, encryption, chunk dispersal, and 5-second challenge loop.
  - `./bin/consumer`: Swarm downloader, EIP-712 ticket streamer, and malicious fraud simulator.
  - `./bin/dht-inspect`: Kademlia routing table inspector and multi-client demand tracker.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 1/13] CRYPTOGRAPHIC INTEGRITY & COMPILED BINARIES        
======================================================================
  Compiling 6 specialized node & inspector binaries...
[✓] All 6 binaries (including Kademlia DHT Inspector) compiled cleanly in ./bin/
```

---

### Checkpoint 2/13: Process Isolation & Port Remediation
* **Purpose**: Prevents port collisions and ensures clean state by inspecting and terminating any orphaned background daemons on TCP ports `8545`, `4001`, `4003`, `4101`, `4102`, `4103`, `4104`, `4201`, `4301`, `4302`, `4303`.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 2/13] PROCESS ISOLATION & CLEAN NETWORK STATE            
======================================================================
[✓] Clean network slate verified across all 10 CDN ports.
```

---

### Checkpoint 3/13: Layer 1 EVM Settlement & Escrow Engine (Terminal 1)
* **Terminal Role**: `[Terminal 1/10] Anvil EVM Blockchain (Port 8545)`
* **Under the Hood**:
  1. Starts a local EVM node via Anvil at `http://127.0.0.1:8545`.
  2. Executes `payments/script/Step1_Setup.s.sol` to deploy `CommitRevealEntropy.sol`, `PaymentChannel.sol`, `SettlementEngine.sol`, and `ProviderRegistry.sol`.
  3. Registers the provider in `ProviderRegistry` and deposits 5.0 ETH into the payment escrow channel for client `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC`.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 3/13] [Terminal 1/10] LAYER 1 EVM SETTLEMENT & ESCROW     
======================================================================
[✓] Anvil EVM is live on http://127.0.0.1:8545
Deploying Smart Contracts & Initializing 5.0 ETH Escrow Deposit...
[✓] Smart Contracts deployed and funded:
  - EntropySource:       0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
  - Escrow Channel:      0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
  - Provider Eth Wallet: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
  - Client Eth Wallet:   0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
```

---

### Checkpoint 4/13: NAT Traversal & Circuit Relay v2 Gateway (Terminal 2)
* **Terminal Role**: `[Terminal 2/10] Circuit Relay v2 (Port 4001)`
* **Under the Hood**:
  - Initializes a libp2p Circuit Relay v2 HOP router allowing NAT-firewalled providers and clients to communicate without direct public IP routing.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 4/13] [Terminal 2/10] [NAT GATEWAY] CIRCUIT RELAY V2     
======================================================================
[✓] Relay Multiaddress: /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWFWUGbJQujTc1C64EHGsHjG5CYtdTUfUUpu77c738Q1eW
```

---

### Checkpoint 5/13: Kademlia DHT Control-Plane Discovery Hub (Terminal 3)
* **Terminal Role**: `[Terminal 3/10] Kademlia DHT Bootstrap Hub (Port 4003)`
* **Under the Hood**:
  - Starts the root DHT node maintaining the Kademlia routing table, rendezvous namespace, and content provider indexing records.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 5/13] [Terminal 3/10] [DHT ROUTER] KADEMLIA BOOTSTRAP    
======================================================================
[✓] Bootstrap Multiaddress: /ip4/127.0.0.1/tcp/4003/p2p/12D3KooWESxNHxkdCQgkhvDffk2n8h1FDPymwe57u8Po4hktzJi7
```

---

### Checkpoint 6/13: 4 Differentiated Storage Provider Tiers (Terminals 4-7)
* **Terminal Roles**:
  - `[Terminal 4]`: Tier-1 Core Storage & Primary Ticket Verifier (Direct Port 4101).
  - `[Terminal 5]`: Tier-2 Edge Cache & NAT-Firewalled Node (Relayed Port 4102 via `--force-relay`).
  - `[Terminal 6]`: Tier-3 Audit Guardian & Availability Prover (Challenge Port 4103).
  - `[Terminal 7]`: Tier-4 Hot Standby Disaster Recovery Replica (Failover Port 4104).
* **Under the Hood**:
  - Each provider connects to the DHT bootstrap, announces its availability under the storage provider namespace, and enables content-addressed chunk storage with ticket verification.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 6/13] [Terminals 4-7/10] 4 DIFFERENTIATED STORAGE TIERS  
======================================================================
[✓] 4 differentiated storage tiers active and verified across ports 4101-4104.

[*] Auditing Kademlia DHT Control-Plane Routing Table & Storage Registrations...
[DHT-INSPECTOR] [✓] Discovered 4 active storage provider(s) on Kademlia DHT:
   [1] Provider Peer ID : 12D3KooWMdnp1ejTXDHAoFmYsfYHrXRekKVNM8tWgX8AJask9RPF
   [2] Provider Peer ID : 12D3KooWEmHsr6sdSnrsKvWsVzJ9cQSKVT9jeJBT5VpPQA7X5dvU
   [3] Provider Peer ID : 12D3KooWDVriJM3ys761G9tWRybJaeiRmqJ8S23kkp3ipWMQGNWQ
   [4] Provider Peer ID : 12D3KooWGMkXtESkyagSRphUtbJg9osYYX2c1U2g3YJmHUoDed4N
```

---

### Checkpoint 7/13: Pre-Flight Financial Ledger Audit
* **Purpose**: Records initial on-chain balances prior to content transfer, verifying that the client holds funds and the escrow channel is funded.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 7/13] PRE-FLIGHT FINANCIAL LEDGER AUDIT                  
======================================================================
  💳 Client Wallet Balance  : 9995.0 ETH
  💳 Provider Wallet Balance: 9997.0 ETH
  🏦 Escrow Channel Deposit : 5.0 ETH
```

---

### Checkpoint 8/13: Publisher Ingestion, Multi-Tier Dispersal & 5s Challenge Loop (Terminal 8)
* **Terminal Role**: `[Terminal 8/10] Publisher Ingestion & Demand Auditor (Port 4201)`
* **Under the Hood**:
  1. Generates 1 MB payload, encrypts with AES-256-GCM, and computes Merkle root.
  2. Plans circular redundant placement across all 4 providers with Replication Factor $R=2$.
  3. Pushes chunks over `/cipher/push/1.0.0` streams and renders the **Publisher Dispersal Table**.
  4. Runs the **5-Second Availability Challenge Loop**: challenges providers for chunk possession, verifies Merkle paths against the root, and dispatches optimistic `PaymentState` vouchers over `/cipher/availability/1.0.0`.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 8/13] [Terminal 8/10] PUBLISHER INGESTION & DEMAND AUDIT 
======================================================================
  Generated 1 MB Random Payload SHA-256: dfe18c7fb64a24fc4dbb8f67fc1dffe30206aa4bff3cd43e5813f140d017cfbe

+---------------------------------------------------------------------------------------------------------+
|                                 PUBLISHER DISPERSAL & PLACEMENT TABLE                                   |
+------------------------------------+---------------+-------------+---------------+----------------------+
| Storage Provider Peer ID           | Replicas (Qty)| Share (%)   | Volume (KB)   | Status               |
+------------------------------------+---------------+-------------+---------------+----------------------+
| 12D3KooWMdnp1ejTXDHAoFmYsfYHrXRekKVNM8tWgX8AJask9RPF | 16 chunks     | 25.0%       | 512.0 KB      | COMMITTED [✓]        |
| 12D3KooWEmHsr6sdSnrsKvWsVzJ9cQSKVT9jeJBT5VpPQA7X5dvU | 16 chunks     | 25.0%       | 512.0 KB      | COMMITTED [✓]        |
| 12D3KooWDVriJM3ys761G9tWRybJaeiRmqJ8S23kkp3ipWMQGNWQ | 16 chunks     | 25.0%       | 512.0 KB      | COMMITTED [✓]        |
| 12D3KooWGMkXtESkyagSRphUtbJg9osYYX2c1U2g3YJmHUoDed4N | 16 chunks     | 25.0%       | 512.0 KB      | COMMITTED [✓]        |
+------------------------------------+---------------+-------------+---------------+----------------------+
| TOTAL CLUSTER REPLICATION (R=2)    | 64 replicas   | 100.0%      | 2048.0 KB     | Invariant: SATISFIED [✓] |
+------------------------------------+---------------+-------------+---------------+----------------------+

[PUBLISHER] [AVAILABILITY] [INFO] Starting Availability challenge engine (Interval: 5s, Rounds: 1)...
[PUBLISHER] [AVAILABILITY] [✓] [Round 1] Provider 12D3KooWMdnp verified chunk possession (Merkle Proof PASS for chunk 0)
[PUBLISHER] [PAYMENT]      [✓] [Round 1] Dispatched Optimistic Payment Voucher to 12D3KooWMdnp (Seq: 1, Cumulative: 25000000000000000 wei)
```

---

### Checkpoint 9/13: Honest Consumer Swarm Download & Metrics Table (Terminal 9)
* **Terminal Role**: `[Terminal 9/10] Honest Consumer 1 (Port 4301)`
* **Under the Hood**:
  1. Resolves manifest from DHT and initiates concurrent chunk swarm across providers.
  2. For each chunk received, signs and streams EIP-712 lottery tickets (`MsgTicket 0x07`).
  3. Verifies chunk ciphertext SHA-256 hashes, decrypts with key, and reassembles original file.
  4. Renders the **Swarm Retrieval Metrics Table** showing per-provider chunk counts, ratios, byte volumes, and download latencies.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 9/13] [Terminal 9/10] HONEST CONSUMER SWARM & TICKETS   
======================================================================
[Progress] 32/32 chunks (100.0%)

+---------------------------------------------------------------------------------------------------------+
|                                    SWARM RETRIEVAL METRICS TABLE                                        |
+------------------------------------+---------------+-------------+---------------+----------------------+
| Provider Peer ID                   | Chunks (Qty)  | Ratio (%)   | Volume (KB)   | Time Taken (ms)      |
+------------------------------------+---------------+-------------+---------------+----------------------+
| 12D3KooWMdnp1ejTXDHAoFmYsfYHrXRekKVNM8tWgX8AJask9RPF | 8 chunks      | 25.0%       | 256.1 KB      | 12 ms                |
| 12D3KooWEmHsr6sdSnrsKvWsVzJ9cQSKVT9jeJBT5VpPQA7X5dvU | 8 chunks      | 25.0%       | 256.1 KB      | 8 ms                 |
| 12D3KooWDVriJM3ys761G9tWRybJaeiRmqJ8S23kkp3ipWMQGNWQ | 8 chunks      | 25.0%       | 256.1 KB      | 7 ms                 |
| 12D3KooWGMkXtESkyagSRphUtbJg9osYYX2c1U2g3YJmHUoDed4N | 8 chunks      | 25.0%       | 256.1 KB      | 5 ms                 |
+------------------------------------+---------------+-------------+---------------+----------------------+
| TOTAL RECONSTRUCTED                | 32 chunks     | 100.0%      | 1024.5 KB     | Overall: 124 ms      |
+------------------------------------+---------------+-------------+---------------+----------------------+
  Original Payload SHA-256 : dfe18c7fb64a24fc4dbb8f67fc1dffe30206aa4bff3cd43e5813f140d017cfbe
  Downloaded File  SHA-256 : dfe18c7fb64a24fc4dbb8f67fc1dffe30206aa4bff3cd43e5813f140d017cfbe
[✓] 100% BIT-FOR-BIT DATA INTEGRITY CONFIRMED!
```

---

### Checkpoint 10/13: Malicious Consumer Fraud Defense & Staking Preservation (Terminal 10)
* **Terminal Role**: `[Terminal 10/10] Malicious Consumer 2 - Fraud Simulation (Port 4303)`
* **Under the Hood**:
  - Malicious consumer attempts to steal chunks by presenting forged EIP-712 signatures (`-simulate-cheat`).
  - Storage providers verify signatures via `ecrecover`, detect forged cryptography, log `[SECURITY SHIELD]`, immediately terminate the transfer stream, and preserve the on-chain escrow balance.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 10/13] [Terminal 10/10] CHEATING DEFENSE & STAKING SHIELD
======================================================================
Simulating Malicious Consumer 2 attempting to steal chunks using FORGED EIP-712 tickets...

[✓] SECURITY SHIELD VERIFIED:
  - Forged EIP-712 payment tickets were detected and REJECTED by storage providers.
  - Chunk transfers were denied to fraudulent client.
  - On-Chain Escrow Channel deposit remains 100% SECURE & PRESERVED against theft!
```

---

### Checkpoint 11/13: Dead-Node Fault Tolerance, Kademlia DHT Alert & Fallback Recovery (Terminal 10)
* **Terminal Role**: `[Terminal 10/10] Failover Recovery Client (Port 4302)`
* **Under the Hood**:
  1. Provider 1 (Port 4101) is killed to simulate a sudden hardware crash.
  2. 2–3 concurrent client search queries query the Kademlia DHT.
  3. The DHT Inspector measures search pressure, detects the missing provider, and broadcasts a **High-Demand Alert** to the publisher.
  4. Consumer 3 connects to the remaining 3 surviving providers (Edge, Audit, Standby) and successfully reconstructs 100% of the data using the $R=2$ replica parity.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 11/13] [Terminal 10/10] [FAILOVER AUDIT] SURVIVOR SWARM  
======================================================================
[!] Killing Provider 1 (Port 4101) to simulate node crash...

[*] 2-3 Clients searching Kademlia DHT for ContentID dfe18c... while Provider 1 is down...
================================================================================
                    KADEMLIA DHT HIGH-DEMAND BROADCAST ALERT                    
================================================================================
Queried ContentID  : dfe18c7fb64a24fc4dbb8f67fc1dffe30206aa4bff3cd43e5813f140d017cfbe
Search Query Count : 3 concurrent client searches
Active Providers   : 3 replica providers responding
Dropped Providers  : 1 provider node(s) offline
Demand/Cap Ratio   : 1.00 (HIGH DEMAND DETECTED)
Kademlia Broadcast : HIGH_CONTENT_DEMAND_ALERT -> Publisher
================================================================================
[✓] SUCCESS: 100% Content reconstructed from surviving replica tiers (R=2)!
```

---

### Checkpoint 12/13: On-Chain EVM Dispute & Raffle Settlement
* **Under the Hood**:
  1. Mines 20 confirmation blocks on Anvil.
  2. Executes `payments/script/Step2_Settle.s.sol`.
  3. Evaluates blockhash entropy from `CommitRevealEntropy.sol`, proves winning chunk index on-chain, and transfers 1.0 ETH from escrow to provider.
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 12/13] ON-CHAIN EVM DISPUTE & RAFFLE SETTLEMENT          
======================================================================
Mining 20 blocks on Anvil to pass confirmation window...
Submitting winning ticket on-chain to EscrowChannel...

== Logs ==
  === STEP 5: ON-CHAIN RAFFLE WINNER CALCULATION ===
    [+] Target Blockhash:    0x6afc5b1c3583e9661cc4f72dcdeb512c77407ea4963573543493db20abc953fb
    [+] Winning Chunk Index: 8
  
  === STEP 6: EXECUTING ON-CHAIN SETTLEMENT ===
    [+] Provider Bal Before: 9997 ETH
    [+] Channel Bal Before:  5 ETH
  ==================================================
  SUCCESS! ROUND SETTLED ON LIVE ANVIL EVM
  ==================================================
    [+] Provider Bal After:  9998 ETH
    [+] Channel Bal After:   4 ETH
    [+] NET PAYOUT EARNED:   1 ETH
  ==================================================
```

---

### Checkpoint 13/13: Daemon Proof of Storage Audit & Final Publisher Repayment
* **Under the Hood**:
  1. Publisher executes a live Availability challenge against surviving provider daemons using `-challenge-cid "$CONTENT_ID"`.
  2. Provider daemons prove continuous retention of data on disk by submitting sibling Merkle audit paths and cryptographic signatures.
  3. Publisher verifies all Merkle proofs (`[✓] PASS`), dispatches final closing vouchers, and executes an on-chain storage reward repayment (0.5 ETH) to the provider.
  4. Displays the final financial audit ledger showing combined earnings (+1.5 ETH total).
* **Terminal Output**:
```text
======================================================================
 [CHECKPOINT 13/13] DAEMON PROOF OF STORAGE & PUBLISHER REPAYMENT    
======================================================================
[*] Publisher auditing continuous Proof of Storage across surviving provider daemons...
[PUBLISHER] [AVAILABILITY] [✓] Loaded manifest for ContentID: dfe18c... (Total Chunks: 32)
[PUBLISHER] [AVAILABILITY] [✓] [Round 1] Provider 12D3KooWNQ4Q verified chunk possession (Merkle Proof PASS for chunk 0)
[PUBLISHER] [PAYMENT]      [✓] [Round 1] Dispatched Optimistic Payment Voucher to 12D3KooWNQ4Q (Seq: 1, Cumulative: 25000000000000000 wei)

[✓] Continuous Proof of Storage Cryptographically Verified across Cluster!
[*] Executing on-chain storage reward repayment from Publisher to Provider...

=== FINAL CONSOLIDATED FINANCIAL LEDGER AUDIT ===
  💳 Client Wallet Balance   : 9994.999980344680397312 ETH
  💳 Publisher Wallet Balance: 9999.497098972427099970 ETH
  💳 Provider Wallet Balance : 9999.499912311833025836 ETH (+1.5 ETH Total Earned!)
     - Lottery Ticket Settlement : +1.0 ETH
     - Storage Proof Repayment   : +0.5 ETH
  🏦 Escrow Channel Deposit  : 2.000000000000000000 ETH

======================================================================
🎉 ALL 13 CHECKPOINTS, PROOF OF STORAGE & REPAYMENT PASSED (100%)!   
======================================================================
```
