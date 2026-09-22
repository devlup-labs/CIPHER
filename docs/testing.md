# CIPHER Testing Strategy & Quick Start Guide

This document provides a modern, fast, and structured guide to testing the CIPHER decentralized content delivery network. It covers unit testing, automated end-to-end role validation, multi-provider replication, physical multi-laptop testing, relay/DCUtR NAT traversal, and manual verification workflows.

---

## ⚡ Quick Start: Fast Automated Testing

All tests can be run locally with CGO disabled:

```bash
# 1. Run all unit and robustness tests (~3s)
CGO_ENABLED=0 go test -count=1 ./test/robustness/... ./internal/...

# 2. [NEW - Phase 4] Run Remote Ingestion & Multi-Provider Replication Test (3 Providers, R=2, Fault Kill)
bash test_remote_push.sh

# 3. Run automated End-to-End Role Separation Test (Publisher, Provider, Client, Bootstrap)
bash test_roles.sh

# 4. Run Single-Provider Independence & Persistence Restart Test
bash test_provider_independent.sh

# 5. Run Legacy Peer A <-> Peer B transfer test
bash test_transfer.sh
```

---

## 🏗️ The 5 CIPHER Roles Under Test

| Role | Binary Path | Responsibilities |
| :--- | :--- | :--- |
| **Publisher** | `cmd/publisher` | Ingests source file, chunks & encrypts via XChaCha20, creates manifest, plans placement, pushes to remote providers via `/cipher/push/1.0.0` with replication $R$, and exits. |
| **Provider** | `cmd/provider` | Standalone daemon hosting Content-Addressed Storage (CAS). Accepts uploads (`/cipher/push/1.0.0`), announces completed batches to DHT (`discovery.Provide`), serves chunks (`/cipher/chunk/1.0.0`), and republishes periodically. |
| **Client** | `cmd/client` | Discovers candidate providers via DHT (or direct dial `-d`), resolves manifest, swarms chunks concurrently via worker pool, verifies ciphertext hashes, decrypts out-of-order, and reassembles payload. |
| **Bootstrap** | `cmd/bootstrap` | Kademlia DHT bootstrap routing node for decentralized provider discovery. |
| **Relay** | `cmd/relay` | Circuit v2 Relay node for NAT traversal and DCUtR hole punching coordination. |

---

## 🧪 Core Test Scenarios

### Scenario 1: Remote Ingestion & Multi-Provider Replication (Phase 4)
**Objective**: Prove that the Publisher can push content to $N$ independent remote providers across the network with Replication Factor $R$, verify the replication invariant, and shut down. Then prove that clients can swarm-download the file even if a provider node is terminated midway.

```text
[Step 1] Bootstrap Node starts on port 48001
[Step 2] 3 Providers start with isolated directories (store_p1, store_p2, store_p3)
[Step 3] Publisher ingests 2 MB test file into an isolated store_pub
[Step 4] Publisher pushes chunks to Provider 1, 2, 3 with Replication R=2:
         - Circular placement assigns each chunk to 2 distinct providers
         - GlobalReplicaTracker confirms all 64 chunks reached >= 2 replicas
         - Providers announce ContentID to DHT on batch completion
         - Publisher exits cleanly (seed=false)
[Step 5] Client 1 discovers providers via DHT, swarms chunks from all 3 providers,
         decrypts, and verifies SHA-256 hash matches 100%
[Step 6] Fault Tolerance: Provider 1 is killed (kill $PROV1_PID)
[Step 7] Client 2 downloads file from surviving Providers 2 & 3 via DHT discovery
         and reassembles with 100% hash equality
```

**Automated Command**:
```bash
bash test_remote_push.sh
```

---

### Scenario 2: Independent Provider Lifecycle (Publisher Offline)
**Objective**: Prove that once content is published to a Provider, the **Publisher is completely removed from the retrieval path**.

```text
[Step 1] Bootstrap Node starts (DHT mesh coordinator)
[Step 2] Publisher ingests 2 MB test file into Provider CAS store
[Step 3] Publisher process is completely KILLED (verified 100% offline)
[Step 4] Provider starts independently with existing CAS store & registers with DHT
[Step 5] Client queries DHT solely with ContentID + Key + Bootstrap address
         - Discovers Provider on DHT (no direct peer address passed)
         - Resolves Manifest from Provider over /cipher/chunk/1.0.0
         - Downloads all chunks in parallel & decrypts payload
         - Validates SHA-256 checksum equality (100% match)
```

**Automated Command**:
```bash
bash test_provider_independent.sh
```

---

### Scenario 3: Provider Persistence & Restart
**Objective**: Prove that a Provider can be stopped, restarted, and immediately re-announce its local CAS store without data corruption or manifest loss.

```text
[Step 1] Stop the active Provider (kill $PROV_PID)
[Step 2] Restart Provider pointing to existing store directory (-store ./provider_store)
[Step 3] Provider startup automatically triggers discovery.StartRepublisher (re-announcing all CIDs)
[Step 4] A new Client queries DHT, discovers restarted Provider, fetches chunks, and reassembles
```

*This is automatically validated as Step 6 in `test_provider_independent.sh`.*

---

