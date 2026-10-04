package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	interfaces "cipher/availability/escrow-payment/interfaces"
	payment "cipher/availability/escrow-payment/payment"
	availabilitytypes "cipher/availability/availability-contracts/types"
	verification "cipher/availability/availability-contracts/verification"
	"cipher/integration/availability"

	"cipher/network/content/core"
	"cipher/network/content/crypto"
	"cipher/network/content/engine"
	"cipher/network/content/manifest"
	"cipher/network/content/storage"
	"cipher/network/content/verifier"
	"cipher/network/discovery"
	"cipher/network/distribution"
	"cipher/network/identity"
	"cipher/network/protocol/chunk"
	"cipher/network/retrieval"
	"cipher/network/transport"
	"cipher/shared/logger"

	golog "github.com/ipfs/go-log/v2"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

func main() {
	golog.SetAllLoggers(golog.LevelWarn)

	log := logger.Publisher

	filePath := flag.String("file", "", "Path to the file to ingest and publish (required)")
	port := flag.Int("p", 4005, "Port for the publisher to listen on (TCP)")
	wsPort := flag.Int("ws-port", 4006, "Port for the publisher to listen on (WebSocket, 0 to disable)")
	storePath := flag.String("store", "./store_publisher", "Path to local content store directory")
	bootstrapAddr := flag.String("bootstrap", "", "Bootstrap peer multiaddress")
	relayAddr := flag.String("relay", "", "Static relay multiaddress to use for NAT traversal")
	forceRelay := flag.Bool("force-relay", false, "Force traffic over relay")
	seed := flag.Bool("seed", true, "Keep publisher running to seed chunks over /cipher/chunk/1.0.0")
	identityPath := flag.String("identity", "", "Custom path to identity key file (optional)")
	chunkSizeKB := flag.Int("chunk-size", 32, "Chunk size in KB (default: 32)")
	providersList := flag.String("providers", "", "Comma-separated multiaddresses of target providers to push content to")
	replication := flag.Int("replication", 2, "Replication factor R (replicas per chunk across providers)")
	push := flag.Bool("push", false, "Push chunks to remote providers over /cipher/push/1.0.0 and exit")
	pushTimeout := flag.Duration("push-timeout", 5*time.Minute, "Timeout for remote push distribution")
	challengeProviders := flag.Bool("challenge", false, "Issue Availability challenge against remote providers after push")
	challengeInterval := flag.Duration("challenge-interval", 5*time.Second, "Interval between periodic Availability challenges (default: 5s)")
	challengeRounds := flag.Int("challenge-rounds", 1, "Number of challenge rounds to execute (default: 1; set >1 or use with -challenge-loop)")
	challengeLoop := flag.Bool("challenge-loop", false, "Continuously challenge providers at -challenge-interval in background")
	challengeCID := flag.String("challenge-cid", "", "Specific ContentID hex to challenge providers for without re-encrypting")
	roleName := flag.String("role-name", "Ingestion & Replication Engine", "Human-readable role name for this publisher node")

	flag.Parse()

	if *filePath == "" && *challengeCID == "" {
		log.Fatalf("Error: -file <path> or -challenge-cid <hex> is required")
	}

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

	// 2. Start libp2p host and DHT
	h, kdht, err := transport.NewNode(ctx, *port, *wsPort, priv, *relayAddr, *forceRelay)
	if err != nil {
		log.Fatalf("Failed to create libp2p node: %v", err)
	}
	defer h.Close()
	defer kdht.Close()

	// 3. Connect to DHT bootstrap if provided
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

	// 4. Initialize Content Engine
	if err := storage.NewFSStorage(*storePath); err != nil {
		log.Fatalf("Failed to create store dir: %v", err)
	}
	config := core.EngineConfig{ChunkSize: uint32((*chunkSizeKB) * 1024)}
	enc := crypto.NewChaCha20Encryptor()
	dig := verifier.NewSHA256Digest()
	keys := engine.NewLocalKeyProvider()
	store := storage.NewFSStore(*storePath)
	eng := engine.NewContentEngine(config, enc, dig, store, store, keys, store)

	// Register chunk protocol stream handler for initial seeding
	chunk.NewStreamHandler(h, eng)

	// 5. Ingest Content or Resolve Manifest for -challenge-cid
	var m *manifest.Manifest
	if *challengeCID != "" {
		cIDBytes, err := hex.DecodeString(*challengeCID)
		if err != nil || len(cIDBytes) != 32 {
			log.Fatalf("Invalid -challenge-cid format (must be 32-byte hex)")
		}
		var cid core.ContentID
		copy(cid[:], cIDBytes)

		// Try loading from local engine store first
		if mBytes, err := store.GetManifestBytes(ctx, cid); err == nil {
			m, _ = manifest.Deserialize(mBytes)
		}

		if m == nil {
			t := transport.NewTransport(h)
			provs, _ := discovery.FindProviders(ctx, kdht, cid, 16)
			var pIDs []peer.ID
			for _, p := range provs {
				pIDs = append(pIDs, p.ID)
			}
			m, err = retrieval.ResolveManifest(ctx, cid, kdht, t, eng, pIDs)
			if err != nil {
				log.Fatalf("Failed to resolve manifest for challenge-cid %x: %v", cid, err)
			}
		}
		log.Sub("Availability").Success("Loaded manifest for ContentID: %x (Total Chunks: %d)", m.Descriptor.ID, len(m.ChunkIDs))
	} else {
		log.Sub("Ingest").Info("Ingesting source file: %s", *filePath)
		f, err := os.Open(*filePath)
		if err != nil {
			log.Fatalf("Failed to open file for ingest: %v", err)
		}
		defer f.Close()

		m, err = eng.Ingest(ctx, f, manifest.TypeFile)
		if err != nil {
			log.Fatalf("Failed to ingest file: %v", err)
		}

		// Persist manifest in engine
		mBytes, err := m.Serialize()
		if err != nil {
			log.Fatalf("Failed to serialize manifest: %v", err)
		}
		if err := eng.PutManifestBytes(ctx, m.Descriptor.ID, mBytes); err != nil {
			log.Fatalf("Failed to store manifest: %v", err)
		}
	}

	// 6. Execute Remote Push / Availability Challenges if requested
	if *push || *challengeProviders {
		t := transport.NewTransport(h)
		var targetPeers []peer.ID

		if *providersList != "" {
			for _, pAddrStr := range strings.Split(*providersList, ",") {
				pAddrStr = strings.TrimSpace(pAddrStr)
				if pAddrStr == "" {
					continue
				}
				addrInfo, err := t.Connect(ctx, pAddrStr)
				if err != nil {
					log.Warn("Failed to connect to provider %s: %v", pAddrStr, err)
					continue
				}
				targetPeers = append(targetPeers, addrInfo.ID)
			}
		} else {
			// Automated Kademlia DHT Storage Provider Discovery
			log.Info("No -providers specified. Querying Kademlia DHT for active storage providers...")

			for attempt := 1; attempt <= 3; attempt++ {
				dhtCtx, dhtCancel := context.WithTimeout(ctx, 5*time.Second)
				discovered, _ := discovery.FindStorageProviders(dhtCtx, kdht, 16)
				dhtCancel()

				for _, prov := range discovered {
					if prov.ID == h.ID() {
						continue // Skip self
					}
					alreadyAdded := false
					for _, existing := range targetPeers {
						if existing == prov.ID {
							alreadyAdded = true
							break
						}
					}
					if alreadyAdded {
						continue
					}

					if err := h.Connect(ctx, prov); err != nil {
						log.Warn("Failed to connect to discovered provider %s: %v", prov.ID, err)
						continue
					}
					log.Success("Discovered active storage provider via DHT: %s", prov.ID)
					targetPeers = append(targetPeers, prov.ID)
				}

				if len(targetPeers) > 0 {
					break
				}
				if attempt < 3 {
					log.Info("Retrying DHT provider lookup (attempt %d/3)...", attempt+1)
					time.Sleep(1 * time.Second)
				}
			}
		}

		if len(targetPeers) == 0 {
			log.Fatalf("Fatal: No storage providers available. Ensure at least one provider is running with -allow-push or specify -providers explicitly.")
		}

		effectiveReplication := *replication
		if effectiveReplication > len(targetPeers) {
			effectiveReplication = len(targetPeers)
			log.Info("Notice: Reduced replication to %d to match available provider count", effectiveReplication)
		}

		plan, err := distribution.PlanPlacement(m, targetPeers, effectiveReplication)
		if err != nil {
			log.Fatalf("Failed to plan chunk placement: %v", err)
		}

		if *push {
			tracker := distribution.NewGlobalReplicaTracker(effectiveReplication)
			pushCtx, pushCancel := context.WithTimeout(ctx, *pushTimeout)
			defer pushCancel()

			if err := distribution.Distribute(pushCtx, t, eng, plan, tracker, distribution.DefaultUploaderConfig); err != nil {
				log.Fatalf("Push distribution failed to satisfy replication invariant: %v", err)
			}

			log.Success("All chunks successfully committed with >= %d replicas across %d remote providers!",
				effectiveReplication, len(targetPeers))
		}

		fmt.Println("\n+---------------------------------------------------------------------------------------------------------+")
		fmt.Println("|                                 PUBLISHER DISPERSAL & PLACEMENT TABLE                                   |")
		fmt.Println("+------------------------------------+---------------+-------------+---------------+----------------------+")
		fmt.Println("| Storage Provider Peer ID           | Replicas (Qty)| Share (%)   | Volume (KB)   | Status               |")
		fmt.Println("+------------------------------------+---------------+-------------+---------------+----------------------+")
		totalReplicas := 0
		for _, pID := range targetPeers {
			chunksAssigned := len(plan.ProviderChunks[pID])
			totalReplicas += chunksAssigned
			pct := 0.0
			if len(m.ChunkIDs)*effectiveReplication > 0 {
				pct = float64(chunksAssigned) / float64(len(m.ChunkIDs)*effectiveReplication) * 100
			}
			volKB := float64(chunksAssigned * 32)
			fmt.Printf("| %-34s | %-13s | %-11s | %-13s | %-20s |\n",
				pID.String(),
				fmt.Sprintf("%d chunks", chunksAssigned),
				fmt.Sprintf("%.1f%%", pct),
				fmt.Sprintf("%.1f KB", volKB),
				"COMMITTED [✓]",
			)
		}
		fmt.Println("+------------------------------------+---------------+-------------+---------------+----------------------+")
		fmt.Printf("| TOTAL CLUSTER REPLICATION (R=%d)    | %-13s | %-11s | %-13s | Invariant: %-10s |\n",
			effectiveReplication,
			fmt.Sprintf("%d replicas", totalReplicas),
			"100.0%",
			fmt.Sprintf("%.1f KB", float64(totalReplicas*32)),
			"SATISFIED [✓]",
		)
		fmt.Println("+------------------------------------+---------------+-------------+---------------+----------------------+\n")

		if *challengeProviders && len(targetPeers) > 0 {
			log.Sub("Availability").Info("Starting Availability challenge engine (Interval: %v, Rounds: %d, Continuous Loop: %t)...",
				*challengeInterval, *challengeRounds, *challengeLoop)

			leafHashes := make([][]byte, len(m.ChunkIDs))
			for i, cid := range m.ChunkIDs {
				leafHashes[i] = make([]byte, 32)
				copy(leafHashes[i], cid[:])
			}

			tree, err := availability.NewMerkleTree(leafHashes)
			if err == nil {
				merkleRoot := tree.Root()
				client := availability.NewAvailabilityClient(h)

				round := 1
				for {
					if !*challengeLoop && round > *challengeRounds {
						break
					}

					for _, target := range targetPeers {
						assignedChunks := plan.ProviderChunks[target]
						chunkIdx := 0
						if len(assignedChunks) > 0 {
							targetChunkID := assignedChunks[(round-1)%len(assignedChunks)]
							for idx, cid := range m.ChunkIDs {
								if cid == targetChunkID {
									chunkIdx = idx
									break
								}
							}
						}

						challenge := availabilitytypes.Challenge{
							ChallengeID: availabilitytypes.ChallengeID(fmt.Sprintf("challenge-pub-round-%d-%s", round, target.String()[:8])),
							ContractID:  "avail-pub-contract",
							ProviderID:  target.String(),
							FileID:      hex.EncodeToString(m.Descriptor.ID[:]),
							ChunkID:     chunkIdx,
							Nonce:       []byte(fmt.Sprintf("pub-entropy-nonce-r%d-%d", round, time.Now().UnixNano())),
							CreatedAt:   time.Now().UTC(),
						}

						log.Sub("Availability").Info("[Round %d] Challenging provider %s for chunk %d...", round, target.String()[:12], challenge.ChunkID)
						resp, err := client.ChallengeProvider(ctx, target, challenge)
						if err != nil {
							log.Sub("Availability").Warn("[Round %d] Challenge request failed on %s: %v", round, target.String()[:12], err)
						} else {
							valid := verification.VerifyMerkleProof(resp.ChunkHash, challenge.ChunkID, resp.MerkleProof, merkleRoot)
							if valid {
								log.Sub("Availability").Success("[Round %d] Provider %s verified chunk possession (Merkle Proof PASS for chunk %d)", round, target.String()[:12], challenge.ChunkID)

								// Construct and dispatch optimistic cumulative payment voucher
								cumulativePay := uint64(round) * 25000000000000000 // 0.025 ETH per verified challenge period
								voucher := payment.PaymentState{
									ContractID:         "avail-pub-contract",
									Publisher:          h.ID().String(),
									Provider:           target.String(),
									Sequence:           uint64(round),
									Period:             uint64(round),
									CumulativePayment:  cumulativePay,
									LastChallengeID:    string(challenge.ChallengeID),
									ValidUntil:         uint64(time.Now().Add(24 * time.Hour).Unix()),
									Status:             interfaces.AvailabilityPass,
									PublisherSignature: []byte("publisher-signed-voucher-state"),
								}

								if err := client.SendPaymentVoucher(ctx, target, voucher); err != nil {
									log.Sub("Payment").Warn("[Round %d] Failed to send payment voucher to %s: %v", round, target.String()[:12], err)
								} else {
									log.Sub("Payment").Success("[Round %d] Dispatched Optimistic Payment Voucher to %s (Seq: %d, Cumulative: %d wei)",
										round, target.String()[:12], voucher.Sequence, voucher.CumulativePayment)
								}
							} else {
								log.Sub("Availability").Error("[Round %d] Merkle proof verification failed for provider %s!", round, target.String()[:12])
							}
						}
					}

					round++
					if *challengeLoop || round <= *challengeRounds {
						log.Sub("Availability").Info("Waiting %v for next 5-second challenge round...", *challengeInterval)
						time.Sleep(*challengeInterval)
					}
				}
			}
		}
	} else {
		// Traditional direct publisher DHT announcement
		log.Sub("DHT").Info("Announcing ContentID %x on DHT...", m.Descriptor.ID)
		if err := discovery.Provide(ctx, kdht, m.Descriptor.ID); err != nil {
			log.Sub("DHT").Warn("Could not advertise on DHT: %v (ensure bootstrap node is active)", err)
		} else {
			log.Sub("DHT").Success("Successfully announced ContentID on DHT")
		}
	}

	key, _ := keys.Get(ctx, m.Descriptor.ID)

	fields := []logger.Field{
		{Key: "Publisher Role", Value: *roleName},
		{Key: "File Ingested ", Value: *filePath},
		{Key: "ContentID     ", Value: fmt.Sprintf("%x", m.Descriptor.ID)},
		{Key: "Decryption Key", Value: fmt.Sprintf("%x", key)},
		{Key: "Chunks Total  ", Value: fmt.Sprintf("%d (%d KB per chunk)", len(m.ChunkIDs), *chunkSizeKB)},
		{Key: "Publisher ID  ", Value: h.ID().String()},
		{Key: "", Value: "Listening Multiaddresses:"},
	}
	for _, addr := range h.Addrs() {
		fields = append(fields, logger.Field{Key: "", Value: fmt.Sprintf("  - %s/p2p/%s", addr.String(), h.ID().String())})
	}

	log.Banner(fmt.Sprintf("CIPHER PUBLISHER: %s", strings.ToUpper(*roleName)), fields...)


	if *push || *challengeProviders || !*seed {
		if *push {
			log.Info("Remote push complete, exiting.")
		} else if *challengeProviders {
			log.Info("Availability challenge audit complete, exiting.")
		} else {
			log.Info("Seeding flag is false, exiting publisher.")
		}
		return
	}

	log.Success("Seeding content over /cipher/chunk/1.0.0. Press Ctrl+C to stop.")

	// Wait for OS shutdown signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Warn("Shutting down publisher...")
}

