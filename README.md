# 🌐 CIPHER

> **CDNs are owned by a handful of companies. CIPHER is the alternative — a decentralized content delivery protocol where math replaces the middleman.**

CIPHER is a peer-to-peer content delivery network designed to bypass centralized intermediaries. It allows users to securely share, verify, and stream files directly between one another using advanced cryptographic integrity checks and automated network hole-punching.

Whether you're behind a restrictive NAT or simply want to share data without relying on a corporate tunnel, CIPHER ensures your data is delivered safely and authentically.

## Highlights

* **Decentralized by Design**: No central authority controls your data.
* **Zero-Config Networking**: Built-in Libp2p hole-punching (DCUtR) connects peers across isolated networks and firewalls when direct connectivity is possible, with Circuit Relay v2 as the relay path.
* **Cryptographic Integrity**: Keccak256-based Merkle trees guarantee that every byte received is exactly what was requested.
* **Secure Transport**: All file chunks are encrypted in transit using XChaCha20-Poly1305 symmetric authenticated encryption.
* **Multi-Protocol**: Supports connection establishment through QUIC, TCP, and Secure WebSockets.

## Usage

CIPHER works by spinning up a **Provider** to serve data and a **Client** to request it. A public Circuit Relay v2 service allows the peers to communicate across separate networks.

### 1. Start the Provider

The provider initializes the data, generates a cryptographic Merkle tree for integrity, connects to the public relay, and reserves a relay slot.

```bash
go run ./cmd/provider -relay /dns4/cipher-al5a.onrender.com/tcp/443/wss/p2p/12D3KooWHc1PMoMsmwnFw1yNDxSgQFhUQidAymtyi9UZC6yQ9zmV --verbose
```

A successful startup should show:

```text
Connected to relay
Successfully reserved slot on relay
Provider Peer ID: <PROVIDER_PEER_ID>
```

Also note the values printed during startup:

```text
Root: <MERKLE_ROOT_HASH>
Chunks: <CHUNK_COUNT>
```

Keep the provider running while the client downloads the file.

### 2. Start the Client

From a different network, run the client using the provider's relay circuit address.

Append:

```text
/p2p-circuit/p2p/<PROVIDER_PEER_ID>
```

to the relay address and provide the Merkle root and chunk count printed by the provider.

```bash
go run ./cmd/client \
  -provider /dns4/cipher-al5a.onrender.com/tcp/443/wss/p2p/12D3KooWHc1PMoMsmwnFw1yNDxSgQFhUQidAymtyi9UZC6yQ9zmV/p2p-circuit/p2p/<PROVIDER_PEER_ID> \
  -root <MERKLE_ROOT_HASH> \
  -chunks <CHUNK_COUNT> \
  --verbose
```

For the default test file and current provider identity, an example command is:

```bash
go run ./cmd/client -provider /dns4/cipher-al5a.onrender.com/tcp/443/wss/p2p/12D3KooWHc1PMoMsmwnFw1yNDxSgQFhUQidAymtyi9UZC6yQ9zmV/p2p-circuit/p2p/12D3KooWBhRgH4ggEm1iTuTznm9zQ1Jt2AEehyUyAAkvBWmbv4ER -root dd8f60d3e7a7b99c19e4561814e7d22c597814de25677b18d4bbc38aaa1aa940 -chunks 4 --verbose
```

> **Note:** The example client command is valid for the shown provider identity and test file. If the provider identity, source file, Merkle root, or chunk count changes, use the latest values printed by the provider.

The client connects through the relay circuit, requests encrypted chunks, receives the key reveal after the protocol exchange, verifies the cryptographic proofs, and reassembles the file locally.

The final output is:

```text
downloaded_file.txt
```

For deep diagnostic logs during transport, append the `--verbose` flag to either command.

## Installation

Currently, CIPHER is built from source. Ensure you have [Go](https://golang.org/doc/install) installed on your system.

Clone the repository and build the binaries:

```bash
# Clone the repository
git clone https://github.com/anshul090706/CIPHER.git
cd CIPHER

# Download dependencies
go mod download

# Build the client and provider executables
go build -o provider ./cmd/provider
go build -o client ./cmd/client
```

The client machine should clone the complete repository because the client executable depends on the internal CIPHER packages.

The client does not need the provider's original `test_file.txt`. After a successful transfer, the received and verified content is written to `downloaded_file.txt`.

## Relay Deployment

CIPHER now includes a deployable Circuit Relay v2 service under:

```text
cmd/relay
```

For a Render Web Service, use:

```text
Language: Go
Branch: main
Root Directory: Leave blank

Build Command:
go build -o cipher-relay ./cmd/relay

Start Command:
./cipher-relay
```

The deployed relay uses a persistent libp2p identity supplied through the `RELAY_PRIVATE_KEY` environment variable. This keeps the Relay Peer ID stable across restarts and redeployments.

Current public relay:

```text
Host:
cipher-al5a.onrender.com

Peer ID:
12D3KooWHc1PMoMsmwnFw1yNDxSgQFhUQidAymtyi9UZC6yQ9zmV
```

Full relay multiaddress:

```text
/dns4/cipher-al5a.onrender.com/tcp/443/wss/p2p/12D3KooWHc1PMoMsmwnFw1yNDxSgQFhUQidAymtyi9UZC6yQ9zmV
```

## Cross-Network Fixes

During end-to-end cross-network testing, the following issues were identified and fixed:

* **Broken relay dependency** — replaced the unavailable public relay dependency with a deployable Circuit Relay v2 service.
* **Relay reservation failure** — explicitly initialized the Circuit Relay v2 service so the HOP protocol is available for provider reservations.
* **mDNS interference** — disabled automatic mDNS discovery in the provider and client for deterministic explicit relay-path communication.
* **Relayed stream timeout** — allowed the CIPHER application stream to operate over libp2p limited relay connections.
* **Unstable relay identity** — added persistent relay identity support through the `RELAY_PRIVATE_KEY` environment variable so the Relay Peer ID remains stable across deployments.

The tested transfer flow is:

```text
Provider
    ↓
Connect to Relay
    ↓
Reserve Relay Slot
    ↓
Client connects through /p2p-circuit/
    ↓
ChunkRequest
    ↓
ChunkResponse
    ↓
LotteryTicket
    ↓
KeyReveal
    ↓
Decrypt Chunk
    ↓
Verify HResp
    ↓
Verify Merkle Proof
    ↓
downloaded_file.txt
```

## Verify the Transfer

On Windows PowerShell, compare the source and downloaded file hashes.

Provider:

```powershell
Get-FileHash test_file.txt -Algorithm SHA256
```

Client:

```powershell
Get-FileHash downloaded_file.txt -Algorithm SHA256
```

Both hashes should match.
