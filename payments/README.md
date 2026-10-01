# CIPHER Payments & Economic Settlement Engine

This directory contains the smart contract layer and economic infrastructure for **CIPHER**, a trust-minimized decentralized Content Delivery Network. 

---

## 1. Overview & Economic Model

In a decentralized CDN, content providers must be economically compensated for caching and serving data chunks without requiring expensive, high-latency on-chain transactions for every single transfer.

CIPHER solves this with **Probabilistic Micropayments (Lottery Tickets)** combined with unidirectional **Payment Channels**:

1. **Unidirectional Payment Channels**: Consumers lock funds (ETH) into `PaymentChannel.sol` dedicated to serving providers.
2. **Zero-Gas Off-Chain Tickets**: For every chunk received via `/cipher/chunk/1.0.0`, the consumer signs an EIP-712 typed `RoundTicket` with small face value and probability $p \in (0, 1]$.
3. **Probabilistic Value**: The expected payment per chunk is $\mathbb{E}[V] = \text{FaceValue} \times \text{TicketProb}$. Consumers pay micro-amounts without gas fees.
4. **Batched On-Chain Settlement**: Providers collect tickets off-chain. Only tickets that satisfy the winning lottery condition against `CommitRevealEntropy.sol` are submitted to `SettlementEngine.sol` for payout from the consumer's channel.

---

## 2. Core Smart Contracts (`src/`)

```text
payments/src/
├── core/
│   ├── PaymentChannel.sol        # Channel escrow, deposits, and withdrawal management
│   ├── SettlementEngine.sol      # EIP-712 ticket verification, lottery math, & provider payouts
│   ├── ProviderRegistry.sol      # Staked provider registry & node identity tracking
│   └── ChunkDisputeResolver.sol  # Merkle tree dispute arbitration for corrupt/withheld chunks
├── randomness/
│   ├── CommitRevealEntropy.sol   # Unbiasable round-based commit-reveal entropy source
│   └── RaffleMath.sol            # Cryptographic lottery modulo arithmetic
├── interfaces/                   # Internal protocol interfaces
└── config/                       # Network parameter configurations
```

### Contract Breakdown:

| Contract | Purpose | Key Functions |
| :--- | :--- | :--- |
| **`PaymentChannel.sol`** | Manages consumer deposits in escrow. Protects funds from double-spending while guaranteeing provider solvency. | `openChannel(provider)`, `deposit(channelId)`, `settleDrawdown(channelId, amount, recipient)` |
| **`SettlementEngine.sol`** | The central clearinghouse. Verifies EIP-712 signatures, checks round entropy, and draws from `PaymentChannel` to pay providers. | `commitRound(roundId, tau, faceValue, hash)`, `claimTicket(ticket, signature)`, `batchSettle(claims[])` |
| **`CommitRevealEntropy.sol`** | Provides verifiable, unbiasable on-chain randomness so neither consumer nor provider can pre-determine winning tickets. | `commit(roundId, commitment)`, `reveal(roundId, secret)`, `getRoundEntropy(roundId)` |
| **`ProviderRegistry.sol`** | Registry for storage nodes. Enforces bonded collateral stakes to protect against Sybil attacks. | `registerProvider{value: stake}()`, `slashProvider(provider, amount)` |
| **`ChunkDisputeResolver.sol`** | Validates Merkle dispute proofs when a provider fails to serve promised data or delivers corrupt bytes. | `disputeChunk(channelId, chunkIndex, merkleProof)` |

---

## 3. Cryptographic Specification: EIP-712 Typed Tickets

Tickets exchanged over the P2P wire protocol are structured and hashed in accordance with **EIP-712**:

### EIP-712 Type Definition:
```solidity
struct RoundTicket {
    bytes32 channelId;
    bytes32 roundId;
    uint64 nonce;
    uint256 faceValue;
    uint256 ticketProb;
}
```

