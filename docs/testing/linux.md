# CIPHER Linux Testing Guide (Ubuntu, Debian, Fedora, Arch)

This guide provides instructions for setting up prerequisites, configuring environment parameters, and executing the complete CIPHER testing suite on Linux distributions (Ubuntu, Debian, Fedora, CentOS, Arch Linux).

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

## 1. Distribution-Specific Dependency Installation

Select your Linux distribution to install all required compilers, cryptographic libraries, and networking utilities:

### A. Ubuntu / Debian / Pop!_OS / Linux Mint
```bash
sudo apt update
sudo apt install -y golang-go build-essential curl git lsof psmisc libssl-dev pkg-config

# Install Foundry (forge, anvil, cast)
curl -L https://foundry.paradigm.xyz | bash
source ~/.bashrc
foundryup
```

### B. Fedora / RHEL / CentOS Stream / Rocky Linux
```bash
sudo dnf install -y golang make gcc gcc-c++ curl git lsof psmisc openssl-devel pkgconfig

# Install Foundry (forge, anvil, cast)
curl -L https://foundry.paradigm.xyz | bash
source ~/.bashrc
foundryup
```

### C. Arch Linux / Manjaro / EndeavourOS
```bash
sudo pacman -Syu --needed go base-devel curl git lsof psmisc openssl

# Install Foundry (forge, anvil, cast)
curl -L https://foundry.paradigm.xyz | bash
source ~/.bashrc
foundryup
```

### D. Verifying Installed Versions
```bash
go version        # Must be 1.22+
forge --version   # Foundry smart contract compiler
anvil --version   # Local EVM ledger
cast --version    # Ethereum RPC client
```

---

## 2. Linux Permissions & System Resource Configuration

### A. File Descriptor Limits (ulimit)
P2P swarms and concurrent chunk streamers open multiple TCP sockets across nodes. Verify that your system's open file limit is at least `4096`:
```bash
# Check current limit
ulimit -n

# Increase limit if below 4096 (temporary for current shell)
ulimit -n 4096

# To persist permanently, add to /etc/security/limits.conf:
# * soft nofile 65536
# * hard nofile 65536
```

### B. Non-Root Port Binding
All CIPHER ports are unprivileged TCP ports (`> 1024`):
- `8545` (Anvil EVM JSON-RPC)
- `4001` (Circuit Relay v2)
- `4003` (Kademlia DHT Bootstrap)
- `4101` - `4104` (4 Differentiated Storage Provider Tiers)
- `4201`, `4202` (Publisher Ingest & Availability Challenge)
- `4301`, `4302`, `4303` (Consumer Swarm, Fault Recovery, and Fraud Simulation)

No `sudo` or root privileges are required for running tests.

### C. Script Permissions
Ensure executable permissions are set on all shell harnesses:
```bash
chmod +x ./local_multiple_terminal_test.sh
chmod +x ./test_workflow.sh
chmod +x ./scripts/*.sh
chmod +x ./tests/e2e/*.sh
chmod +x ./payments/*.sh
```

---

## 3. Test Execution Workflows (Relative Paths)

All test commands should be executed from the repository root:

### 1. Complete 13-Checkpoint Architecture Runner (Automated CI & Headless Mode)
On Linux, `local_multiple_terminal_test.sh` automatically detects the non-macOS environment and executes all 10 nodes in unified background orchestrator mode:
```bash
# Automated non-stop mode (executes all 13 checkpoints from start to finish)
./local_multiple_terminal_test.sh --single --auto

# Step-by-step checkpoint mode (pauses after each milestone for inspection)
./local_multiple_terminal_test.sh --single
```

### 2. Master Verification Pipeline
Executes Solidity tests, Go cryptography tests, Availability engine tests, and live Anvil EVM settlement:
```bash
./test_workflow.sh
```

### 3. Smart Contract Verification (Foundry Forge)
```bash
# Payment channels, escrow, and EIP-712 settlement contracts (46 tests)
(cd payments && forge test -v)

# Availability escrow and proof dispute contracts (7 tests)
(cd availability/escrow-payment/escrow && forge test -v)
```

### 4. Go Unit & Cryptographic Test Suites
```bash
# Run all unit tests across the entire repository
go test ./network/...

# Test EIP-712 Typed Data Hashing & Signature Recovery
go test -v ./network/payments

# Test Chunk Wire Protocol & Adversarial Ticket Rejection
go test -v ./network/protocol/chunk -run TestChunkProtocol_TicketPayment
```

### 5. Cross-Domain End-to-End Test Suites
```bash
# Live Anvil EVM P2P Transfer & On-Chain Settlement
./tests/e2e/payment_transfer_anvil_e2e.sh

# Remote Ingestion & Multi-Provider Replication (R=2, dead-node kill)
./tests/e2e/remote_push_e2e.sh

# DHT Swarming & Provider Discovery
./tests/e2e/roles_e2e.sh

# Standalone Provider Independence & Crash-Recovery
./tests/e2e/provider_lifecycle_e2e.sh

# Dynamic Multi-Epoch Availability, Demand Spikes & Slashing Simulator
./scripts/run_simulation.sh
```

---

## 4. Live Inspection & Daemon Monitoring

During test execution, all background daemons write dedicated logs into the repository root:

```bash
# Monitor Provider 1 (Tier-1 Core Storage)
tail -f provider1.log

# Monitor Publisher Ingestion & Challenge Loop
tail -f publisher.log

# Monitor Consumer Chunk Swarming & Payment Ticket Streaming
tail -f consumer.log

# Inspect Kademlia DHT routing table and active provider index
./bin/dht-inspect -bootstrap "/ip4/127.0.0.1/tcp/4003/p2p/<BOOTSTRAP_PEER_ID>" -list-providers
```

---

## 5. Troubleshooting & Linux-Specific FAQs

### Q1: `bind: address already in use` (Port 8545, 4001, 4101, etc.)
* **Remedy**: Terminate any lingering processes using `fuser` or `lsof`:
  ```bash
  fuser -k 8545/tcp 4001/tcp 4003/tcp 4101/tcp 4102/tcp 4103/tcp 4104/tcp 4201/tcp 4301/tcp 4302/tcp 4303/tcp
  ```

### Q2: Loopback Firewall Block (UFW / iptables)
* **Remedy**: Ensure local loopback traffic (`127.0.0.1`) is permitted:
  ```bash
  sudo ufw allow in on lo
  sudo ufw allow out on lo
  ```
