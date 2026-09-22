package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strings"

	"cipher/network/content/core"
	"cipher/network/content/crypto"
	"cipher/network/content/engine"
	"cipher/network/content/storage"
	"cipher/network/content/verifier"
	"cipher/network/discovery"
	"cipher/network/identity"
	"cipher/network/protocol/chunk"
	"cipher/network/protocol/push"
	"cipher/network/transport"

	golog "github.com/ipfs/go-log/v2"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

func main() {
	golog.SetAllLoggers(golog.LevelWarn)

	port := flag.Int("p", 4001, "Port for the provider to listen on (TCP)")
	wsPort := flag.Int("ws-port", 4002, "Port for the provider to listen on (WebSocket, 0 to disable)")
	storePath := flag.String("store", "./provider_store", "Path to local content store directory")
	bootstrapAddr := flag.String("bootstrap", "", "Bootstrap peer multiaddress")
	relayAddr := flag.String("relay", "", "Static relay multiaddress to use for NAT traversal")
	forceRelay := flag.Bool("force-relay", false, "Force traffic over relay")
	republishHours := flag.Int("republish-interval", 12, "Interval in hours for DHT republisher")
	corruptProb := flag.Float64("test-corrupt-prob", 0.0, "Probability (0.0 to 1.0) of sending corrupt chunk for testing")
	identityPath := flag.String("identity", "", "Custom path to identity key file (optional)")
	allowPush := flag.Bool("allow-push", true, "Enable /cipher/push/1.0.0 remote ingestion protocol")
	pushAuthPolicy := flag.String("push-auth-policy", "open", "Push authorization policy: 'open' or 'allowlist'")
	pushAllowedPublishers := flag.String("push-allowed-publishers", "", "Comma-separated list of allowed publisher peer IDs (for allowlist policy)")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize identity
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

	// 2. Start libp2p host & DHT
	h, kdht, err := transport.NewNode(ctx, *port, *wsPort, priv, *relayAddr, *forceRelay)
	if err != nil {
		log.Fatalf("Failed to create provider node: %v", err)
	}
	defer h.Close()
	defer kdht.Close()

	// 3. Connect to DHT bootstrap if specified
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

	// Connect to relay if specified
	if *relayAddr != "" {
		relayInfo, err := peer.AddrInfoFromString(*relayAddr)
		if err == nil {
			if err := h.Connect(ctx, *relayInfo); err != nil {
				log.Printf("Warning: Failed to connect to relay: %v", err)
			} else {
				if res, err := client.Reserve(ctx, h, *relayInfo); err == nil {
					h.ConnManager().Protect(relayInfo.ID, "relay")
					log.Printf("[✓] Connected to relay and reserved slot (expires: %s)", res.Expiration.String())
				}
			}
		}
	}

	// 4. Initialize Content Store and Engine
	if err := storage.NewFSStorage(*storePath); err != nil {
		log.Fatalf("Failed to create store dir: %v", err)
	}
	config := core.EngineConfig{ChunkSize: 32 * 1024}
	enc := crypto.NewChaCha20Encryptor()
	dig := verifier.NewSHA256Digest()
	keys := engine.NewLocalKeyProvider()
	store := storage.NewFSStore(*storePath)
	eng := engine.NewContentEngine(config, enc, dig, store, store, keys, store)

	// Apply testing flags
	if *corruptProb > 0 {
		chunk.TestCorruptProb = *corruptProb
		log.Printf("[TESTING] Corrupt probability set to %.2f", *corruptProb)
	}

	// 5. Register Data-Plane Stream Handler (/cipher/chunk/1.0.0)
	chunk.NewStreamHandler(h, eng)

	// 6. Register Ingestion Stream Handler (/cipher/push/1.0.0)
	var allowedPublishersList []peer.ID
	if *pushAllowedPublishers != "" {
		for _, pidStr := range strings.Split(*pushAllowedPublishers, ",") {
			pidStr = strings.TrimSpace(pidStr)
			if pid, err := peer.Decode(pidStr); err == nil {
				allowedPublishersList = append(allowedPublishersList, pid)
			}
		}
	}
	push.NewStreamHandler(h, eng, kdht, *allowPush, push.AuthPolicy(*pushAuthPolicy), allowedPublishersList)

	// 7. Start Control-Plane DHT Announcements
	if *allowPush {
		discovery.StartStorageProviderHeartbeat(ctx, kdht, 10*time.Minute)
	}

	interval := time.Duration(*republishHours) * time.Hour
	discovery.StartRepublisher(ctx, kdht, store, interval)

	manifests, _ := store.ListManifests(ctx)

	fmt.Println("\n================= CIPHER PROVIDER =================")
	fmt.Printf("Provider Peer ID: %s\n", h.ID().String())
	fmt.Printf("Store Location  : %s\n", *storePath)
	fmt.Printf("Hosted Manifests: %d\n", len(manifests))
	fmt.Println("Listening Addresses:")
	for _, addr := range h.Addrs() {
		fmt.Printf("  - %s/p2p/%s\n", addr.String(), h.ID().String())
	}
	fmt.Println("===================================================")
	log.Println("Provider is ready and serving content. Press Ctrl+C to stop.")

	// Wait for OS shutdown signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Println("Shutting down provider...")
}
