# Testing Strategy

The CIPHER project employs a comprehensive testing strategy to ensure reliability and correctness of the P2P networking components.

## Running Tests Locally

To run all unit tests for the project without CGO dependencies:

```bash
CGO_ENABLED=0 go test ./...
```

To run tests with verbose output:

```bash
CGO_ENABLED=0 go test -v ./...
```

## Testing Layers

1. **Unit Testing**:
   - Focuses on testing individual modules in isolation.
   - Example: Testing the `transport.NewNode` initialization ensures that a libp2p host can be correctly allocated and bound to a port.

2. **Integration Testing** (Planned):
   - Focuses on testing interactions between the `peer` and `relay` nodes.
   - Verifies end-to-end connectivity, stream multiplexing, and protocol correctness over local temporary networks.

3. **Manual Testing**:
   - **Persistent Identity**: To manually test that node identities persist across restarts:
     1. Build the peer binary: `go build -o bin/peer ./cmd/peer`
     2. Run the peer node: `./bin/peer`
     3. Note the outputted `Peer ID` and shut down the peer (Ctrl+C).
     4. Run the peer node again: `./bin/peer`
     5. Verify that the outputted `Peer ID` matches exactly with the previous run.

   - **Direct P2P Connection (Chunk Transport)**: To manually test the `/cipher/chunk/1.0.0` protocol:
     1. Build the peer binary: `go build -o bin/peer ./cmd/peer`
     2. Start **Peer A (Listener/Seeder)** and ingest a file: `./bin/peer -p 55555 -store ./store_a -ingest test.mp4`
     3. Take note of Peer A's full multiaddress, the generated `ContentID` (e.g., `abcd12345...`), and the `Key`.
     4. Open a second terminal window. Start **Peer B (Downloader)**, passing Peer A's address and fetching the content: 
        `./bin/peer -store ./store_b -d /ip4/127.0.0.1/tcp/55555/p2p/12D3... -fetch abcd12345... -key <KEY> -reassemble out.mp4`
     5. Observe Peer A's logs: It should log the ingestion and stream requests for manifest and chunks.
     6. Observe Peer B's logs: It should log resolving the manifest, downloading the chunks sequentially, and successfully reassembling the file.

   - **Relay Node Deployment**: To manually test the relay functionality (without full NAT traversal logic yet, just connection):
     1. Build the relay binary: `go build -o bin/relay ./cmd/relay`
     2. Start the relay: `./bin/relay`
     3. Ensure it outputs `Relay Service Started!` and prints its multiaddresses (e.g., `/ip4/127.0.0.1/tcp/4001/p2p/...`).
     4. To verify routing traffic, proceed to the **Relay Connectivity** step below.

   - **Relay Connectivity & Content Transfer**: To verify that a peer (e.g., Mac) can connect to another peer (e.g., Windows) through a public relay (e.g., Ubuntu server) and transfer content:
     1. Start your public relay node on the Ubuntu server and copy its multiaddress (e.g., `/ip4/<PUBLIC-IP>/tcp/4001/p2p/<RELAY-ID>`).
     2. On **Peer A** (Listener, e.g., Windows), start the peer, provide the relay, and ingest a file:
        `./bin/peer -p 55555 -store ./store_a -relay /ip4/<PUBLIC-IP>/tcp/4001/p2p/<RELAY-ID> -ingest test.mp4`
     3. Take note of Peer A's generated Relay Multiaddress, the returned `ContentID`, and the `Key`.
     4. On **Peer B** (Dialer, e.g., Mac), start Peer B, provide the static relay, dial Peer A's relayed address, and fetch the content:
        `./bin/peer -store ./store_b -relay /ip4/<PUBLIC-IP>/tcp/4001/p2p/<RELAY-ID> -d /ip4/<PUBLIC-IP>/tcp/4001/p2p/<RELAY-ID>/p2p-circuit/p2p/<PEER-A-ID> -fetch <ContentID> -key <KEY> -reassemble out.mp4`
     5. **Verify Hole Punching**: After the initial relayed connection, both peers should log `[DCUtR] Hole Punch Event: StartHolePunch`. Following a successful attempt, you should see `EndHolePunch`, and the network connection will naturally upgrade to direct routing.
     6. **Verify Transfer Integrity**: Peer B will download all encrypted chunks sequentially, verifying their hashes, and finally reassemble them into `out.mp4`.

   - **Relay Fallback**: To verify that the system can still transfer data over the relay when direct connection fails (or is disabled):
     1. Start **Peer A** with a relay multiaddress and `-p 50001`:
        `./bin/peer -p 50001 -store ./store_a -relay /ip4/<PUBLIC-IP>/tcp/4001/p2p/<RELAY-ID> -ingest test.mp4`
     2. Start **Peer B** with the `-force-relay` flag to disable hole punching:
        `./bin/peer -p 50002 -store ./store_b -relay /ip4/<PUBLIC-IP>/tcp/4001/p2p/<RELAY-ID> -d /ip4/<PUBLIC-IP>/tcp/4001/p2p/<RELAY-ID>/p2p-circuit/p2p/<PEER-A-ID> -force-relay -fetch <ContentID> -key <KEY> -reassemble out.mp4`
     3. Verify that the file transfers successfully. The network logs should explicitly say `1) Relay [...]` when listing active connections to the target peer.

   - **Direct Upgrade (Hole Punching)**: To verify natural direct upgrade:
     1. Start Peer A and B with a relay, but without the `-force-relay` flag.
     2. Start a large transfer (e.g. 50-100MB).
     3. During the transfer, you should observe `[DCUtR] Hole Punch Event: StartHolePunch` and `EndHolePunch`. 
     4. Note that for a single stream, the underlying connection used is decided when the stream is created. Any *new* stream after successful hole punching will use the direct connection (you will see `Path: Direct` in the transfer logs). For large transfers initiated before hole punch completes, they may finish over Relay depending on the exact libp2p stream multiplexer routing. Wait for `EndHolePunch` to finish and initiate another transfer to verify it says `Path: Direct` and is much faster.

    - **Integrity Verification**: 
     1. This is automatically tested on every block fetched by the client via the verify-then-store mechanism.
     2. If a chunk gets corrupted, the client drops the transfer with a hash mismatch error before the chunk touches the disk.

