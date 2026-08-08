package manager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"cipher/internal/content/core"
	"cipher/internal/content/engine"
	"cipher/internal/transfer/scheduler"
	"cipher/internal/transport"
)

type Progress struct {
	TotalChunks     int
	CompletedChunks int
	BytesDownloaded uint64
	Started         time.Time
	ETA             time.Duration
}

type TransferManager struct {
	SessionManager SessionManager
	Engine         *engine.ContentEngine
	Transport      *transport.Transport
}

func NewTransferManager(sm SessionManager, eng *engine.ContentEngine, t *transport.Transport) *TransferManager {
	return &TransferManager{
		SessionManager: sm,
		Engine:         eng,
		Transport:      t,
	}
}

func (tm *TransferManager) Download(ctx context.Context, contentID core.ContentID, chunkIDs []core.ChunkID, peers []peer.ID) error {
	if len(peers) == 0 {
		return fmt.Errorf("no peers provided")
	}

	// 2. Setup Session
	sess, err := tm.SessionManager.Open(contentID)
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}

	if sess == nil {
		sess = &TransferSession{
			ContentID:   contentID,
			TargetPeer:  peers[0], // Historic artifact, keep for now or drop later
			Status:      StatusInProgress,
			StartedAt:   time.Now(),
			Completed:   make([]bool, len(chunkIDs)),
			TotalChunks: len(chunkIDs),
		}
		if err := tm.SessionManager.Save(sess); err != nil {
			return fmt.Errorf("failed to save new session: %w", err)
		}
	} else {
		// Session exists, ensure completed array matches manifest
		if len(sess.Completed) == 0 {
			sess.Completed = make([]bool, len(chunkIDs))
			sess.TotalChunks = len(chunkIDs)
		}
		sess.Status = StatusInProgress
		if err := tm.SessionManager.Save(sess); err != nil {
			return fmt.Errorf("failed to save resumed session: %w", err)
		}
	}

	// 3. Setup Scheduler tasks
	var tasks []scheduler.ChunkTask
	completedCount := 0
	for i, chunkID := range chunkIDs {
		if i < len(sess.Completed) && sess.Completed[i] {
			completedCount++
			continue
		}
		// Also double-check engine
		has, _ := tm.Engine.HasChunk(ctx, chunkID)
		if has {
			sess.Completed[i] = true
			completedCount++
			continue
		}
		tasks = append(tasks, scheduler.ChunkTask{
			Index:    i,
			ChunkID:  chunkID,
			Attempts: 0,
		})
	}
	tm.SessionManager.Save(sess)

	log.Printf("Starting download: %d chunks remaining", len(tasks))
	
	if len(tasks) == 0 {
		sess.Status = StatusCompleted
		tm.SessionManager.Save(sess)
		return nil
	}

	// 4. Setup Sources
	var sources []scheduler.Source
	for _, p := range peers {
		sources = append(sources, scheduler.Source{
			PeerID:    p,
			Available: nil, // Assume all chunks are available initially
		})
	}

	// 5. Run Scheduler
	sched := scheduler.NewScheduler(tm.Transport, tm.Engine, 3) // MaxAttempts = 3
	
	completions := make(chan scheduler.WorkerResult, len(tasks))
	errCh := make(chan error, 1)
	
	go func() {
		errCh <- sched.Run(ctx, tasks, sources, completions)
		close(completions)
	}()

	// Tracking metrics
	peerContributions := make(map[string]int)

	// 6. Handle Progress
	for res := range completions {
		if res.Task.Index < len(sess.Completed) {
			sess.Completed[res.Task.Index] = true
		}
		completedCount++
		
		if res.PeerID != "" {
			peerContributions[res.PeerID]++
		}

		if err := tm.SessionManager.Save(sess); err != nil {
			log.Printf("[TransferManager] Warning: failed to save session: %v", err)
		}

		fmt.Printf("\r\033[K[Progress] %d/%d chunks (%.1f%%)", completedCount, sess.TotalChunks, float64(completedCount)/float64(sess.TotalChunks)*100)
	}

	fmt.Println()
	if sess.TotalChunks > 0 {
		fmt.Println("\n--- Peer Contribution Metrics ---")
		for peerID, count := range peerContributions {
			fmt.Printf("Peer %s: %d chunks (%.1f%%)\n", peerID, count, float64(count)/float64(sess.TotalChunks)*100)
		}
		fmt.Println("---------------------------------")
	}

	schedErr := <-errCh
	if schedErr != nil {
		sess.Status = StatusFailed
		tm.SessionManager.Save(sess)
		return fmt.Errorf("transfer failed: %w", schedErr)
	}

	sess.Status = StatusCompleted
	tm.SessionManager.Save(sess)

	return nil
}
