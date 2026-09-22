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
- Payment, escrow, and settlement contracts belong in [`payments/contracts/`](payments/contracts/).

## Repository layout

```text
CIPHER/
├── docs/                   Protocol and architecture documentation
├── network/                Decentralized CDN networking
│   ├── cmd/                Network services (bootstrap, relay, peer)
│   ├── content/            Chunking, CAS storage, encryption, Merkle verification
│   ├── discovery/          Kademlia DHT routing & provider discovery
│   ├── distribution/       Multi-provider replication engine (circular stride)
│   ├── protocol/           Wire protocols (/cipher/chunk/1.0.0, /cipher/push/1.0.0)
│   ├── transfer/           Worker pool, transfer sessions, work-stealing scheduler
│   └── transport/          libp2p host, stream, and circuit v2 relay management
├── availability/           Availability mechanisms
│   └── contracts/          Availability-specific smart contracts
├── payments/               Economic settlement
│   └── contracts/          Payment and escrow smart contracts
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

### 1. Build Binaries
```bash
# Build all nodes and network services
go build -o bin/publisher ./nodes/publisher
go build -o bin/provider ./nodes/provider
go build -o bin/consumer ./nodes/consumer
go build -o bin/bootstrap ./network/cmd/bootstrap
go build -o bin/relay ./network/cmd/relay
```

### 2. Run Tests
```bash
# Run all unit tests
CGO_ENABLED=0 go test ./network/...

# Run End-to-End Multi-Provider Replication Test
./tests/e2e/remote_push_e2e.sh

# Run Role-Based DHT Swarming Test
./tests/e2e/roles_e2e.sh

# Run Provider Persistence & Independence Test
./tests/e2e/provider_lifecycle_e2e.sh
```

### 3. Running Nodes in Production / Staging
```bash
# Start DHT Bootstrap Node
./bin/bootstrap -p 4003

# Start Circuit Relay v2 Node
./bin/relay

# Start Provider Node
./bin/provider -p 4101 -store ./p1_store -bootstrap "<BOOTSTRAP_MULTIADDR>" -relay "<RELAY_MULTIADDR>" -allow-push=true

# Publish & Push Content Across Providers with Replication R=2
./bin/publisher -file ./sample.mp4 -bootstrap "<BOOTSTRAP_MULTIADDR>" -replication 2 -push

# Download Content as a Consumer
./bin/consumer -fetch "<CONTENT_ID>" -key "<KEY>" -out ./downloaded.mp4 -bootstrap "<BOOTSTRAP_MULTIADDR>"
```
