# CIPHER

CIPHER is a trust-minimized decentralized content delivery network (CDN). Publishers authenticate content, independent providers cache and serve it, and consumers discover providers and verify the content they receive.

CIPHER is designed for content delivery and caching, not permanent decentralized storage.

## Roles

- **Publisher** — divides content into chunks, commits to them with a Merkle root, and replicates across providers.
- **Provider** — advertises, caches, and serves publisher-authenticated content as an independent, permissionless node.
- **Consumer** — discovers suitable providers via DHT, fetches content, and verifies chunks using Merkle proofs and digests.

Content authentication proves that a chunk belongs to publisher-authenticated content. It does not, by itself, prove delivery or receipt.

## Development domains

CIPHER is divided into three independently owned domains:

- [`network/`](network/) handles peer connectivity, discovery, requests, transfer, demand metadata, cache announcements, and provider selection.
- [`availability/`](availability/) handles availability agreements, epochs, challenges, proof verification, and availability outcomes.
- [`payments/`](payments/) handles accounting, authorization, escrow, collateral, settlement, refunds, withdrawals, and payment disputes.

Smart contracts stay with the domain that owns their behavior:

- Availability-specific contracts belong in [`availability/contracts/`](availability/contracts/).
- Payment, escrow, and settlement contracts belong in [`payments/src/`](payments/src/).

## Repository layout

```text
CIPHER/
├── docs/                   Protocol and architecture documentation
├── network/                Decentralized CDN networking & payment bindings
│   ├── cmd/                Network services (bootstrap, relay, peer)
│   ├── content/            Chunking, CAS storage, encryption, Merkle verification
│   ├── discovery/          Kademlia DHT routing & provider discovery
│   ├── distribution/       Multi-provider replication engine (circular stride)
│   ├── identity/           Ed25519 (P2P) and Secp256k1 (Ethereum) key management
│   ├── payments/           Go contract bindings (abigen), EIP-712 signer, and client
│   ├── protocol/           Wire protocols (/cipher/chunk/1.0.0 + MsgTicket, /cipher/push/1.0.0)
│   ├── transfer/           Worker pool, transfer sessions, work-stealing scheduler
│   └── transport/          libp2p host, stream, and circuit v2 relay management
├── availability/           Availability mechanisms
│   └── contracts/          Availability-specific smart contracts
├── payments/               Economic settlement & smart contracts
│   ├── src/                Payment, escrow, lottery entropy, and settlement contracts
│   ├── script/             Deployment, setup, and settlement scripts
│   └── test/               Foundry Forge test suites
├── shared/                 Stable shared definitions and utilities
├── nodes/                  Publisher, provider, and consumer applications
│   ├── publisher/          Content ingestion and remote push distributor
│   ├── provider/           Decentralized storage node & chunk server
│   └── consumer/           Content discovery, swarming download & reassembly
├── integration/            Cross-domain composition and adapters
├── tests/                  Cross-module and end-to-end tests
│   ├── e2e/                Automated integration & lifecycle test scripts
│   └── robustness/         Fuzzing & boundary verification tests
├── scripts/                Development and operational utilities
├── config/                 Shared configuration templates
└── docker/                 Optional local container environment
```

## Quick Start

