package distribution

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"cipher/network/content/core"
	"cipher/network/content/manifest"
)

// PlacementPlan contains the deterministic chunk-to-provider assignments.
type PlacementPlan struct {
	ContentID      core.ContentID
	Manifest       *manifest.Manifest
	Assignments    map[core.ChunkID][]peer.ID
	ProviderChunks map[peer.ID][]core.ChunkID
	Replication    int
}

// PlanPlacement calculates circular redundant chunk assignments for N providers with replication R.
func PlanPlacement(m *manifest.Manifest, providers []peer.ID, replication int) (*PlacementPlan, error) {
	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers provided for placement")
	}
	if replication <= 0 {
		return nil, fmt.Errorf("replication factor must be greater than zero, got %d", replication)
	}
	if replication > len(providers) {
		return nil, fmt.Errorf("replication factor %d cannot exceed provider count %d", replication, len(providers))
	}
	if len(m.ChunkIDs) == 0 {
		return nil, fmt.Errorf("manifest contains zero chunks")
	}

	plan := &PlacementPlan{
		ContentID:      m.Descriptor.ID,
		Manifest:       m,
		Assignments:    make(map[core.ChunkID][]peer.ID),
		ProviderChunks: make(map[peer.ID][]core.ChunkID),
		Replication:    replication,
	}

	for _, p := range providers {
		plan.ProviderChunks[p] = make([]core.ChunkID, 0)
	}

	n := len(providers)
	for i, chunkID := range m.ChunkIDs {
		plan.Assignments[chunkID] = make([]peer.ID, 0, replication)
		for r := 0; r < replication; r++ {
			targetIndex := (i + r) % n
			p := providers[targetIndex]
			plan.Assignments[chunkID] = append(plan.Assignments[chunkID], p)
			plan.ProviderChunks[p] = append(plan.ProviderChunks[p], chunkID)
		}
	}

	return plan, nil
}
