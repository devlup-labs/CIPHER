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
- [ ] Milestone 13: Provider Selection & Dynamic Scoring Engine (RTT/bandwidth-based scheduler prioritization)
- [ ] Milestone 14: Provider Reputation & Byzantine Fault Hardening (Corrupt chunk quarantine & peer blacklisting)
- [ ] Milestone 15: Tiered Caching & Dynamic Edge Replication (LRU RAM cache + demand-driven CDN replication)
- [ ] Milestone 16: Observability, Metrics & Telemetry Suite (Prometheus metrics & CLI dashboard)

---

### Command Quick-Start Reference

```bash
# Push distribution across remote providers with R=2 replication:
CGO_ENABLED=0 go run cmd/publisher/main.go \
  -file test_files/test.mp4 \
  -push \
  -providers "<P1_ADDR>,<P2_ADDR>,<P3_ADDR>" \
  -replication 2 \
  -bootstrap "<BOOTSTRAP_MULTIADDR>" \
  -seed=false

# Client swarm download via DHT:
CGO_ENABLED=0 go run cmd/client/main.go \
  -bootstrap "<BOOTSTRAP_MULTIADDR>" \
  -fetch "<CONTENT_ID>" \
  -key "<DECRYPTION_KEY>" \
  -out downloaded.mp4
```
