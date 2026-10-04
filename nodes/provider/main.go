package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"cipher/integration/availability"
	payment "cipher/availability/escrow-payment/payment"
	"cipher/network/content/core"
	"cipher/network/content/crypto"
	"cipher/network/content/engine"
	"cipher/network/content/manifest"
	"cipher/network/content/storage"
	"cipher/network/content/verifier"
	"cipher/network/discovery"
	"cipher/network/identity"
	"cipher/network/payments"
	"cipher/network/protocol/chunk"
	"cipher/network/protocol/push"
	"cipher/network/transport"
	"cipher/shared/logger"

	"crypto/ed25519"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	golog "github.com/ipfs/go-log/v2"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

func main() {
	golog.SetAllLoggers(golog.LevelWarn)

	log := logger.Provider

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

	// Payments and Availability protocol flags
	ethRPC := flag.String("eth-rpc", "", "Ethereum JSON-RPC URL (e.g. http://127.0.0.1:8545)")
	entropyAddr := flag.String("entropy-addr", "", "CommitRevealEntropy contract address (hex)")
	providerEthKey := flag.String("eth-key", "", "Provider Ethereum private key (hex, optional)")
	enableAvailability := flag.Bool("availability", true, "Enable Availability challenge handler (/cipher/availability/1.0.0)")
	roleName := flag.String("role-name", "Core Storage Provider", "Human-readable role name for this provider node")
	_ = providerEthKey

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
		log.Sub("DHT").Success("Bootstrap complete. Routing table has %d peers", len(kdht.RoutingTable().ListPeers()))
	}

	// Connect to relay if specified
	if *relayAddr != "" {
		relayInfo, err := peer.AddrInfoFromString(*relayAddr)
		if err == nil {
			if err := h.Connect(ctx, *relayInfo); err != nil {
				log.Warn("Failed to connect to relay: %v", err)
			} else {
				if res, err := client.Reserve(ctx, h, *relayInfo); err == nil {
					h.ConnManager().Protect(relayInfo.ID, "relay")
					log.Success("Connected to relay and reserved slot (expires: %s)", res.Expiration.String())
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
		log.Warn("[TESTING] Corrupt probability set to %.2f", *corruptProb)
	}

	// 5. Register Data-Plane Stream Handler (/cipher/chunk/1.0.0)
	streamHandler := chunk.NewStreamHandler(h, eng)
	if *ethRPC != "" && *entropyAddr != "" {
		log.Sub("Payment").Info("Configuring payment ticket verification against CommitRevealEntropy %s...", *entropyAddr)
		ethClient, err := ethclient.DialContext(ctx, *ethRPC)
		if err != nil {
			log.Fatalf("Failed to dial Ethereum RPC: %v", err)
		}
		defer ethClient.Close()

		chainID, err := ethClient.ChainID(ctx)
		if err != nil {
			log.Fatalf("Failed to retrieve chain ID: %v", err)
		}

		verifierSigner := payments.NewTicketSigner(nil, chainID, common.HexToAddress(*entropyAddr))

		var ticketsMu sync.Mutex
		var storedTickets []*payments.SignedTicket

		streamHandler.SetTicketHandler(func(ticket *payments.SignedTicket) error {
			ticketsMu.Lock()
			defer ticketsMu.Unlock()

			if !verifierSigner.Verify(ticket.Ticket, ticket.Signature, ticket.Ticket.Sender) {
				log.Sub("Payment").Error("[SECURITY SHIELD] REJECTED FRAUDULENT TICKET from %s! Forged EIP-712 signature detected. Chunk transfer DENIED.", ticket.Ticket.Sender.Hex())
				return fmt.Errorf("invalid ticket signature from %s", ticket.Ticket.Sender.Hex())
			}

			storedTickets = append(storedTickets, ticket)
			log.Sub("Payment").Success("Received & verified authentic ticket #%d (Chunk #%s, Value %s wei, from %s)",
				len(storedTickets), ticket.Ticket.LocalIndex.String(), ticket.Ticket.FaceValue.String(), ticket.Ticket.Sender.Hex())
			return nil
		})
	}

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

	// 8. Register Availability Protocol Handler (/cipher/availability/1.0.0)
	if *enableAvailability {
		rawPriv, err := priv.Raw()
		if err == nil {
			var edPriv ed25519.PrivateKey
			if len(rawPriv) == 64 {
				edPriv = ed25519.PrivateKey(rawPriv)
			} else if len(rawPriv) == 32 {
				edPriv = ed25519.NewKeyFromSeed(rawPriv)
			}

			if len(edPriv) == ed25519.PrivateKeySize {
				proofEngine, err := availability.NewProviderProofEngine(h.ID().String(), edPriv, store)
				if err == nil {
					// Register all existing manifests in proof engine
					for _, mid := range manifests {
						if mData, err := store.GetManifestBytes(ctx, mid); err == nil {
							if m, err := manifest.Deserialize(mData); err == nil {
								var rawChunks [][]byte
								for _, cid := range m.ChunkIDs {
									if ch, err := store.GetChunk(ctx, cid); err == nil {
										rawChunks = append(rawChunks, ch.Data)
									}
								}
								if len(rawChunks) > 0 {
									_, _ = proofEngine.RegisterFileChunks(hex.EncodeToString(m.Descriptor.ID[:]), rawChunks)
								}
							}
						}
					}

					availability.NewAvailabilityStreamHandler(h, proofEngine, func(voucher payment.PaymentState) {
						log.Sub("Availability").Success("Received signed payment voucher (Seq: %d, Cumulative: %d wei)",
							voucher.Sequence, voucher.CumulativePayment)
					})
					log.Sub("Availability").Info("Stream handler active on %s", availability.AvailabilityProtocolID)
				}
			}
		}
	}

	fields := []logger.Field{
		{Key: "Provider Role   ", Value: *roleName},
		{Key: "Provider Peer ID", Value: h.ID().String()},
		{Key: "Store Location  ", Value: *storePath},
		{Key: "Hosted Manifests", Value: fmt.Sprintf("%d", len(manifests))},
		{Key: "Push Ingestion  ", Value: fmt.Sprintf("%t (policy: %s)", *allowPush, *pushAuthPolicy)},
		{Key: "", Value: "Listening Multiaddresses:"},
	}
	for _, addr := range h.Addrs() {
		fields = append(fields, logger.Field{Key: "", Value: fmt.Sprintf("  - %s/p2p/%s", addr.String(), h.ID().String())})

	}

	log.Banner(fmt.Sprintf("CIPHER PROVIDER: %s", strings.ToUpper(*roleName)), fields...)
	log.Success("Provider [%s] is ready and serving content. Press Ctrl+C to stop.", *roleName)

	// Wait for OS shutdown signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Warn("Shutting down provider...")
}

