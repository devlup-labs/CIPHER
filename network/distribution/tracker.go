package distribution

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"cipher/network/content/core"
)

type ReplicaState uint8

const (
	ReplicaPending ReplicaState = iota
	ReplicaUploading
	ReplicaCommitted
	ReplicaFailed
)

// GlobalReplicaTracker tracks chunk replicas across all providers to guarantee the R-replication invariant.
type GlobalReplicaTracker struct {
	mu        sync.RWMutex
	state     map[core.ChunkID]map[peer.ID]ReplicaState
	requiredR int
}

func NewGlobalReplicaTracker(requiredR int) *GlobalReplicaTracker {
	return &GlobalReplicaTracker{
		state:     make(map[core.ChunkID]map[peer.ID]ReplicaState),
		requiredR: requiredR,
	}
}

func (t *GlobalReplicaTracker) SetStatus(chunkID core.ChunkID, p peer.ID, status ReplicaState) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state[chunkID] == nil {
		t.state[chunkID] = make(map[peer.ID]ReplicaState)
	}
	t.state[chunkID][p] = status
}

func (t *GlobalReplicaTracker) CommittedCount(chunkID core.ChunkID) int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	count := 0
	for _, status := range t.state[chunkID] {
		if status == ReplicaCommitted {
			count++
		}
	}
	return count
}

func (t *GlobalReplicaTracker) IsChunkSatisfied(chunkID core.ChunkID) bool {
	return t.CommittedCount(chunkID) >= t.requiredR
}

// IsComplete returns true if and only if EVERY chunk in chunkIDs has >= requiredR committed replicas.
func (t *GlobalReplicaTracker) IsComplete(chunkIDs []core.ChunkID) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, id := range chunkIDs {
		committedCount := 0
		for _, status := range t.state[id] {
			if status == ReplicaCommitted {
				committedCount++
			}
		}
		if committedCount < t.requiredR {
			return false
		}
	}
	return true
}

func (t *GlobalReplicaTracker) GetSummary(chunkIDs []core.ChunkID) (satisfied int, total int) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	total = len(chunkIDs)
	for _, id := range chunkIDs {
		count := 0
		for _, status := range t.state[id] {
			if status == ReplicaCommitted {
				count++
			}
		}
		if count >= t.requiredR {
			satisfied++
		}
	}
	return satisfied, total
}
