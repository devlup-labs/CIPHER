# CIPHER Windows Testing Guide (WSL2 / Git Bash)

This guide provides instructions for setting up, installing prerequisites, and executing the complete CIPHER test suite on **Windows 10 / Windows 11**.

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

## 1. Recommended Environment: WSL2 (Windows Subsystem for Linux 2)

The recommended environment for running CIPHER on Windows is **WSL2** running Ubuntu 22.04 LTS or Ubuntu 24.04 LTS. WSL2 provides native Linux kernel support, high-performance POSIX file I/O, and zero-configuration networking with localhost port forwarding.

### Step 1: Install WSL2 & Ubuntu
Open **PowerShell as Administrator** and run:
```powershell
wsl --install -d Ubuntu
```
Restart your computer if prompted. After rebooting, complete the initial Ubuntu username and password setup.

### Step 2: Configure Git Line Endings (CRLF vs LF)
To ensure shell scripts retain POSIX `LF` line endings on Windows:
```bash
git config --global core.autocrlf input
```

---

## 2. Installing Required Dependencies & External Resources

Inside your WSL2 Ubuntu terminal (or Git Bash), run the following installation commands:

### A. Install Build Tools & Utilities
```bash
sudo apt update
sudo apt install -y build-essential curl git lsof psmisc libssl-dev
```

### B. Install Go (v1.22 or higher)
```bash
# Option 1: Using Ubuntu apt
sudo apt install -y golang-go

# Option 2: Direct official binary installation (Recommended for latest version)
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
rm go1.22.5.linux-amd64.tar.gz

# Add Go to PATH in ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

### C. Install Foundry Toolchain (`forge`, `anvil`, `cast`)
Foundry provides the local EVM blockchain engine and smart contract compilation suite:
```bash
curl -L https://foundry.paradigm.xyz | bash
source ~/.bashrc
foundryup

# Verify Foundry binaries
forge --version
anvil --version
cast --version
```

---

## 3. Windows Permissions & Networking Considerations

### A. Windows Defender Firewall
When running the Anvil EVM node (`:8545`) or CIPHER libp2p nodes for the first time, Windows Defender Firewall may display a prompt asking:
> *"Windows Defender Firewall has blocked some features of this app"*

* **Action Required**: Check **Private networks** (e.g. home or work network) and click **Allow access**.

### B. WSL2 Localhost Port Forwarding
WSL2 automatically maps `127.0.0.1` inside WSL to `localhost` on your Windows host. All 10 internal CIPHER ports (`8545`, `4001`, `4003`, `4101-4104`, `4201`, `4301-4303`) bind cleanly to the loopback interface without requiring manual port redirection.

### C. File Execution Permissions
Ensure all shell scripts have the executable permission bit set:
```bash
chmod +x ./local_multiple_terminal_test.sh
chmod +x ./test_workflow.sh
chmod +x ./scripts/*.sh
chmod +x ./tests/e2e/*.sh
chmod +x ./payments/*.sh
```

---

## 4. Test Execution Workflows (Relative Paths)

All test commands should be executed from the repository root:

### 1. Complete 13-Checkpoint Architecture Test (Universal Single-Terminal Mode)
On Windows/WSL2, `local_multiple_terminal_test.sh` automatically defaults to single-terminal background orchestration:
```bash
# Automated non-stop CI mode (all 13 checkpoints)
./local_multiple_terminal_test.sh --single --auto

# Interactive mode (pauses between checkpoints for inspection)
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

# Availability escrow and slashing contracts (7 tests)
(cd availability/escrow-payment/escrow && forge test -v)
```

### 4. Go Unit & Cryptographic Test Suites
```bash
# All unit tests across network subsystem
go test ./network/...

# Adversarial EIP-712 payment ticket tests
go test -v ./network/protocol/chunk -run TestChunkProtocol_TicketPayment
```

### 5. Cross-Domain End-to-End Test Suites
```bash
# Live Anvil EVM P2P Transfer & On-Chain Settlement
./tests/e2e/payment_transfer_anvil_e2e.sh

# Remote Ingestion & Multi-Provider Replication (R=2, dead-node kill)
./tests/e2e/remote_push_e2e.sh

# Standalone Provider Independence & Crash-Recovery Test
./tests/e2e/provider_lifecycle_e2e.sh

# Dynamic Availability, Demand Spikes & Slashing Simulator
./scripts/run_simulation.sh
```

---

## 5. Troubleshooting & Windows-Specific FAQs

### Q1: `bash: ./local_multiple_terminal_test.sh: /usr/bin/env: bad interpreter: No such file or directory` / `\r: command not found`
* **Cause**: The file was cloned with Windows `CRLF` line endings instead of Unix `LF`.
* **Fix**: Run `dos2unix` or convert line endings with `sed`:
  ```bash
  sudo apt install -y dos2unix
  dos2unix local_multiple_terminal_test.sh test_workflow.sh tests/e2e/*.sh
  ```

### Q2: Port 8545 is already in use
* **Cause**: An orphaned Anvil process from a previous test run is still active.
* **Fix**: Use the built-in cross-platform `fuser` utility:
  ```bash
  fuser -k 8545/tcp
  ```

### Q3: How to view live logs of background nodes during testing?
All daemon logs are saved directly in the repository root. Open a separate WSL2 terminal tab and run:
```bash
tail -f provider1.log
tail -f publisher.log
tail -f consumer.log
```
