# CIPHER Roadmap

## Phase 1: Foundation
- [x] Milestone 1: Repository Setup (Project structure, libp2p basic integration)
- [x] Milestone 2: Persistent Peer Identity (Ed25519 identity generation and persistence)
- [x] Milestone 3: Direct P2P File Transfer (Custom protocol `/cipher/filetransfer/1.0.0` over basic streams)

## Phase 2: Enhanced Connectivity
- [x] Milestone 4: Bare Relay Service (Basic `circuitv2` relay node creation)
- [x] Milestone 5: Relay Connectivity (Routing application streams over `circuitv2` limited connections)
- [x] Milestone 6: Hole Punching & DCUtR (Upgrading relay connections to direct TCP/UDP via hole punching)

## Phase 3: Protocol Refinement
- [x] Milestone 7: Content Engine Foundation (Chunking, Cryptography, Content-Addressed Storage, Manifests)
- [x] Milestone 8: Content-Addressed Protocol & Integration
- [x] Milestone 9: Reliable Content Transfer (Session Management, Resume, Retry)
- [x] Milestone 10: Multi-peer Swarming & Chunk Scheduling *(Validated across multiple devices & NATs via public relays!)*

## Phase 4: Decentralization & Scaling
- [x] Milestone 11: Decentralized Discovery (Kademlia DHT server mode, CID generation, periodic Provider Republisher daemon)
- [x] Milestone 12: Remote Ingestion & Multi-Provider Replication Protocol (`/cipher/push/1.0.0`, circular placement planner, global replica invariant tracker)

## Phase 5: Economic Incentives & Settlement (Current)
- [x] Milestone 13: Dual Identity Management (Secp256k1 Ethereum keys alongside libp2p Ed25519 keys)
- [x] Milestone 14: Foundry Smart Contracts & Go Bindings (`PaymentChannel`, `SettlementEngine`, `CommitRevealEntropy`, `ProviderRegistry`)
- [x] Milestone 15: Streaming Micro-Payments Protocol (`MsgTicket 0x07` EIP-712 chunk-for-ticket verification)
- [x] Milestone 16: Live Anvil EVM Integration & Master Workflow (`./test_workflow.sh`)
- [ ] Milestone 17: Windowed Credit Buffering & Autonomous Background Settlement Daemons
- [ ] Milestone 18: On-Chain Merkle Dispute Arbitration (`ChunkDisputeResolver`) & Availability Slashing

---

### Command Quick-Start Reference

```bash
# Execute master end-to-end workflow (Solidity + Go Unit + Wire + Live Anvil Settlement):
./test_workflow.sh

# Push distribution across remote providers with R=2 replication:
go run nodes/publisher/main.go \
  -file test_files/test.mp4 \
  -push \
  -providers "<P1_ADDR>,<P2_ADDR>,<P3_ADDR>" \
  -replication 2 \
  -bootstrap "<BOOTSTRAP_MULTIADDR>"

# Consumer download via DHT with streaming EIP-712 payment tickets:
go run nodes/consumer/main.go \
  -bootstrap "<BOOTSTRAP_MULTIADDR>" \
  -fetch "<CONTENT_ID>" \
  -key "<DECRYPTION_KEY>" \
  -out downloaded.mp4 \
  --eth-rpc "http://127.0.0.1:8545" \
  --eth-key "<CLIENT_ETH_PRIVKEY>" \
  --entropy-addr "<ENTROPY_CONTRACT_ADDR>" \
  --provider-eth-addr "<PROVIDER_ETH_ADDR>"
```
