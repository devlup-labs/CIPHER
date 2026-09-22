package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cipher/network/content/core"
	"cipher/network/content/crypto"
	"cipher/network/content/engine"
	"cipher/network/content/storage"
	"cipher/network/content/verifier"
	"cipher/network/discovery"
	"cipher/network/identity"
	"cipher/network/retrieval"
	"cipher/network/transfer/manager"
	"cipher/network/transfer/scheduler"
	"cipher/network/transport"

	golog "github.com/ipfs/go-log/v2"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

func main() {
	golog.SetAllLoggers(golog.LevelWarn)

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
		fmt.Printf("Session %x cancelled.\n", cID)
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

	if *bootstrapAddr != "" {
		bootstrapInfo, err := peer.AddrInfoFromString(*bootstrapAddr)
		if err != nil {
			log.Fatalf("Invalid bootstrap address: %v", err)
		}
		if err := discovery.Bootstrap(ctx, kdht, h, []peer.AddrInfo{*bootstrapInfo}); err != nil {
			log.Fatalf("Failed to bootstrap DHT: %v", err)
		}
		log.Printf("[DHT] Bootstrap complete. Routing table has %d peers", len(kdht.RoutingTable().ListPeers()))
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
		log.Printf("[TESTING] Throttling enabled (2MB/s)")
	}

	t := transport.NewTransport(h)
	var targetPeers []peer.ID

	// 4. Control Plane: Resolve Providers (Direct or via DHT)
	if *target != "" {
		log.Printf("Connecting directly to target provider(s): %s", *target)
		for _, targetStr := range strings.Split(*target, ",") {
			targetStr = strings.TrimSpace(targetStr)
			if targetStr == "" {
				continue
			}
			addrInfo, err := t.Connect(ctx, targetStr)
			if err != nil {
				log.Printf("Warning: Failed to connect to provider %s: %v", targetStr, err)
				continue
			}
			targetPeers = append(targetPeers, addrInfo.ID)
		}
		if len(targetPeers) == 0 {
			log.Fatalf("Fatal: Could not connect to any specified target providers")
		}
	} else {
		log.Printf("[DHT] Querying DHT control-plane for providers of ContentID %x...", contentID)
		providers, err := discovery.FindProviders(ctx, kdht, contentID, 5)
		if err != nil {
			log.Fatalf("[DHT] Provider discovery failed: %v", err)
		}
		if len(providers) == 0 {
			log.Fatalf("[DHT] No providers found for ContentID %x on DHT", contentID)
		}

		for _, p := range providers {
			log.Printf("[DHT] Discovered provider: %s", p.ID)
			if err := t.ConnectPeer(ctx, p); err != nil {
				log.Printf("[DHT] Failed to connect to provider %s: %v", p.ID, err)
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
	log.Printf("Resolving manifest for ContentID %x from %d provider(s)...", contentID, len(targetPeers))
	m, err := retrieval.ResolveManifest(ctx, contentID, kdht, t, eng, targetPeers)
	if err != nil {
		log.Fatalf("Failed to resolve manifest: %v", err)
	}
	log.Printf("[✓] Manifest resolved! Total chunks: %d", len(m.ChunkIDs))

	// 7. Data Plane: Parallel Swarming Chunk Download
	log.Printf("Downloading %d chunks from %d provider(s)...", len(m.ChunkIDs), len(targetPeers))
	tm := manager.NewTransferManager(sm, eng, t)
	if err := tm.Download(ctx, contentID, m.ChunkIDs, targetPeers); err != nil {
		log.Fatalf("Download failed: %v", err)
	}
	log.Printf("[✓] All %d chunks downloaded and verified successfully!", len(m.ChunkIDs))

	// 8. Content Engine: Decrypt & Reassemble
	if *reassembleOut != "" {
		if *keyHex == "" {
			log.Printf("Warning: No decryption key provided (-key). Attempting reassembly with cached keys...")
		}
		outF, err := os.Create(*reassembleOut)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer outF.Close()

		if err := eng.Reassemble(ctx, m, outF); err != nil {
			log.Fatalf("Reassembly failed: %v", err)
		}
		log.Printf("[✓] Content decrypted and reassembled to: %s", *reassembleOut)
	}
}
