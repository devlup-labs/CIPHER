package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cipher/internal/content/core"
	"cipher/internal/content/crypto"
	"cipher/internal/content/engine"
	"cipher/internal/content/manifest"
	"cipher/internal/content/storage"
	"cipher/internal/content/verifier"
	"cipher/internal/identity"
	"cipher/internal/protocol/chunk"
	"cipher/internal/transfer/manager"
	"cipher/internal/transfer/scheduler"
	"cipher/internal/transport"

	"encoding/hex"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	golog "github.com/ipfs/go-log/v2"
)

func main() {
	// Enable all libp2p debug logging to diagnose AutoNAT
	golog.SetAllLoggers(golog.LevelDebug)
	target := flag.String("d", "", "Target peer multiaddress to dial (e.g. /ip4/127.0.0.1/tcp/55555/p2p/Qm...)")
	port := flag.Int("p", 0, "Port to listen on (default 0 for random)")
	relayAddr := flag.String("relay", "", "Static relay multiaddress to use for NAT traversal")
	forceRelay := flag.Bool("force-relay", false, "Disable hole punching and force traffic over the relay")
	storePath := flag.String("store", "./content_store", "Path to the local content store directory")
	
	// Milestone 8 flags
	ingestFile := flag.String("ingest", "", "Path to file to ingest locally")
	fetchID := flag.String("fetch", "", "ContentID to fetch from target peer")
	reassembleOut := flag.String("reassemble", "", "Output path to reassemble the fetched ContentID")
	keyHex := flag.String("key", "", "Decryption key (hex) for reassembly")
	resumeID := flag.String("resume", "", "ContentID to resume downloading")
	transferStatus := flag.Bool("transfer-status", false, "List all active transfer sessions")
	cancelID := flag.String("cancel", "", "ContentID to cancel and delete the transfer session")

	// Testing Flags
	throttle := flag.String("throttle", "", "Throttle speed (e.g., 2MB) per second")
	corruptProb := flag.Float64("test-corrupt-prob", 0.0, "Probability (0.0 to 1.0) of sending a corrupt chunk for testing")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priv, err := identity.LoadOrCreate()
	if err != nil {
		log.Fatalf("Failed to load or create identity: %v", err)
	}

	h, err := transport.NewNode(ctx, *port, priv, *relayAddr, *forceRelay)
	if err != nil {
		log.Fatalf("Failed to create libp2p node: %v", err)
	}

	// Setup Content Engine
	if err := storage.NewFSStorage(*storePath); err != nil {
		log.Fatalf("Failed to create store dir: %v", err)
	}
	config := core.EngineConfig{ChunkSize: 32 * 1024}
	enc := crypto.NewChaCha20Encryptor()
	dig := verifier.NewSHA256Digest()
	keys := engine.NewLocalKeyProvider()
	store := storage.NewFSStore(*storePath)
	// Passing engineLogger isn't supported yet, removing it.
	eng := engine.NewContentEngine(config, enc, dig, store, store, keys)
	
	// Apply testing flags
	if *corruptProb > 0 {
		chunk.TestCorruptProb = *corruptProb
		log.Printf("[TESTING] Chunk corruption probability set to %.2f", *corruptProb)
	}
	if *throttle == "2MB" {
		// 2MB/s = 8 chunks/sec (256KB each). Sleep 125ms per chunk.
		scheduler.TestThrottle = 500 * time.Millisecond
		log.Printf("[TESTING] Throttling enabled (2MB/s)")
	}

	chunk.NewStreamHandler(h, eng)

	sm, err := manager.NewFileSessionManager(*storePath + "/sessions")
	if err != nil {
		log.Fatalf("Failed to create session manager: %v", err)
	}

	if *transferStatus {
		sessions, err := sm.List()
		if err != nil {
			log.Fatalf("Failed to list sessions: %v", err)
		}
		if len(sessions) == 0 {
			fmt.Println("No active transfer sessions.")
		} else {
			fmt.Println("Transfer Sessions:")
			for _, s := range sessions {
				fmt.Printf(" - ContentID: %x | Status: %s | Progress: %d/%d chunks | Target: %s\n", s.ContentID, s.Status, s.CompletedCount(), s.TotalChunks, s.TargetPeer.String())
			}
		}
		return
	}

	if *cancelID != "" {
		cIDBytes, _ := hex.DecodeString(*cancelID)
		var cID core.ContentID
		copy(cID[:], cIDBytes)
		sm.Delete(cID)
		fmt.Printf("Session %x cancelled.\n", cID)
		return
	}

	// Setup protocol handler
	chunk.NewStreamHandler(h, eng)

	log.Printf("Peer initialized with ID: %s", h.ID().String())
	log.Println("Listening on the following local addresses:")
	for _, addr := range h.Addrs() {
		log.Printf("  - %s/p2p/%s", addr.String(), h.ID().String())
	}

	if *relayAddr != "" {
		relayInfo, err := peer.AddrInfoFromString(*relayAddr)
		if err == nil {
			// Proactively connect and explicitly reserve a slot on the relay
			if err := h.Connect(ctx, *relayInfo); err != nil {
				log.Printf("Warning: Failed to connect to relay: %v", err)
			} else {
				if res, err := client.Reserve(ctx, h, *relayInfo); err != nil {
					log.Printf("Warning: Failed to reserve slot on relay: %v", err)
				} else {
					h.ConnManager().Protect(relayInfo.ID, "relay") // Prevent idle timeout
					log.Printf("\n[✓] Successfully connected to relay and reserved slot!")
					log.Printf("    Reservation Expiration: %s", res.Expiration.String())
					log.Printf("    Relay Peer ID: %s", relayInfo.ID.String())
					log.Printf("Your Relayed Multiaddress (Share this with peers to connect to you):")
					log.Printf("  - %s/p2p-circuit/p2p/%s\n", *relayAddr, h.ID().String())
				}
			}
		}
	}

	if *ingestFile != "" {
		log.Printf("Ingesting file: %s", *ingestFile)
		f, err := os.Open(*ingestFile)
		if err != nil {
			log.Fatalf("Failed to open ingest file: %v", err)
		}
		defer f.Close()

		m, err := eng.Ingest(ctx, f, manifest.TypeFile)
		if err != nil {
			log.Fatalf("Failed to ingest: %v", err)
		}
		// Save manifest bytes to engine memory so it can be served
		mBytes, _ := m.Serialize()
		eng.PutManifestBytes(ctx, m.Descriptor.ID, mBytes)
		
		key, _ := keys.Get(ctx, m.Descriptor.ID)
		log.Printf("[✓] Ingest complete!")
		log.Printf("    ContentID: %x", m.Descriptor.ID)
		log.Printf("    Key: %x", key)
	}

	var targetContentIDHex string
	if *fetchID != "" {
		targetContentIDHex = *fetchID
	} else if *resumeID != "" {
		targetContentIDHex = *resumeID
	}

	if *target != "" && targetContentIDHex != "" {
		log.Printf("Dialing target peer(s): %s", *target)
		t := transport.NewTransport(h)
		
		var targetPeers []peer.ID
		for _, targetStr := range strings.Split(*target, ",") {
			targetStr = strings.TrimSpace(targetStr)
			if targetStr == "" {
				continue
			}
			addrInfo, err := t.Connect(ctx, targetStr)
			if err != nil {
				log.Printf("Warning: Failed to connect to target %s: %v", targetStr, err)
				continue
			}
			targetPeers = append(targetPeers, addrInfo.ID)
		}

		if len(targetPeers) == 0 {
			log.Fatalf("Fatal: Could not connect to any target peers")
		}

		cIDBytes, err := hex.DecodeString(targetContentIDHex)
		if err != nil || len(cIDBytes) != 32 {
			log.Fatalf("Invalid ContentID hex")
		}
		var contentID core.ContentID
		copy(contentID[:], cIDBytes)

		if *keyHex != "" {
			kBytes, err := hex.DecodeString(*keyHex)
			if err != nil || len(kBytes) != 32 {
				log.Fatalf("Invalid key hex format or length (must be 32 bytes)")
			}
			keys.Put(ctx, contentID, kBytes)
		}

		// Resolve manifest from the first peer
		client, err := chunk.NewClient(ctx, t, targetPeers[0], eng)
		if err != nil {
			log.Fatalf("Failed to create chunk client: %v", err)
		}
		
		log.Printf("Resolving manifest for ContentID: %x from %s", contentID, targetPeers[0])
		mData, err := client.Resolve(ctx, contentID)
		client.Close() // Can close after resolve, workers will make their own
		
		if err != nil {
			log.Fatalf("Failed to resolve manifest: %v", err)
		}

		m, err := manifest.Deserialize(mData)
		if err != nil {
			log.Fatalf("Failed to deserialize manifest: %v", err)
		}

		log.Printf("Downloading %d chunks from %d peers...", len(m.ChunkIDs), len(targetPeers))
		
		tm := manager.NewTransferManager(sm, eng, t)
		if err := tm.Download(ctx, contentID, m.ChunkIDs, targetPeers); err != nil {
			log.Fatalf("Download failed: %v", err)
		}

		log.Printf("[✓] Download complete!")

		if *reassembleOut != "" {
			outF, err := os.Create(*reassembleOut)
			if err != nil {
				log.Fatalf("Failed to create out file: %v", err)
			}
			defer outF.Close()
			if err := eng.Reassemble(ctx, m, outF); err != nil {
				log.Fatalf("Reassemble failed: %v", err)
			}
			log.Printf("[✓] Reassembled to: %s", *reassembleOut)
		}
	}

	// Wait for termination signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println()
	log.Println("Shutting down peer...")
	if err := h.Close(); err != nil {
		log.Fatalf("Failed to close host: %v", err)
	}
}