### Scenario 4: Multi-Provider Swarming (Candidate Provider Model)
**Objective**: Prove that the Client's `TransferManager` and `Scheduler` parallelize chunk downloads across multiple independent providers, gracefully handling partial chunk assignments.

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
                    └───►   Client  ◄───┘
```

#### Candidate Miss Semantics:
When a client requests a chunk from Provider A that Provider A does not host, Provider A returns `ErrChunkNotFound (0x02)`. The client scheduler identifies this as a **candidate miss**, marks Provider A in `MissedPeers` for that task, and immediately routes the task to Provider B without incrementing failure penalties.

---

### Scenario 5: Public Relay & DCUtR Hole Punching (NAT Traversal)
**Objective**: Connect across NATs/firewalls via a public `circuitv2` relay and verify automatic, seamless upgrade to direct TCP/UDP sockets via DCUtR.

1. **Start Public Relay (on cloud VM / Azure / Ubuntu)**:
   ```bash
   go run cmd/relay/main.go
   ```
   *Copy the relay multiaddress:* `/ip4/<PUBLIC_IP>/tcp/4001/p2p/<RELAY_PEER_ID>`
2. **Start Provider behind NAT**:
   ```bash
   go run cmd/provider/main.go -p 4001 -store ./provider_store -relay "/ip4/<PUBLIC_IP>/tcp/4001/p2p/<RELAY_ID>" -bootstrap "<BOOTSTRAP_ADDR>"
   ```
3. **Start Client on separate machine / network**:
   ```bash
   go run cmd/client/main.go -relay "/ip4/<PUBLIC_IP>/tcp/4001/p2p/<RELAY_ID>" -bootstrap "<BOOTSTRAP_ADDR>" -fetch "<CID>" -key "<KEY>" -out output.mp4
   ```
4. **Verify DCUtR Upgrade**:
   - Initial connection is established over the limited relay circuit.
   - DCUtR executes simultaneous hole punch in the background (`[DCUtR] Hole Punch Event: StartHolePunch` $\rightarrow$ `EndHolePunch`).
   - Transfer shifts to high-throughput direct socket (`Path: Direct`).

---

## 💻 Hardware & Multi-Laptop Testing Guide

For genuine real-world validation with physical machines:

### Option A: The 2-Laptop Setup (Minimum Hardware)
* **Laptop 1** (e.g. `192.168.1.50` on Home Wi-Fi):
  - Terminal 1: `./bin/bootstrap -p 48001 -ws-port 0`
  - Terminal 2: `./bin/provider -p 48010 -ws-port 0 -store ./store_p1 -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>"`
  - Terminal 3: `./bin/provider -p 48020 -ws-port 0 -store ./store_p2 -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>"`
* **Laptop 2** (e.g. `192.168.1.80` or Mobile Hotspot):
  - Push content:
    ```bash
    ./bin/publisher -file test.mp4 -push -replication 2 \
      -providers "/ip4/192.168.1.50/tcp/48010/p2p/<P1_ID>,/ip4/192.168.1.50/tcp/48020/p2p/<P2_ID>" \
      -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>" -seed=false
    ```
  - Retrieve content via DHT:
    ```bash
    ./bin/client -bootstrap "/ip4/192.168.1.50/tcp/48001/p2p/<BOOT_ID>" \
      -fetch "<CONTENT_ID>" -key "<KEY>" -out downloaded.mp4
    ```

### Option B: The 3-Machine Setup (VPS + 2 NAT Laptops)
* **Cloud VM / VPS** (Public IP `20.x.x.x`): Runs Bootstrap (`:4003`) and Relay (`:4001`).
* **Laptop A** (Home Wi-Fi): Runs Providers 1 & 2 pointing to VPS Relay and Bootstrap.
* **Laptop B** (Mobile Hotspot): Runs Publisher & Client. Traverses NAT via DCUtR hole punching to download directly from Laptop A.

---

## 📊 Unit & Robustness Test Suites

```bash
# Push wire protocol suite (framing, limits, serializers)
CGO_ENABLED=0 go test -v ./internal/protocol/push/...

# Distribution & Replication suite (circular planner, replica tracker)
CGO_ENABLED=0 go test -v ./internal/distribution/...

# Content Engine suite (chunking, ChaCha20, digests, CAS)
CGO_ENABLED=0 go test -v ./internal/content/...

# Chunk protocol suite (wire framing, handlers, ACKs)
CGO_ENABLED=0 go test -v ./internal/protocol/chunk/...

# 1000-iteration randomized robustness gauntlet
CGO_ENABLED=0 go test -v ./test/robustness/...
```

---

## 🛠️ CLI Quick Reference

```bash
# Push file to remote providers with R=2 replication and exit:
go run cmd/publisher/main.go -file <file> -push -replication 2 -providers "<P1_ADDR>,<P2_ADDR>" -bootstrap <boot_addr> -seed=false

# Provider: Start storage node with push enabled:
go run cmd/provider/main.go -p 4001 -ws-port 4002 -store <store_dir> -bootstrap <boot_addr> -allow-push=true

# Client: Fetch content via DHT swarming and reassemble:
go run cmd/client/main.go -bootstrap <boot_addr> -fetch <CID> -key <KEY_HEX> -out <output_path>

# Client: Query active download sessions:
go run cmd/client/main.go -store <store_dir> -status

# Client: Cancel a download session:
go run cmd/client/main.go -store <store_dir> -cancel <CID>
```
