package distribution

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"

	"cipher/network/content/core"
	"cipher/network/content/manifest"
)

func generateTestPeerIDs(n int) []peer.ID {
	peers := make([]peer.ID, n)
	for i := 0; i < n; i++ {
		peers[i] = peer.ID(fmt.Sprintf("peer-%d", i))
	}
	return peers
}

func generateTestManifest(chunkCount int) *manifest.Manifest {
	chunkIDs := make([]core.ChunkID, chunkCount)
	for i := 0; i < chunkCount; i++ {
		_, _ = rand.Read(chunkIDs[i][:])
	}

	var cid core.ContentID
	_, _ = rand.Read(cid[:])

	return &manifest.Manifest{
		Descriptor: manifest.ContentDescriptor{
			ID: cid,
		},
		ChunkIDs: chunkIDs,
	}
}

func TestPlanPlacement(t *testing.T) {
	providers := generateTestPeerIDs(3)
	m := generateTestManifest(10)
	replication := 2

	plan, err := PlanPlacement(m, providers, replication)
	if err != nil {
		t.Fatalf("PlanPlacement failed: %v", err)
	}

	if plan.Replication != 2 {
		t.Errorf("expected replication 2, got %d", plan.Replication)
	}

	// 1. Every chunk must have exactly R distinct providers
	for _, chunkID := range m.ChunkIDs {
		assigned, ok := plan.Assignments[chunkID]
		if !ok {
			t.Fatalf("chunk %x not found in assignments", chunkID)
		}
		if len(assigned) != replication {
			t.Errorf("chunk %x has %d replicas, want %d", chunkID, len(assigned), replication)
		}
		if assigned[0] == assigned[1] {
			t.Errorf("chunk %x assigned to duplicate providers: %s", chunkID, assigned[0])
		}
	}

	// 2. Verify total chunk assignments equals chunkCount * replication
	totalAssigned := 0
	for _, chunks := range plan.ProviderChunks {
		totalAssigned += len(chunks)
	}
	if totalAssigned != len(m.ChunkIDs)*replication {
		t.Errorf("total chunk assignments %d, want %d", totalAssigned, len(m.ChunkIDs)*replication)
	}
}

func TestPlanPlacementValidation(t *testing.T) {
	m := generateTestManifest(5)
	providers := generateTestPeerIDs(2)

	// Replication > providers
	_, err := PlanPlacement(m, providers, 3)
	if err == nil {
		t.Errorf("expected error when replication > providers, got nil")
	}

	// Zero providers
	_, err = PlanPlacement(m, nil, 1)
	if err == nil {
		t.Errorf("expected error when providers is empty, got nil")
	}

	// Zero replication
	_, err = PlanPlacement(m, providers, 0)
	if err == nil {
		t.Errorf("expected error when replication <= 0, got nil")
	}
}