## Reliable Content Transfer (Milestone 9)

Milestone 9 introduces client-side transfer session tracking, retry policies, and resume capability while maintaining a strictly stateless `/cipher/chunk/1.0.0` protocol.

- [☑️] **Local Session Persistence**: `TransferSession` correctly persists progress in `sessions/` directory.
- [☑️] **Skip Logic (Idempotency)**: The downloader verifies local chunk presence via `HasChunk()` and successfully skips downloading already retrieved chunks.
- [☑️] **Resume Interrupted Downloads**: Killing a peer (`Ctrl+C`) halfway through a transfer and running `-resume` correctly starts from the remaining chunks without duplicates.
- [☑️] **Retry Policy**: Transient network failures safely trigger an exponential backoff instead of terminating the download.
- [☑️] **Session Management CLI**: `-transfer-status` and `-cancel` accurately list and clear local session state.

## Content Transport Integration Test Matrix (Milestone 8)

The transport layer leverages `/cipher/chunk/1.0.0` to explicitly separate network transmission from filesystem storage. It is considered validated only after all the following tests pass:

- [☑️] **Message Serialization**: Symmetric Message envelope (Version, Type, Payload) handles valid and malformed structures correctly.
- [☑️] **Verify-Then-Store**: Chunks are verified against their SHA-256 ContentID before being written to the engine store.
- [☑️] **Interruption Handling**: Peer disconnects midway through manifest or chunk transfers are caught and resources are freed safely.
- [☑️] **Decoupled Cryptography**: Reassembling downloaded chunks correctly fails without the explicit out-of-band encryption Key.
- [☑️] **Transport Benchmark**: Baseline transfer rate established (~490 MB/s over loopback).

