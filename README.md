# 🌐 CIPHER

> **CDNs are owned by a handful of companies. CIPHER is the alternative — a decentralized content delivery protocol where math replaces the middleman.**

CIPHER is a peer-to-peer content delivery network designed to bypass centralized intermediaries. It allows users to securely share, verify, and stream files directly between one another using advanced cryptographic integrity checks and automated network hole-punching.

Whether you're behind a restrictive NAT or simply want to share data without relying on a corporate tunnel, CIPHER ensures your data is delivered safely and authentically.

##  Highlights

*   **Decentralized by Design**: No central authority controls your data.
*   **Zero-Config Networking**: Built-in Libp2p hole-punching (DCUtR) effortlessly connects peers across isolated networks and firewalls.
*   **Cryptographic Integrity**: Keccak256-based Merkle trees guarantee that every byte received is exactly what was requested.
*   **Secure Transport**: All file chunks are encrypted in transit using robust XChaCha20-Poly1305 symmetric authenticated encryption.
*   **Multi-Protocol**: Automatically negotiates the best connection via QUIC, TCP, or Secure WebSockets.
* **Ethereum-ready Incentives**: Providers are registered with an Ethereum identity, with MVP support for provider registration, staking, and challenge handling. These components currently use in-memory state for local testing and are intended to be replaced by Ethereum smart contracts.

##  Usage

CIPHER works by spinning up a **Provider** to serve data, and a **Client** to request it. Here is a minimal example of how to securely transfer a file across two completely separate networks using a public relay:

**1. Start the Provider**
The provider will initialize the data, generate a cryptographic Merkle tree for integrity, and reserve a connection slot on a public relay.

```bash
go run ./cmd/provider -relay /dns4/relay-torrentium-3zok.onrender.com/tcp/443/wss/p2p/12D3KooWEBxhvkASAJtmdeKWiWWhdXCzwXEVvSMpjuY8YrDAi68Z
```
*Note the `Root` hash and the `Peer ID` outputted in the terminal.*

**2. Start the Client**
From a different network, run the client. Provide the provider's full relay address (append `/p2p-circuit/p2p/<PROVIDER_PEER_ID>`) and the Merkle root hash.

```bash
go run ./cmd/client \
  -provider /dns4/relay-torrentium-3zok.onrender.com/tcp/443/wss/p2p/12D3KooWEBxhvkASAJtmdeKWiWWhdXCzwXEVvSMpjuY8YrDAi68Z/p2p-circuit/p2p/<PROVIDER_PEER_ID> \
  -root <MERKLE_ROOT_HASH> \
  -chunks 4
```

The client will automatically negotiate a direct hole-punched connection if possible, stream the encrypted chunks, verify the cryptographic proof, and reassemble the file locally. 

*(For deep diagnostic logs during transport, simply append the `--verbose` flag to either command!)*

##  Installation

Currently, CIPHER is built from source. Ensure you have [Go](https://golang.org/doc/install) installed on your system.

Clone the repository and build the binaries:

```bash
# Clone the repository
git clone https://github.com/1amKhush/CIPHER.git
cd CIPHER

# Build the client and provider executables
go build -o provider ./cmd/provider
go build -o client ./cmd/client
```

##  Feedback & Contributing

We are building CIPHER to make a lasting impact on how data is distributed across the internet, and community feedback is everything. 

If you find this project interesting, have a feature request, or run into an issue, please don't hesitate to open an [Issue](https://github.com/1amKhush/CIPHER/issues) or start a discussion. Pull requests are always welcome! 
