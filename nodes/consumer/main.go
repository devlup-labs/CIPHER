package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"cipher/network/content/core"
	"cipher/network/content/crypto"
	"cipher/network/content/engine"
	"cipher/network/content/storage"
	"cipher/network/content/verifier"
	"cipher/network/discovery"
	"cipher/network/identity"
	"cipher/network/payments"
	"cipher/network/protocol/chunk"
	"cipher/network/retrieval"
	"cipher/network/transfer/manager"
	"cipher/network/transfer/scheduler"
	"cipher/network/transport"
	"cipher/shared/logger"

	"github.com/ethereum/go-ethereum/common"
	golog "github.com/ipfs/go-log/v2"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

func main() {
	golog.SetAllLoggers(golog.LevelWarn)

	log := logger.Consumer

	fetchID := flag.String("fetch", "", "ContentID to fetch (hex)")
	resumeID := flag.String("resume", "", "ContentID to resume downloading (hex)")
	keyHex := flag.String("key", "", "Decryption key (32-byte hex) for reassembly")
	reassembleOut := flag.String("out", "", "Output path to reassemble the decrypted file")

	port := flag.Int("p", 5001, "Port for the client to listen on (TCP)")
	wsPort := flag.Int("ws-port", 5002, "Port for the client to listen on (WebSocket, 0 to disable)")
	storePath := flag.String("store", "./client_store", "Path to local client cache store directory")

	target := flag.String("d", "", "Optional direct provider multiaddress(es), comma-separated. If omitted, providers are discovered via DHT.")
	bootstrapAddr := flag.String("bootstrap", "", "Bootstrap peer multiaddress")
	relayAddr := flag.String("relay", "", "Static relay multiaddress for NAT traversal")
	forceRelay := flag.Bool("force-relay", false, "Force traffic over relay")

	transferStatus := flag.Bool("status", false, "List all active transfer sessions")
	cancelID := flag.String("cancel", "", "ContentID to cancel and delete the transfer session")
	identityPath := flag.String("identity", "", "Custom path to identity key file (optional)")
	throttle := flag.String("throttle", "", "Throttle speed (e.g., 2MB) for testing")

	// Payments protocol flags
	ethRPC := flag.String("eth-rpc", "", "Ethereum JSON-RPC URL (e.g. http://127.0.0.1:8545)")
	ethKey := flag.String("eth-key", "", "Ethereum private key (hex)")
	entropyAddr := flag.String("entropy-addr", "", "CommitRevealEntropy contract address (hex)")
	providerEthAddr := flag.String("provider-eth-addr", "", "Provider Ethereum payout address (hex)")
	roundID := flag.Int64("round-id", 1, "Payment round ID")
	roundFaceValue := flag.String("round-face-value", "1000000000000000000", "Round face value in wei (default 1 ETH)")
	roleName := flag.String("role-name", "Swarm Consumer Client", "Human-readable role name for this consumer node")
	simulateCheat := flag.Bool("simulate-cheat", false, "Simulate malicious consumer generating forged/tampered EIP-712 payment tickets")
	keepAlive := flag.Bool("keep-alive", false, "Keep consumer running as local edge seeder / in-memory cache after download")

	flag.Parse()

	// 1. Session management commands that do not need network
	sm, err := manager.NewFileSessionManager(*storePath + "/sessions")
	if err != nil {
		log.Fatalf("Failed to initialize session manager: %v", err)
	}

	if *transferStatus {
		sessions, err := sm.List()
		if err != nil {
			log.Fatalf("Failed to list sessions: %v", err)
		}
		if len(sessions) == 0 {
			fmt.Println("No active transfer sessions.")
		} else {
			fmt.Println("Active Transfer Sessions:")
			for _, s := range sessions {
				fmt.Printf(" - ContentID: %x | Status: %s | Progress: %d/%d chunks | Target: %s\n",
					s.ContentID, s.Status, s.CompletedCount(), s.TotalChunks, s.TargetPeer.String())
			}
		}
		return
	}

	if *cancelID != "" {
		cIDBytes, err := hex.DecodeString(*cancelID)
		if err != nil || len(cIDBytes) != 32 {
			log.Fatalf("Invalid ContentID for cancel (must be 32 bytes hex)")
		}
		var cID core.ContentID
		copy(cID[:], cIDBytes)
		sm.Delete(cID)
		log.Success("Session %x cancelled.", cID)
		return
	}

	targetContentIDHex := *fetchID
	if targetContentIDHex == "" {
		targetContentIDHex = *resumeID
	}

	if targetContentIDHex == "" {
		log.Fatal("Must specify -fetch <ContentID> or -resume <ContentID> (or use -status / -cancel)")
	}

	cIDBytes, err := hex.DecodeString(targetContentIDHex)
	if err != nil || len(cIDBytes) != 32 {
		log.Fatalf("Invalid ContentID hex (must be 32 bytes)")
	}
	var contentID core.ContentID
	copy(contentID[:], cIDBytes)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize identity & network
	var priv libp2pcrypto.PrivKey
	if *identityPath != "" {
		priv, err = identity.LoadOrCreateFromPath(*identityPath)
	} else {
		priv, err = identity.LoadOrCreate()
	}
	if err != nil {
		log.Fatalf("Failed to load or create identity: %v", err)
	}

	h, kdht, err := transport.NewNode(ctx, *port, *wsPort, priv, *relayAddr, *forceRelay)
	if err != nil {
		log.Fatalf("Failed to create client libp2p node: %v", err)
	}
	defer h.Close()
	defer kdht.Close()

	fields := []logger.Field{
		{Key: "Consumer Role   ", Value: *roleName},
		{Key: "Consumer Peer ID", Value: h.ID().String()},
		{Key: "Target ContentID", Value: targetContentIDHex},
		{Key: "Cache Directory ", Value: *storePath},
	}
	log.Banner(fmt.Sprintf("CIPHER CONSUMER: %s", strings.ToUpper(*roleName)), fields...)

	if *bootstrapAddr != "" {
		bootstrapInfo, err := peer.AddrInfoFromString(*bootstrapAddr)
		if err != nil {
			log.Fatalf("Invalid bootstrap address: %v", err)
		}
		if err := discovery.Bootstrap(ctx, kdht, h, []peer.AddrInfo{*bootstrapInfo}); err != nil {
			log.Fatalf("Failed to bootstrap DHT: %v", err)
		}
		log.Sub("DHT").Success("Bootstrap complete. Routing table has %d peers", len(kdht.RoutingTable().ListPeers()))
	}

	if *relayAddr != "" {
		relayInfo, err := peer.AddrInfoFromString(*relayAddr)
		if err == nil {
			if err := h.Connect(ctx, *relayInfo); err == nil {
				client.Reserve(ctx, h, *relayInfo)
			}
		}
	}

	// 3. Initialize Content Engine
	if err := storage.NewFSStorage(*storePath); err != nil {
		log.Fatalf("Failed to create store dir: %v", err)
	}
	config := core.EngineConfig{ChunkSize: 32 * 1024}
	enc := crypto.NewChaCha20Encryptor()
	dig := verifier.NewSHA256Digest()
	keys := engine.NewLocalKeyProvider()
	store := storage.NewFSStore(*storePath)
	eng := engine.NewContentEngine(config, enc, dig, store, store, keys, store)

	if *throttle == "2MB" {
		scheduler.TestThrottle = 500 * time.Millisecond
		log.Warn("[TESTING] Throttling enabled (2MB/s)")
	}

	t := transport.NewTransport(h)
	var targetPeers []peer.ID

	// 4. Control Plane: Resolve Providers (Direct or via DHT)
	if *target != "" {
		log.Info("Connecting directly to target provider(s): %s", *target)
		for _, targetStr := range strings.Split(*target, ",") {
			targetStr = strings.TrimSpace(targetStr)
			if targetStr == "" {
				continue
			}
			addrInfo, err := t.Connect(ctx, targetStr)
			if err != nil {
				log.Warn("Failed to connect to provider %s: %v", targetStr, err)
				continue
			}
			targetPeers = append(targetPeers, addrInfo.ID)
		}
		if len(targetPeers) == 0 {
			log.Fatalf("Fatal: Could not connect to any specified target providers")
		}
	} else {
		log.Sub("DHT").Info("Querying DHT control-plane for providers of ContentID %x...", contentID)
		providers, err := discovery.FindProviders(ctx, kdht, contentID, 5)
		if err != nil {
			log.Sub("DHT").Fatalf("Provider discovery failed: %v", err)
		}
		if len(providers) == 0 {
			log.Sub("DHT").Fatalf("No providers found for ContentID %x on DHT", contentID)
		}

		for _, p := range providers {
			log.Sub("DHT").Success("Discovered provider: %s", p.ID)
			if err := t.ConnectPeer(ctx, p); err != nil {
				log.Sub("DHT").Warn("Failed to connect to provider %s: %v", p.ID, err)
				continue
			}
			targetPeers = append(targetPeers, p.ID)
		}
		if len(targetPeers) == 0 {
			log.Fatalf("Found providers on DHT, but failed to establish data connection to any")
		}
	}

	// 5. Store decryption key if provided
	if *keyHex != "" {
		kBytes, err := hex.DecodeString(*keyHex)
		if err != nil || len(kBytes) != 32 {
			log.Fatalf("Invalid key format (must be 32-byte hex)")
		}
		keys.Put(ctx, contentID, kBytes)
	}

	// 6. Data Plane: Resolve Manifest
	log.Info("Resolving manifest for ContentID %x from %d provider(s)...", contentID, len(targetPeers))
	m, err := retrieval.ResolveManifest(ctx, contentID, kdht, t, eng, targetPeers)
	if err != nil {
		log.Fatalf("Failed to resolve manifest: %v", err)
	}
	log.Success("Manifest resolved! Total chunks: %d", len(m.ChunkIDs))

	// 7. Data Plane: Parallel Swarming Chunk Download
	log.Sub("Swarm").Info("Downloading %d chunks from %d provider(s)...", len(m.ChunkIDs), len(targetPeers))
	tm := manager.NewTransferManager(sm, eng, t)

	// Configure payment ticket generation if Ethereum parameters provided
	if *ethRPC != "" && *ethKey != "" && *entropyAddr != "" && *providerEthAddr != "" {
		log.Sub("Payment").Info("Initializing Ethereum client & EIP-712 ticket signer...")
		payClient, err := payments.NewPaymentClient(ctx, *ethRPC, *ethKey, payments.ContractAddresses{
			EntropySource: common.HexToAddress(*entropyAddr),
		})
		if err != nil {
			log.Fatalf("Failed to initialize Ethereum payment client: %v", err)
		}
		defer payClient.Close()

		faceVal, ok := new(big.Int).SetString(*roundFaceValue, 10)
		if !ok {
			log.Fatalf("Invalid round face value: %s", *roundFaceValue)
		}

		pAddr := common.HexToAddress(*providerEthAddr)
		rID := big.NewInt(*roundID)

		// Map chunk IDs to local indices in manifest
		chunkIdxMap := make(map[core.ChunkID]uint64, len(m.ChunkIDs))
		for idx, cid := range m.ChunkIDs {
			chunkIdxMap[cid] = uint64(idx)
		}

		var ticketMu sync.Mutex
		tm.SetTicketGenerator(func(chunkID core.ChunkID) (*payments.SignedTicket, error) {
			ticketMu.Lock()
			defer ticketMu.Unlock()

			idx, exists := chunkIdxMap[chunkID]
			if !exists {
				return nil, fmt.Errorf("chunk %x not in manifest", chunkID)
			}

			ticket := payments.RoundTicket{
				Sender:      payClient.Address,
				Recipient:   pAddr,
				RoundID:     rID,
				LocalIndex:  new(big.Int).SetUint64(idx),
				FaceValue:   faceVal,
				WinProb:     big.NewInt(100000000000000000), // 0.1
				SenderNonce: big.NewInt(12345),
			}

			sig, err := payClient.Signer.SignTicket(ticket)
			if err != nil {
				return nil, fmt.Errorf("failed to sign ticket for chunk index %d: %w", idx, err)
			}

			if *simulateCheat {
				// Deliberately tamper with signature bytes to simulate fraudulent/counterfeit payment ticket
				for i := range sig {
					sig[i] ^= 0xFF
				}
				log.Sub("Payment").Warn("[FRAUD SIMULATION] Generated FORGED/TAMPERED EIP-712 ticket (invalid sig: %x...)", sig[:8])
			} else {
				log.Sub("Payment").Info("Generated & signed authentic ticket for chunk #%d (sig: %x...)", idx, sig[:8])
			}

			return &payments.SignedTicket{
				Ticket:    ticket,
				Signature: sig,
			}, nil
		})
	}

	if err := tm.Download(ctx, contentID, m.ChunkIDs, targetPeers); err != nil {
		log.Fatalf("Download failed: %v", err)
	}
	log.Success("All %d chunks downloaded and verified successfully!", len(m.ChunkIDs))

	// 8. Content Engine: Decrypt & Reassemble
	if *reassembleOut != "" {
		if *keyHex == "" {
			log.Warn("No decryption key provided (-key). Attempting reassembly with cached keys...")
		}
		outF, err := os.Create(*reassembleOut)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer outF.Close()

		if err := eng.Reassemble(ctx, m, outF); err != nil {
			log.Fatalf("Reassembly failed: %v", err)
		}
		log.Success("Content decrypted and reassembled to: %s", *reassembleOut)
	}

	// 9. Keep-Alive P2P Edge Seeder & In-Memory Cache Mode
	if *keepAlive {
		chunk.NewStreamHandler(h, eng)
		fields := []logger.Field{
			{Key: "Consumer Role   ", Value: *roleName + " [ACTIVE SEEDER]"},
			{Key: "Cache Directory ", Value: *storePath},
			{Key: "Decrypted Output", Value: *reassembleOut},
			{Key: "Cached Chunks   ", Value: fmt.Sprintf("%d chunks retained in local cache", len(m.ChunkIDs))},
			{Key: "Edge Protocols  ", Value: "/cipher/chunk/1.0.0, /cipher/pull/1.0.0"},
		}
		log.Banner("CIPHER CONSUMER: ACTIVE EDGE SEEDER & CACHE", fields...)
		log.Success("Consumer is actively listening and re-seeding content to peer swarm. Press Ctrl+C to stop.")

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Warn("Shutting down consumer edge seeder...")
	}
}

