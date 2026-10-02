package distribution

import (
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"

	"cipher/network/content/core"
)

func TestGlobalReplicaTracker(t *testing.T) {
	chunk1 := core.ChunkID{1}
	chunk2 := core.ChunkID{2}
	p1 := peer.ID("peer-1")
	p2 := peer.ID("peer-2")
	p3 := peer.ID("peer-3")

	tracker := NewGlobalReplicaTracker(2)

	// Initially not complete
	chunks := []core.ChunkID{chunk1, chunk2}
	if tracker.IsComplete(chunks) {
		t.Fatalf("expected incomplete on empty tracker")
	}

	// Commit chunk1 to p1
	tracker.SetStatus(chunk1, p1, ReplicaCommitted)
	if tracker.IsChunkSatisfied(chunk1) {
		t.Fatalf("chunk1 should need 2 replicas, has 1")
	}

	// Commit chunk1 to p2 -> chunk1 satisfied
	tracker.SetStatus(chunk1, p2, ReplicaCommitted)
	if !tracker.IsChunkSatisfied(chunk1) {
		t.Fatalf("chunk1 should be satisfied with 2 replicas")
	}

	// But overall not complete because chunk2 has 0 replicas
	if tracker.IsComplete(chunks) {
		t.Fatalf("overall tracker should not be complete without chunk2")
	}

	// Commit chunk2 to p2 and p3
	tracker.SetStatus(chunk2, p2, ReplicaCommitted)
	tracker.SetStatus(chunk2, p3, ReplicaCommitted)

	// Now both satisfied
	if !tracker.IsComplete(chunks) {
		t.Fatalf("tracker should be complete when all chunks have >= 2 replicas")
	}

	satisfied, total := tracker.GetSummary(chunks)
	if satisfied != 2 || total != 2 {
		t.Errorf("summary mismatch: satisfied=%d, total=%d", satisfied, total)
	}
}

func TestTrackerRandomChunks(t *testing.T) {
	chunks := make([]core.ChunkID, 20)
	for i := range chunks {
		_, _ = rand.Read(chunks[i][:])
	}

	tracker := NewGlobalReplicaTracker(2)
	p1 := peer.ID("peer-a")
	p2 := peer.ID("peer-b")

	for _, c := range chunks {
		tracker.SetStatus(c, p1, ReplicaCommitted)
		tracker.SetStatus(c, p2, ReplicaCommitted)
	}

	if !tracker.IsComplete(chunks) {
		t.Errorf("expected all 20 chunks to be complete with 2 replicas")
	}
}
