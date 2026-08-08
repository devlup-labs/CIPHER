# CIPHER P2P Transport

This repository contains the P2P transport module for the CIPHER project.

### 🚀 Current Status (Milestone 10 Achieved!)
The CIPHER transport layer has successfully achieved **End-to-End Multi-Peer Swarming**. The system can now ingest files into a decoupled Content Engine (which chunks, encrypts via XChaCha20, and hashes via SHA-256), and dynamically download those encrypted chunks concurrently from multiple peers across the Internet. 

It natively supports NAT traversal by leveraging libp2p `circuitv2` public relays and automatically upgrades to high-speed direct TCP/UDP connections in the background via **DCUtR (Hole Punching)**.

## Structure

- `cmd/peer`: Entry point for standard peers.
- `cmd/relay`: Entry point for relay peers.
- `internal/transport`: libp2p host creation, connection management, and hole punching (DCUtR).
- `internal/content`: The Content Engine Foundation (Chunking, Encryption, Hashing, Manifests, Storage).
- `internal/identity`: Persistent Ed25519 identity generation.
- `internal/protocol`: Custom libp2p protocol definitions.
- `configs/`: Configuration files.
- `data/`: Data storage directories.

## Documentation

- [Architecture](docs/architecture.md): Overview of the project's modular design and libp2p network topology.
- [Testing Strategy](docs/testing.md): Instructions and guidelines for unit and integration testing.
- [Relay Deployment](docs/relay_deployment.md): Instructions on how to deploy a public relay node on Linux.

## Setup

```bash
# Build the project into the bin directory
go build -o bin/peer ./cmd/peer
go build -o bin/relay ./cmd/relay
go build -o bin/content-test ./cmd/content-test

# Run core transport network tests
go test ./...

# Run Content Engine robustness tests
go test -v ./test/robustness/...
```