> [!IMPORTANT]
> **Notice for Local Testing:**
> Local multi-terminal testing, role simulations, and verification suites must be executed against the **`local`** branch of the repository:
> 👉 **[https://github.com/devlup-labs/CIPHER/tree/local](https://github.com/devlup-labs/CIPHER/tree/local)**
>
> ```bash
> git clone https://github.com/devlup-labs/CIPHER.git
> cd CIPHER
> git checkout local
> ```

### 1. Build Binaries
```bash
# Build all nodes and network services
go build -o bin/publisher ./nodes/publisher
go build -o bin/provider ./nodes/provider
go build -o bin/consumer ./nodes/consumer
go build -o bin/bootstrap ./network/cmd/bootstrap
go build -o bin/relay ./network/cmd/relay
```

## 🧪 Testing Locally

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

### Dedicated Platform Testing Guides
Detailed installation instructions, system permissions, external dependencies, and OS-specific run commands:

| Platform | Guide | Recommended Environment |
| :--- | :--- | :--- |
| **Unified** | [`docs/testing/general.md`](docs/testing/general.md) | Architecture, verification topology, global prerequisites & tools |
| **Windows** | [`docs/testing/windows.md`](docs/testing/windows.md) | WSL2 (Ubuntu 22.04/24.04 LTS) or Git Bash with native Foundry |
| **Linux** | [`docs/testing/linux.md`](docs/testing/linux.md) | Ubuntu / Debian / Fedora / Arch Linux (native POSIX background orchestration) |
| **macOS** | [`docs/testing/mac.md`](docs/testing/mac.md) | Apple Silicon & Intel (10-window desktop perimeter tiling or single-terminal) |

---

### Local Test Execution Modes

```bash
# 1. Complete 13-Checkpoint Architecture Orchestrator (10 Specialized Node Roles)
./local_multiple_terminal_test.sh --auto            # Option 1: macOS 10-window desktop tiled layout (Automated non-stop)
./local_multiple_terminal_test.sh                  # Option 2: macOS 10-window desktop tiled layout (Interactive step-by-step)
./local_multiple_terminal_test.sh --single --auto   # Option 3: Universal single-terminal mode (Automated non-stop CI - Linux, macOS, WSL2)
./local_multiple_terminal_test.sh --single          # Option 4: Universal single-terminal mode (Interactive step-by-step)

# 2. Master Verification Pipeline (Solidity + Unit + Adversarial Wire + Live Anvil Settlement)
./test_workflow.sh

# 3. Smart Contract Verification (Foundry Forge)
(cd payments && forge test)                                 # Payments, escrow, and EIP-712 settlement (46 tests)
(cd availability/escrow-payment/escrow && forge test)       # Availability escrow and slashing contracts (7 tests)

# 4. Go Unit, Cryptography & Adversarial Wire Protocol Suites
go test ./network/...                                       # Full network unit test suite
go test -v ./network/protocol/chunk -run TestChunkProtocol  # Adversarial wire streaming & ticket forge rejection

# 5. Cross-Domain End-to-End System Integration Tests
./tests/e2e/payment_transfer_anvil_e2e.sh                   # Live Anvil EVM P2P chunk transfer & settlement
./tests/e2e/remote_push_e2e.sh                              # Remote push & replication (R=2, dead-node kill)
./tests/e2e/roles_e2e.sh                                    # Role-based DHT provider discovery & swarm retrieval
./tests/e2e/provider_lifecycle_e2e.sh                       # Standalone provider crash-recovery & persistence
./scripts/run_simulation.sh                                 # Multi-epoch availability, demand spikes & slashing simulation
```

---

## 🚀 Running Nodes in Staging / Production

```bash
# 1. Start DHT Bootstrap Node
./bin/bootstrap -p 4003

# 2. Start Circuit Relay v2 Node (NAT Traversal Proxy)
./bin/relay

# 3. Start Storage Provider (with EVM ticket verification & CAS storage)
./bin/provider -p 4101 -store ./p1_store -bootstrap "<BOOTSTRAP_MULTIADDR>" \
  --eth-rpc "http://127.0.0.1:8545" --entropy-addr "<ENTROPY_CONTRACT_ADDR>"

# 4. Ingest, Encrypt & Push Content Across Providers with Replication R=2
./bin/publisher -file ./sample.mp4 -bootstrap "<BOOTSTRAP_MULTIADDR>" -replication 2 -push

# 5. Swarm Download Content as a Consumer (streaming EIP-712 payment tickets)
./bin/consumer -fetch "<CONTENT_ID>" -key "<KEY>" -out ./downloaded.mp4 -bootstrap "<BOOTSTRAP_MULTIADDR>" \
  --eth-rpc "http://127.0.0.1:8545" --eth-key "<CLIENT_ETH_PRIVKEY>" \
  --entropy-addr "<ENTROPY_CONTRACT_ADDR>" --provider-eth-addr "<PROVIDER_ETH_ADDR>"
```

---

## 📚 Documentation & Technical References

The complete architectural specifications, protocol designs, and test logs are maintained in the [`docs/`](docs/) directory:

### Architecture & System Specifications
* 🏛️ **[System Architecture & Security Model](docs/architecture.md)** — Core primitives, trust-minimization invariants, and domain ownership boundaries.
* 📦 **[Availability Subsystem & Proof Engine](docs/availability_integration.md)** — Ingestion, Merkle proof trees, availability epochs, challenges, and collateral slashing.
* 💳 **[Network & Economic Settlement Integration](docs/network_payments_integration.md)** — Dual identity (Ed25519 + Secp256k1), EIP-712 lottery tickets, and on-chain arbitration.
* 🌐 **[Circuit Relay v2 & NAT Deployment](docs/relay_deployment.md)** — HOP protocol reservations, hole-punching proxy configuration, and production topology.
* 🗺️ **[Development Roadmap & Milestones](docs/roadmap.md)** — Phased delivery plan and future protocol enhancements.

### Testing & Verification Guides
* 🎯 **[13-Step Local Testing & Architecture Execution Guide](docs/13_step_local_testing_execution.md)** — Checkpoint-by-checkpoint technical breakdown with live terminal outputs, logs, and cryptographic audits.
* 🧪 **[Master Testing Strategy & Quick Start](docs/testing.md)** — Complete test inventory, adversarial scenarios, and physical multi-laptop test setup.
* 🌐 **[General Testing Architecture & Matrix](docs/testing/general.md)** — Unified topology, prerequisite matrix, and repository-relative conventions.
* 🪟 **[Windows Testing Guide (WSL2 / Git Bash)](docs/testing/windows.md)** — WSL2 configuration, build tools, firewall permissions, and run commands.
* 🐧 **[Linux Testing Guide (Ubuntu, Fedora, Arch)](docs/testing/linux.md)** — Native Linux packages, `ulimit` tuning, non-root ports, and CI automation.
* 🍏 **[macOS Testing Guide (Apple Silicon & Intel)](docs/testing/mac.md)** — Homebrew setup, AppleScript window tiling, automation permissions, and single-terminal fallback.

