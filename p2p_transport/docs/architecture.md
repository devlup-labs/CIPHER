# CIPHER P2P Architecture

## Overview
The CIPHER project is built on top of [libp2p](https://libp2p.io/), utilizing a modular approach to separate concerns across the transport, identity, cryptography, and protocol layers.

## Modules

### `cmd/` (Entrypoints)
- **peer**: The standard client node participating in the network.
- **relay**: A specialized node designed to relay traffic between peers that may be behind NATs or firewalls.

### `internal/` (Core Logic)
- **transport**: Manages the initialization of libp2p hosts, connection establishment, and multiplexing.
- **identity**: Handles peer ID generation, key management, and cryptographic identities. It includes a persistent identity system that ensures a node's `PeerID` remains constant across restarts by storing an Ed25519 private key in the user's OS-level configuration directory (`~/.config/cipher/`, `Library/Application Support/CIPHER/`, or `AppData/Roaming/CIPHER/`).
- **crypto**: Provides standard cryptographic primitives for the broader application.
- **protocol**: Defines the custom network protocol IDs and handling used by CIPHER, currently including `/cipher/chunk/1.0.0` for direct, content-addressed P2P communication.
- **content**: The Content Engine Foundation. A completely decoupled pipeline that manages data ingestion above the transport layer. It contains:
  - **chunker**: Splits files into variable-sized or fixed-sized chunks based on engine config.
  - **crypto**: Encrypts and decrypts chunks independently using XChaCha20-Poly1305.
  - **verifier**: Hashes and verifies chunk/file integrity using SHA-256.
  - **manifest**: Manages immutable content capabilities (Chunk IDs, Hashes, Descriptors) decoupled from decryption rights.
  - **storage**: Defines `ChunkSource` and `ChunkSink` interfaces. Currently implemented using local, content-addressed files.
  - **engine**: The coordinator that wires the pipeline together (ingest and reassembly).
- **transfer**: High-level orchestrator of network downloads with strict separation of concerns:
  - **manager**: Owns transfer lifecycle, session state (resume/retry), and progress tracking without blocking.
  - **scheduler**: Owns work distribution. It maintains a lock-free queue and spins up a worker pool to fetch chunks concurrently from multiple source peers.

### Content Engine Data Flow

The Content Engine completely decouples the application's data ingestion from the transport layer. It operates as an independent pipeline designed for a future decentralized encrypted CDN.

```mermaid
graph TD
    App[Application] -->|Raw File/Stream| Chunker
    
    subgraph Content Engine
        Chunker[Chunker<br/>Splits into 256KB chunks] --> Crypto
        Crypto[Crypto<br/>Encrypts via XChaCha20-Poly1305] --> Verifier
        Verifier[Verifier<br/>Hashes ciphertext via SHA-256] --> Storage
        Storage[(ContentStore<br/>Local FS-backed CAS)]
        
        ManifestBuilder[Manifest Builder]
        Verifier -.->|Generates Chunk IDs| ManifestBuilder
    end
    
    ManifestBuilder -->|Outputs| Manifest[Immutable Manifest]
```

#### 1. Core Data Structures
- **Chunk & Headers**: Data is explicitly split into `ChunkHeader` and `Data` payload. This ensures headers (containing metadata like `Version`, `PlainSize`, and `CipherSize`) can be evaluated independently of the encrypted payload during network transmission.
- **Strong Typing**: The engine uses strict `[32]byte` types for `ChunkID` and `ContentID` to avoid string-encoding bugs, maximize comparison speed, and reduce heap allocations.

#### 2. Cryptography & Verification
- **Chunk Encryption (XChaCha20-Poly1305)**: Every chunk is encrypted independently. The engine generates secure 192-bit (24-byte) nonces automatically. Because chunks are encrypted independently, the engine natively supports out-of-order decryption and random access required by swarming protocols.
- **Content-Addressing**: Chunks are identified strictly by the SHA-256 hash of their **ciphertext**. 

#### 3. Content-Addressed Storage (CAS)
- **Decoupled Sinks**: The storage layer is abstracted behind `ChunkSource` and `ChunkSink` interfaces, ensuring the engine never touches the filesystem directly.
- **Sharded Local Storage**: The current `FSStore` implementation shards chunks using their hex-encoded hash prefixes (e.g., `store/ab/cd/abcdef123...`). This architecture mimics Git/IPFS to prevent filesystem degradation and inode limits when millions of chunks are persisted.

#### 4. Immutable Manifests
The `manifest` module generates a cryptographic capability file after ingestion. It intentionally decouples the **content description** (the ordered `ChunkIDs` and tree root) from the **decryption rights** (the content key). This permits the system to distribute the manifest publicly for swarming while restricting the decryption key to authorized users.

#### 5. Local Session Management & Multi-Device Swarming
Instead of requiring servers to maintain download states, CIPHER utilizes a strictly **client-side session architecture** for resume, recovery, and multi-device swarming. The network transfer pipeline has been successfully proven to parallelize chunk requests across public relay networks and direct hole-punched paths across different physical devices.

The network transfer pipeline follows a strict hierarchy:

```text
Application (CLI, Gateway)
        │
        ▼
TransferManager (Session, Progress, Retry Policy)
        │
        ▼
Scheduler (Queue, Source Maps)
        │
        ▼
Workers (Concurrent Fetching)
        │
        ▼
Chunk Protocol (Stateless Network Comms)
        │
        ▼
Content Engine (Verification & Storage)
```

The `TransferManager` tracks progress using a boolean bitset and persists it locally (e.g., `sessions/<ContentID>.json`). Upon restart or failure, missing chunks are automatically queued into the `Scheduler`, which assigns them to available `Worker`s connecting to seed peers. The underlying `/cipher/chunk/1.0.0` protocol remains fully stateless and completely decoupled from orchestration logic.

## Network Topology
The network utilizes a hybrid peer-to-peer topology where standard peers connect to one another directly if possible, or fallback to utilizing `relay` nodes for NAT traversal and connectivity routing.

### Transport Abstraction
To ensure that the application logic remains decoupled from the underlying transport state, CIPHER implements a `Transport` abstraction layer (`internal/transport/stream.go`). The application simply calls `Transport.Connect()` and `Transport.OpenStream()`. The Transport layer delegates the connection selection to libp2p, which autonomously manages background transport upgrades (e.g., Hole Punching) and provides the application with the best available path.

### Relay Connectivity & DCUtR Hole Punching (Circuit v2)
CIPHER leverages go-libp2p's `circuitv2` protocol for initial relaying and NAT traversal coordination. By default, `circuitv2` establishes **limited** (transient) connections.

While the application can fall back to communicating over these limited relay connections using `network.WithAllowLimitedConn`, CIPHER employs **DCUtR (Direct Connection Upgrade through Relay)**. DCUtR opportunistically runs in the background upon the establishment of a relay circuit. It coordinates a UDP/TCP hole punch, and upon success, automatically establishes a direct connection. Subsequent application streams naturally prefer this direct, high-throughput path over the temporary relay circuit.

#### Hole Punching State Machine
The sequence below illustrates how an application immediately connects via a relay, whilst DCUtR transparently upgrades the transport layer in the background.

```mermaid
sequenceDiagram
    participant App as Application
    participant Transport
    participant Libp2p
    participant Relay as Circuit v2 Relay
    participant Remote as Target Peer (NAT)

    App->>Transport: Connect(Peer)
    Transport->>Libp2p: Connect(Relay Multiaddr)
    Libp2p->>Relay: Reserve & Connect
    Relay->>Remote: Route Connection
    Libp2p-->>Transport: Connected (Limited)
    Transport-->>App: OK

    App->>Transport: OpenStream(Peer)
    Transport->>Libp2p: NewStream(WithAllowLimitedConn)
    Libp2p-->>App: Stream (via Relay)
    
    note over Libp2p, Remote: DCUtR Background Process
    Libp2p->>Remote: StartHolePunchEvt
    Libp2p->>Remote: HolePunchAttemptEvt (Simultaneous Dial)
    Remote-->>Libp2p: Direct Connection Established!
    Libp2p->>Libp2p: EndHolePunchEvt (Success)
    
    note right of App: App unaware of transport shift
    App->>Transport: OpenStream(Peer)
    Transport->>Libp2p: NewStream()
    Libp2p-->>App: Stream (via Direct Transport)
```

```mermaid
graph TD
    subgraph Direct Connection
        PeerA["Peer A (Listener)"] <-->|"/cipher/chunk/1.0.0"| PeerB["Peer B (Dialer)"]
    end
    subgraph Relayed Connection
        PeerC["Peer C (NAT)"] -->|Circuit v2 Reservation| Relay["Relay Node"]
        Relay -->|Relayed Stream| PeerD["Peer D (NAT)"]
    end
```
