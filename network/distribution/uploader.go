package distribution

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"cipher/network/content/core"
	"cipher/network/content/engine"
	"cipher/network/protocol/push"
	"cipher/network/transport"
)

type UploaderConfig struct {
	MaxRetriesPerChunk int
	FailoverRounds     int
}

var DefaultUploaderConfig = UploaderConfig{
	MaxRetriesPerChunk: 3,
	FailoverRounds:     2,
}

// Distribute pushes all assigned chunks to target providers according to the placement plan,
// enforcing that every chunk reaches >= plan.Replication committed replicas.
func Distribute(
	ctx context.Context,
	t *transport.Transport,
	eng *engine.ContentEngine,
	plan *PlacementPlan,
	tracker *GlobalReplicaTracker,
	cfg UploaderConfig,
) error {
	manifestBytes, err := eng.GetManifestBytes(ctx, plan.ContentID)
	if err != nil {
		return fmt.Errorf("failed to get manifest bytes: %w", err)
	}

	allChunks := plan.Manifest.ChunkIDs

	// Initialize tracker state for initial plan
	for chunkID, providers := range plan.Assignments {
		for _, p := range providers {
			tracker.SetStatus(chunkID, p, ReplicaPending)
		}
	}

	log.Printf("[Distribution] Beginning upload for ContentID %x across %d providers (Replication R=%d)...",
		plan.ContentID, len(plan.ProviderChunks), plan.Replication)

	// Phase 1: Upload initial assignments in parallel across providers
	var wg sync.WaitGroup
	var activeProviders []peer.ID

	for p, assigned := range plan.ProviderChunks {
		if len(assigned) == 0 {
			continue
		}
		activeProviders = append(activeProviders, p)
		wg.Add(1)

		go func(targetPeer peer.ID, chunkList []core.ChunkID) {
			defer wg.Done()
			uploadToProvider(ctx, t, eng, plan.ContentID, targetPeer, chunkList, manifestBytes, tracker, cfg.MaxRetriesPerChunk)
		}(p, assigned)
	}

	wg.Wait()

	// Phase 2: Check invariant and perform failover reassignments if needed
	satisfied, total := tracker.GetSummary(allChunks)
	log.Printf("[Distribution] Initial round completed: %d/%d chunks satisfied replication >= %d",
		satisfied, total, plan.Replication)

	if tracker.IsComplete(allChunks) {
		log.Printf("[Distribution] [✓] Invariant satisfied! All %d chunks have >= %d committed replicas.", total, plan.Replication)
		return nil
	}

	// Attempt failover reassignments
	for round := 1; round <= cfg.FailoverRounds && !tracker.IsComplete(allChunks); round++ {
		log.Printf("[Distribution] Starting failover recovery round %d...", round)

		// Group missing chunks by target candidate providers
		reassignments := make(map[peer.ID][]core.ChunkID)
		for _, cid := range allChunks {
			needed := plan.Replication - tracker.CommittedCount(cid)
			if needed <= 0 {
				continue
			}

			// Find providers that don't have this chunk and haven't failed on it
			for _, p := range activeProviders {
				if needed <= 0 {
					break
				}
				tracker.mu.RLock()
				status := tracker.state[cid][p]
				tracker.mu.RUnlock()

				if status != ReplicaCommitted && status != ReplicaUploading {
					reassignments[p] = append(reassignments[p], cid)
					tracker.SetStatus(cid, p, ReplicaPending)
					needed--
				}
			}
		}

		if len(reassignments) == 0 {
			log.Printf("[Distribution] No viable alternate providers available for failover")
			break
		}

		var failoverWg sync.WaitGroup
		for p, chunks := range reassignments {
			if len(chunks) == 0 {
				continue
			}
			failoverWg.Add(1)
			go func(targetPeer peer.ID, chunkList []core.ChunkID) {
				defer failoverWg.Done()
				uploadToProvider(ctx, t, eng, plan.ContentID, targetPeer, chunkList, manifestBytes, tracker, cfg.MaxRetriesPerChunk)
			}(p, chunks)
		}
		failoverWg.Wait()
	}

	// Final verification of the invariant
	if !tracker.IsComplete(allChunks) {
		satisfied, total := tracker.GetSummary(allChunks)
		return fmt.Errorf("replication invariant violation: only %d/%d chunks reached required %d replicas",
			satisfied, total, plan.Replication)
	}

	log.Printf("[Distribution] [✓] Invariant successfully satisfied: all %d chunks have >= %d committed replicas.",
		len(allChunks), plan.Replication)
	return nil
}

func uploadToProvider(
	ctx context.Context,
	t *transport.Transport,
	eng *engine.ContentEngine,
	contentID core.ContentID,
	targetPeer peer.ID,
	chunks []core.ChunkID,
	manifestBytes []byte,
	tracker *GlobalReplicaTracker,
	maxRetries int,
) {
	log.Printf("[Distribution] Connecting push stream to provider %s (%d chunks assigned)...", targetPeer, len(chunks))

	client, err := push.NewClient(ctx, t, targetPeer)
	if err != nil {
		log.Printf("[Distribution] Failed to connect to provider %s: %v", targetPeer, err)
		for _, cid := range chunks {
			tracker.SetStatus(cid, targetPeer, ReplicaFailed)
		}
		return
	}
	defer client.Close()

	// 1. Send manifest with assigned chunk list
	if err := client.SendManifest(ctx, contentID, chunks, manifestBytes); err != nil {
		log.Printf("[Distribution] Provider %s rejected manifest: %v", targetPeer, err)
		for _, cid := range chunks {
			tracker.SetStatus(cid, targetPeer, ReplicaFailed)
		}
		return
	}

	committedInBatch := 0

	// 2. Stream chunks with per-chunk retry
	for _, chunkID := range chunks {
		chunkData, err := eng.GetChunk(ctx, chunkID)
		if err != nil {
			log.Printf("[Distribution] Failed to read chunk %x from local engine: %v", chunkID, err)
			tracker.SetStatus(chunkID, targetPeer, ReplicaFailed)
			continue
		}

		tracker.SetStatus(chunkID, targetPeer, ReplicaUploading)

		var sendErr error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			sendErr = client.SendChunk(ctx, contentID, chunkData)
			if sendErr == nil {
				break
			}
			time.Sleep(100 * time.Millisecond * time.Duration(attempt))
		}

		if sendErr != nil {
			log.Printf("[Distribution] Failed to upload chunk %x to %s after %d attempts: %v",
				chunkID, targetPeer, maxRetries, sendErr)
			tracker.SetStatus(chunkID, targetPeer, ReplicaFailed)
		} else {
			tracker.SetStatus(chunkID, targetPeer, ReplicaCommitted)
			committedInBatch++
		}
	}

	// 3. Finalize batch if all assigned chunks succeeded
	if committedInBatch == len(chunks) {
		if err := client.SendBatchComplete(ctx, contentID); err != nil {
			log.Printf("[Distribution] BatchComplete failed for provider %s: %v", targetPeer, err)
		} else {
			log.Printf("[Distribution] [✓] Successfully completed batch and triggered DHT announcement on %s", targetPeer)
		}
	} else {
		log.Printf("[Distribution] Incomplete batch on %s (%d/%d committed). Omitting BatchComplete.",
			targetPeer, committedInBatch, len(chunks))
	}
}
