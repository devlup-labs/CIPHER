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

	"cipher/network/content/core"
	"cipher/network/content/crypto"
	"cipher/network/content/engine"
	"cipher/network/content/manifest"
	"cipher/network/content/storage"
	"cipher/network/content/verifier"
	"cipher/network/discovery"
	"cipher/network/identity"
	"cipher/network/protocol/chunk"
	"cipher/network/retrieval"
	"cipher/network/transfer/manager"
	"cipher/network/transfer/scheduler"
	"cipher/network/transport"

	"encoding/hex"

	golog "github.com/ipfs/go-log/v2"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

func main() {
	// Enable all libp2p debug logging to diagnose AutoNAT
	// Reduce libp2p internal logging (DHT, etc.)
	golog.SetAllLoggers(golog.LevelWarn)

	// Set all the flags
	target := flag.String("d", "", "Optional target peer multiaddress. If omitted, providers are discovered through the DHT. Target peer multiaddress to dial (e.g. /ip4/127.0.0.1/tcp/55555/p2p/Qm...)")
	port := flag.Int("p", 4001, "Port for the peer to listen on (TCP)")
	wsPort := flag.Int("ws-port", 4002, "Port for the peer to listen on (WebSocket, 0 to disable)")
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

	bootstrapAddr := flag.String("bootstrap", "", "Bootstrap peer multiaddress")

	identityPath := flag.String("identity", "", "Custom path to identity key file (optional)")
	throttle := flag.String("throttle", "", "Throttle speed (e.g., 2MB) per second")
	corruptProb := flag.Float64("test-corrupt-prob", 0.0, "Probability (0.0 to 1.0) of sending a corrupt chunk for testing")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load or create private key for the newNode
	var priv libp2pcrypto.PrivKey
	var err error
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
		log.Fatalf("Failed to create libp2p node: %v", err)
	}

	// Bootstrap the node if a bootstrap address is provided
	if *bootstrapAddr != "" {
		bootstrapInfo, err := peer.AddrInfoFromString(*bootstrapAddr)
		if err != nil {
			log.Fatalf(
				"Invalid bootstrap address: %v",
				err,
			)
		}

		seeds := []peer.AddrInfo{
			*bootstrapInfo,
		}

		if err := discovery.Bootstrap(
			ctx,
			kdht,
			h,
			seeds,
		); err != nil {
			log.Fatalf(
				"Failed to bootstrap DHT: %v",
				err,
			)
		}

		log.Printf(
			"[DHT] Routing table now contains %d peers",
			len(kdht.RoutingTable().ListPeers()),
		)
	}

	defer kdht.Close()

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
	eng := engine.NewContentEngine(config, enc, dig, store, store, keys, store)

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

	chunk.NewStreamHandler(h, eng) // mp duplicate, have called it again later

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

	// Start DHT Republisher for persistent provider lifecycle
	discovery.StartRepublisher(ctx, kdht, store, 12*time.Hour)

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

		// Ingest reads a file, chunks it, encrypts it, stores it, and returns the manifest.
		m, err := eng.Ingest(ctx, f, manifest.TypeFile)
		if err != nil {
			log.Fatalf("Failed to ingest: %v", err)
		}

		// Save manifest bytes to engine memory so it can be served
		mBytes, _ := m.Serialize()
		eng.PutManifestBytes(ctx, m.Descriptor.ID, mBytes)

		// Advertise/ broadcast the content on the DHT

		log.Printf(
			"[DHT] Providing ContentID %s",
			m.Descriptor.ID,
		)

		if err := discovery.Provide(
			ctx,
			kdht,
			m.Descriptor.ID,
		); err != nil {
			log.Printf(
				"[DHT] Failed to advertise content %s: %v",
				m.Descriptor.ID,
				err,
			)
		} else {
			log.Printf(
				"[DHT] Successfully advertised content %s",
				m.Descriptor.ID,
			)
		}

		key, _ := keys.Get(ctx, m.Descriptor.ID)
		log.Printf("[✓] Ingest complete!")
		log.Printf("    ContentID: %x", m.Descriptor.ID)
		log.Printf("    Key: %x", key)

		log.Printf("\n--- To download this file on another peer (Peer B), run: ---")
		wsAddr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/ws/p2p/%s", *wsPort, h.ID())
		if *wsPort == 0 {
			wsAddr = fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/p2p/%s", *port, h.ID())
		}
		
		fmt.Printf("CGO_ENABLED=0 go run cmd/peer/main.go \\\n" +
			"  -p 5001 \\\n" +
			"  -ws-port 5002 \\\n" +
			"  -store ./store_b \\\n" +
			"  -d \"%s\" \\\n" +
			"  -fetch \"%x\" \\\n" +
			"  -key \"%x\" \\\n" +
			"  -reassemble \"downloaded_file\"\n", wsAddr, m.Descriptor.ID, key)
		log.Printf("-----------------------------------------------------------\n")
	}

	var targetContentIDHex string
	if *fetchID != "" {
		targetContentIDHex = *fetchID
	} else if *resumeID != "" {
		targetContentIDHex = *resumeID
	}

	if targetContentIDHex != "" {
		cIDBytes, err := hex.DecodeString(targetContentIDHex)
		if err != nil || len(cIDBytes) != 32 {
			log.Fatalf("Invalid ContentID hex")
		}
		var contentID core.ContentID
		copy(contentID[:], cIDBytes)

		// Create a transport instance for connecting to peers
		t := transport.NewTransport(h)

		var targetPeers []peer.ID

		// If a target peer is specified, connect to it. Otherwise, search for providers of the content ID on the DHT, and then connect to them.
		if *target != "" {

			log.Printf("Dialing target peer(s): %s", *target)

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
		} else {
			log.Printf(
				"[DHT] No target peer supplied. Searching for providers of %x...",
				contentID,
			)

			providers, err := discovery.FindProviders(
				ctx,
				kdht,
				contentID,
				3,
			)
			if err != nil {
				log.Fatalf(
					"[DHT] Failed to find providers: %v",
					err,
				)
			}

			if len(providers) == 0 {
				log.Fatalf(
					"[DHT] No providers found for ContentID %x",
					contentID,
				)
			}

			for _, provider := range providers {

				log.Printf(
					"[DHT] Found provider: %s",
					provider.ID,
				)

				// Connect using your existing Transport.
				if err := t.ConnectPeer(ctx, provider); err != nil {
					log.Printf(
						"[DHT] Failed to connect to provider %s: %v",
						provider.ID,
						err,
					)
					continue
				}

				targetPeers = append(
					targetPeers,
					provider.ID,
				)
			}

			if len(targetPeers) == 0 {
				log.Fatalf(
					"[DHT] Found providers, but could not connect to any of them",
				)
			}
		}

		if *keyHex != "" {
			kBytes, err := hex.DecodeString(*keyHex)
			if err != nil || len(kBytes) != 32 {
				log.Fatalf("Invalid key hex format or length (must be 32 bytes)")
			}
			keys.Put(ctx, contentID, kBytes)
		}

		// ResolveManifest is a new function that encapsulates the logic of resolving the manifest from the target peers.
		m, err := retrieval.ResolveManifest(ctx, contentID, kdht, t, eng, targetPeers)
		if err != nil {
			log.Fatalf("Failed to resolve manifest: %v", err)
		}

		// Setup Transfer Manager and start download
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