## Content Engine Test Matrix (Milestone 7)

The Content Engine Foundation completely decouples the filesystem from the transport layer. It is tested strictly via local automation before any P2P network integration.

To run the Content Engine test suite:
```bash
go test -v ./internal/content/...
```

The engine is considered validated only after all the following automated tests pass:

- [☑️] **Chunking**: Streams are correctly split into precisely sized chunks based on dynamic configuration.
- [☑️] **Encryption (XChaCha20-Poly1305)**: Every chunk is independently encrypted in-place using a 192-bit nonce, modifying the `CipherSize` correctly.
- [☑️] **Integrity & Hashing**: `Digest` (SHA-256) correctly hashes ciphertexts to yield `ChunkID`s, and securely verifies them before decryption.
- [☑️] **Decoupled Storage**: The `ChunkSource` and `ChunkSink` interfaces successfully store and retrieve chunks from the filesystem using content-addressed filenames.
- [☑️] **End-to-End Reassembly**: A large data stream is successfully ingested, chunked, encrypted, hashed, stored, retrieved, decrypted, and accurately reassembled back into its original sequence using the immutable `Manifest`.
- [☑️] **Corruption Handling**: Altering a stored chunk reliably triggers a `hash mismatch` or decryption failure upon reassembly.

### Robustness Tests (Advanced Validation)

While the above matrix represents the minimum acceptance criteria for the Content Engine, the following robustness tests validate its readiness for a decentralized network environment:

- [☑️] **Boundary Testing**: File sizes precisely hitting chunk boundaries (e.g., 0B, 1B, 31KB, 32KB, 32KB+1B, 1MB, 100MB).
- [ ] **Randomized Inputs**: End-to-end ingest and reassembly of randomly generated binary files (not just text) with verification via SHA-256.
- [ ] **Out-of-Order Reconstruction**: Reassembling a stream from chunks requested and returned in a completely random sequence.
- [ ] **Duplicate Chunk Handling**: Simulating swarming by providing duplicate chunks and verifying no corruption or duplicate writes occur.
- [☑️] **Missing Chunk Handling**: Explicitly deleting a chunk and asserting the engine cleanly returns `ErrMissingChunk` rather than panicking or producing partial output.
- [☑️] **Cryptographic Rejection**: Providing the wrong decryption key and ensuring the engine fails securely.
- [ ] **Manifest Validation**: Rejecting invalid manifests (duplicate IDs, invalid offsets, malformed JSON, etc) before any processing begins.
- [ ] **Determinism**: Asserting that ingesting the exact same file twice with the same key produces the expected object behavior (independent ciphertexts due to nonces).
- [ ] **Performance & Streaming**: Ingesting a 500MB file and asserting that memory remains strictly bounded (no full-file buffering).
- [ ] **API Contracts**: Running identical test suites against any future `ChunkSource` (e.g. Memory, IPFS, S3) to ensure seamless swapping.
- [☑️] **The 1000-Iteration Gauntlet**: Running 1000 loops of: Generate random file -> Random chunk size -> Random key -> Ingest -> Shuffle chunks -> Restore -> Compare SHA-256.

### Manual Testing (Content Engine CLI)

To manually test the pipeline, you can use the `content-test` CLI utility. As development progresses, this tool supports extensive commands for debugging and manual validation.

1. **Build the Test CLI**:
   ```bash
   go build -o bin/content-test ./cmd/content-test
   ```

2. **Ingest a File**:
   Create a sample file and run the ingest command:
   ```bash
   mkdir -p test_files
   echo "This is a test file for the content engine." > test_files/sample.txt
   ./bin/content-test -ingest test_files/sample.txt
   ```
   *Expected Output*: The CLI will log the total chunks created and save `test_files/manifest.json`. Check the `./test_files/content_store/` directory to see the stored encrypted chunks sharded by their SHA-256 hashes.

