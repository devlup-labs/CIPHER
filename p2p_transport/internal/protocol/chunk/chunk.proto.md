# CIPHER Chunk Transport Protocol (v1.0.0)

## Overview
The Chunk Transport Protocol (`/cipher/chunk/1.0.0`) operates strictly over libp2p streams and is designed for high-throughput, resilient transfer of cryptographic chunks and capability manifests.

This protocol explicitly avoids any knowledge of local filesystems, filenames, or application-level routing. It deals exclusively with byte sequences identified by their SHA-256 hashes (`ChunkID` and `ContentID`).

## Message Framing & Envelope
To reduce libp2p stream setup overhead, the protocol multiplexes multiple requests and responses over a **single, long-lived bidirectional stream** between peers.

All messages share a unified, symmetric envelope delineated by an unsigned 32-bit (Varint) length prefix:

```text
Message Envelope
----------------
MessageVersion  uint16
MessageType     uint8
Payload         []byte
```

### Message Types (`uint8`)
- `0x01`: `REQUEST_MANIFEST`
- `0x02`: `MANIFEST`
- `0x03`: `REQUEST_CHUNK`
- `0x04`: `CHUNK`
- `0x05`: `ACK`
- `0x06`: `ERROR`

## Payload Definitions

### `REQUEST_MANIFEST`
Sent by a peer seeking to resolve a capability manifest.
- `ContentID [32]byte`: The hash identifying the requested Manifest.

### `MANIFEST`
Sent in response to a `REQUEST_MANIFEST`.
- `ContentID [32]byte`: The requested ID.
- `ManifestData []byte`: The serialized `Manifest` JSON/binary structure.

### `REQUEST_CHUNK`
Sent by a peer requesting a specific chunk from the host.
- `ChunkID [32]byte`: The requested chunk hash.

### `CHUNK`
The wire structure makes chunks entirely self-describing, ensuring the receiver has all necessary metadata to verify and decrypt.
```text
Chunk Header
------------
Version      uint16
ChunkID      [32]byte
Index        uint32
Offset       int64
PlainSize    uint32
CipherSize   uint32
Nonce        [24]byte

Payload
-------
Ciphertext   []byte (Length == CipherSize)
```

### `ACK`
Sent by the receiver after successfully validating and storing a chunk.
- `ChunkID [32]byte`: The chunk successfully written.
- `Status uint8`: Future-proof status flag (currently always `0` for Success; errors are sent via `ERROR`).

### `ERROR`
Sent by either peer to explicitly signal a protocol failure, rather than relying on stream closure.
- `ErrorCode uint8`
- `Message string`

**Explicit Error Codes**:
- `0x01`: `ERR_CONTENT_NOT_FOUND`
- `0x02`: `ERR_CHUNK_NOT_FOUND`
- `0x03`: `ERR_INVALID_MANIFEST`
- `0x04`: `ERR_PERMISSION_DENIED`
- `0x05`: `ERR_INTERNAL`
- `0x06`: `ERR_INTEGRITY_MISMATCH` (Sent if a received chunk fails SHA-256 verification)

## State Transitions & Workflows

### Standard Download Flow (Single Stream)
1. **Peer A** dials **Peer B** and opens `/cipher/chunk/1.0.0`.
2. **Peer A** sends `REQUEST_MANIFEST(ContentID)`.
3. **Peer B** responds with `MANIFEST(ManifestData)`.
4. **Peer A** parses the `Manifest` and enters the Download phase.
5. Over the **same stream**, **Peer A** pipelines `REQUEST_CHUNK(ChunkID)` for each required chunk.
6. **Peer B** reads chunks from its local `ContentEngine` and streams back `CHUNK(...)` envelopes.
7. As **Peer A** receives each `CHUNK`:
   - It hashes the `Ciphertext` to **verify** it matches the `ChunkID`.
   - If verification passes, it calls `Store.PutChunk()`.
   - It sends `ACK(ChunkID)` back over the stream.
   - If verification fails, it drops the chunk and sends `ERROR(ERR_INTEGRITY_MISMATCH)`.
8. Once all chunks are verified and stored, the download resolves.
9. (Application level decides when/if to `Reassemble()` the content independently of the download process).

### Testing Guarantees
- **Interrupted Transfer**: Disconnecting mid-transfer must cleanly abort the stream. Reconnecting must safely resume from the next missing chunk without corrupting the store.
- **Invalid Peer Request**: Requesting a non-existent `ContentID` must return `ERR_CONTENT_NOT_FOUND` cleanly over the stream.