### Domain Separator:
* **Name**: `CIPHER Payment Protocol`
* **Version**: `1.0.0`
* **ChainID**: Dynamic (e.g., `31337` for Anvil, `1` for Mainnet)
* **Verifying Contract**: Address of deployed `SettlementEngine`

### Verification Pipeline:
1. Provider receives `MsgTicket (0x07)` over the chunk stream.
2. Provider computes EIP-712 hash:
   $$\text{Digest} = \text{keccak256}("\backslash x19\backslash x01" \parallel \text{DomainSeparator} \parallel \text{hashStruct}(\text{Ticket}))$$
3. Recover signer address via `ecrecover(Digest, v, r, s)`.
4. Validate signer matches consumer channel depositor and channel balance $\ge \text{FaceValue}$.

---

## 4. Go Ethereum Integration (`abigen`)

Go bindings are auto-generated from compiled Solidity artifacts into isolated subpackages to prevent struct name collisions:

```text
network/payments/
├── bindings/
│   ├── paymentchannel/     # Go bindings for PaymentChannel.sol
│   ├── settlementengine/   # Go bindings for SettlementEngine.sol
│   ├── commitreveal/       # Go bindings for CommitRevealEntropy.sol
│   └── providerregistry/   # Go bindings for ProviderRegistry.sol
├── client.go               # High-level ethclient wrapper & channel querying
├── signer.go               # Go EIP-712 hashing & Secp256k1 signing
└── signer_test.go          # Cryptographic unit tests verifying Solidity compatibility
```

### Regenerating Bindings:
```bash
# Compile contracts
forge build

# Generate Go bindings via abigen
abigen --abi out/PaymentChannel.sol/PaymentChannel.json --bin out/PaymentChannel.sol/PaymentChannel.bin --pkg paymentchannel --out ../network/payments/bindings/paymentchannel/payment_channel.go
abigen --abi out/SettlementEngine.sol/SettlementEngine.json --bin out/SettlementEngine.sol/SettlementEngine.bin --pkg settlementengine --out ../network/payments/bindings/settlementengine/settlement_engine.go
abigen --abi out/CommitRevealEntropy.sol/CommitRevealEntropy.json --bin out/CommitRevealEntropy.sol/CommitRevealEntropy.bin --pkg commitreveal --out ../network/payments/bindings/commitreveal/commit_reveal.go
abigen --abi out/ProviderRegistry.sol/ProviderRegistry.json --bin out/ProviderRegistry.sol/ProviderRegistry.bin --pkg providerregistry --out ../network/payments/bindings/providerregistry/provider_registry.go
```

---

## 5. Deployment & Testing Scripts (`script/`)

### Setup Script (`script/Step1_Setup.s.sol`)
Deploys core infrastructure to a running EVM node (Anvil or testnet) and deposits consumer funds:
- Deploys `CommitRevealEntropy`.
- Deploys `PaymentChannel`.
- Deploys `SettlementEngine` linked to the channel and entropy source.
- Registers provider in `ProviderRegistry`.
- Opens a payment channel funded with test ETH (e.g. 1 ETH).
- Outputs contract addresses and derived channel IDs.

### Settlement Script (`script/Step2_Settle.s.sol`)
Simulates or executes on-chain lottery round claims:
- Submits provider's accumulated winning tickets.
- Validates payout transfer from consumer channel balance to provider wallet.

---

## 6. How to Run Tests

### Run Foundry Solidity Test Suites:
```bash
forge test
```
*Executes all 46 contract unit, fuzz, and invariant test cases.*

### Run Go Unit & EIP-712 Cryptographic Tests:
```bash
go test -v ./network/payments/...
```

### Run Full Anvil End-to-End P2P Transfer & Settlement:
```bash
./tests/e2e/payment_transfer_anvil_e2e.sh
```

### Run Master CI/CD Workflow:
```bash
./test_workflow.sh
```
