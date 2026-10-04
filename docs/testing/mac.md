# CIPHER macOS Testing Guide (Apple Silicon & Intel)

This guide provides instructions for setting up dependencies, configuring macOS security automation permissions, and executing the complete CIPHER test suite on macOS (Apple Silicon M1/M2/M3/M4 and Intel Macs).

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

## 1. Installing Required Dependencies via Homebrew

### Step 1: Install Homebrew (if not already installed)
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Step 2: Install Go & Essential Utilities
```bash
brew install go git curl

# Verify Go installation (must be 1.22+)
go version
```

### Step 3: Install Foundry Toolchain (`forge`, `anvil`, `cast`)
Foundry provides the local Anvil EVM blockchain node and smart contract compiler:
```bash
curl -L https://foundry.paradigm.xyz | bash
source ~/.zshrc    # or ~/.bash_profile
foundryup

# Verify Foundry installation
forge --version
anvil --version
cast --version
```

---

## 2. macOS Security & Automation Permissions

### A. Terminal.app Desktop Perimeter Tiling Mode
On macOS, [`local_multiple_terminal_test.sh`](../../local_multiple_terminal_test.sh) includes a window tiling engine that uses AppleScript (`osascript`) to dynamically launch and position 10 specialized terminal windows around the perimeter of your desktop:

```text
+---------------------+---------------------+---------------------+---------------------+
| [Slot 1] ANVIL EVM  | [Slot 2] RELAY V2   | [Slot 3] BOOTSTRAP  | [Slot 4] TIER-1 CORE|
| Port: 8545          | Port: 4001 (HOP)    | Port: 4003 (DHT)    | Port: 4101 (Direct) |
+---------------------+---------------------+---------------------+---------------------+
| [Slot 5] TIER-2 EDGE|                                           | [Slot 6] TIER-3 AUDIT
| Port: 4102          |           [ CENTER DESKTOP AREA ]         | Port: 4103          |
| NAT-Relayed Cache   |        (Main IDE / Controller Workspace)  | Availability Prover |
+---------------------+                                           +---------------------+
| [Slot 7] TIER-4 STBY| [Slot 8] PUBLISHER  | [Slot 9] CONSUMER 1 | [Slot 10] CONSUMER 2|
| Port: 4104          | Port: 4201          | Port: 4301          | Port: 4302 / 4303   |
+---------------------+---------------------+---------------------+---------------------+
```

### B. Required macOS Permission Authorization
When executing the multi-window orchestrator for the first time, macOS will display a security prompt:
> **"Terminal wants access to control Terminal"** or **"Terminal wants access to control System Events"**

1. Click **OK / Allow** on the modal dialog.
2. If the prompt was missed or blocked, enable it manually:
   - Open **System Settings ** -> **Privacy & Security** -> **Automation**.
   - Under **Terminal** (or your terminal emulator), toggle **System Events** and **Terminal** to **ON**.

### C. Single-Terminal Fallback Mode
If using a terminal that does not support AppleScript window tiling (such as iTerm2, Alacritty, Warp, or VS Code integrated terminal), run the orchestrator with the `--single` flag:
```bash
./local_multiple_terminal_test.sh --single
```

---

## 3. Test Execution Workflows (Relative Paths)

All test commands should be executed from the repository root:

### 1. Complete 13-Checkpoint Architecture Orchestrator
```bash
# Option 1: macOS 10-Window Desktop Tiled Layout (Automated non-stop execution)
./local_multiple_terminal_test.sh --auto

# Option 2: macOS 10-Window Desktop Tiled Layout (Interactive step-by-step with checkpoints)
./local_multiple_terminal_test.sh

# Option 3: Universal Single-Terminal Mode (Automated non-stop CI)
./local_multiple_terminal_test.sh --single --auto

# Option 4: Universal Single-Terminal Mode (Interactive step-by-step with checkpoints)
./local_multiple_terminal_test.sh --single
```

### 2. Master Verification Pipeline
```bash
./test_workflow.sh
```

### 3. Smart Contract Test Suites (Foundry Forge)
```bash
# Payments, lottery entropy, and settlement contracts (46 tests)
(cd payments && forge test -v)

# Availability escrow and proof dispute contracts (7 tests)
(cd availability/escrow-payment/escrow && forge test -v)
```

### 4. Go Unit & Cryptographic Test Suites
```bash
# Run all unit tests across the network subsystem
go test ./network/...

# Test EIP-712 Typed Data Hashing & Signature Recovery
go test -v ./network/payments

# Test Chunk Wire Protocol & Adversarial Ticket Rejection
go test -v ./network/protocol/chunk -run TestChunkProtocol_TicketPayment
```

### 5. Cross-Domain End-to-End System Tests
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

# Multi-Window Dynamic Simulator (macOS Terminal.app)
./scripts/run_simulation.sh --multi-window
```

---

## 4. Troubleshooting & macOS-Specific FAQs

### Q1: `osascript: execution error: Not authorized to send Apple events to Terminal`
* **Fix**: Grant Automation permissions in **System Settings -> Privacy & Security -> Automation -> Terminal -> System Events (ON)**, or run with `./local_multiple_terminal_test.sh --single`.

### Q2: Port 8545 or 4001 already in use
* **Fix**: Terminate the listening PID using `lsof`:
  ```bash
  kill -9 $(lsof -ti tcp:8545) 2>/dev/null || true
  ```

### Q3: Window tiling positions on multiple monitors / 4K external displays
* **Behavior**: The script queries desktop bounds via `osascript` Finder geometry and scales the 10 slots proportionally. If using multiple displays, the windows tile across the primary display.