3. **Reassemble the File**:
   Using only the generated manifest and the chunks in the local store, reassemble the original file:
   ```bash
   ./bin/content-test -out test_files/restored.txt
   ```
   *Expected Output*: The `restored.txt` file will be created and its contents should perfectly match `sample.txt`.

4. **Future CLI Commands (Planned)**:
   For debugging future distributed networks, the CLI will support advanced validation:
   ```bash
   ./bin/content-test verify test_manifest.json
   ./bin/content-test inspect test_manifest.json
   ./bin/content-test list
   ./bin/content-test corrupt <chunkID>
   ./bin/content-test delete <chunkID>
   ```

## Multi-peer Swarming (Milestone 10)

Milestone 10 introduced a deterministic chunk scheduler that enables concurrent downloading from multiple peers using a worker pool.

- [☑️] **Concurrent Downloading**: The scheduler distributes `ChunkTask` assignments to available workers via a thread-safe `ChunkQueue`, executing parallel downloads.
- [☑️] **Architectural Isolation**: `TransferManager` cleanly owns session state, lifecycle, and progress, while `Scheduler` solely handles work distribution. 
- [☑️] **Resilient Retry Policy**: Failed chunk fetches correctly push tasks back to the queue (up to a max attempt limit) rather than failing the entire download or infinitely looping.
- [☑️] **Multiple Target Input**: The CLI accepts a comma-separated list of multiaddresses (e.g., `-d <addr1>,<addr2>`) to dial multiple seeds simultaneously.

### Manual Swarm Testing (Phases)
1. **Phase 1: Concurrency**
   - *Important Note:* The Content Engine encrypts chunks independently using random 24-byte nonces. You cannot ingest the same file on three peers to simulate swarming, as they will generate different ContentIDs. Instead, ingest it once and have the others download it to become identical seeds.
   - **Step 1:** Start Peer A and ingest the file: 
     `CIPHER_CONFIG_DIR=/tmp/peerA ./bin/peer -p 55555 -store ./storeA -ingest test.mp4`
   - **Step 2:** Start Peer B and fetch from Peer A to become a seed:
     `CIPHER_CONFIG_DIR=/tmp/peerB ./bin/peer -p 55556 -store ./storeB -d <Peer A Address> -fetch <ContentID>`
   - **Step 3:** Start Peer C and fetch from Peer A (or B) to become a seed:
     `CIPHER_CONFIG_DIR=/tmp/peerC ./bin/peer -p 55557 -store ./storeC -d <Peer A Address> -fetch <ContentID>`
   - **Step 4:** Start the Downloader (Peer D) and fetch from all three seeds concurrently:
     `CIPHER_CONFIG_DIR=/tmp/peerD ./bin/peer -store ./storeD -d <PeerA_Addr>,<PeerB_Addr>,<PeerC_Addr> -fetch <ContentID> -reassemble out.mp4`
   - *Result Validation (Cross-Platform/Network):* Successfully verified downloading across completely separate networks (e.g. 1000km+ long distance) using a dedicated Azure VM as a `circuitv2` public relay. The `libp2p` instances flawlessly negotiated AutoNAT reachability and executed a DCUtR Hole Punch, shifting the data transfer off the Azure relay and onto a direct P2P TCP/UDP connection between the peers.
2. **Phase 2: Recovery**
   - Disconnect Peer B midway through the download.
   - Verify the scheduler correctly handles the network error, requeues the chunks, and finishes the download via Peers A and C.
3. **Phase 3: Partial Availability (Future)**
   - Test scheduling when peers only possess specific chunk ranges.

## Continuous Integration
Tests are intended to be executed automatically upon pull requests via standard CI pipelines to maintain code quality across iterations.
