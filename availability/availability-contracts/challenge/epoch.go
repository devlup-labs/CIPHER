package challenge

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	availabilitytypes "cipher/availability/availability-contracts/types"
)

var epochStore = struct {
	sync.RWMutex
	epochs map[availabilitytypes.EpochID]*availabilitytypes.Epoch
	latest map[string]availabilitytypes.EpochID
}{
	epochs: make(map[availabilitytypes.EpochID]*availabilitytypes.Epoch),
	latest: make(map[string]availabilitytypes.EpochID),
}

// CreateEpoch starts a fresh epoch with every chunk eligible for selection.
func CreateEpoch(providerID, fileID string, totalChunks int) (availabilitytypes.EpochID, error) {
	if providerID == "" || fileID == "" {
		return "", errors.New("provider ID and file ID are required")
	}
	if totalChunks <= 0 {
		return "", errors.New("total chunks must be positive")
	}
	randomID := make([]byte, 16)
	if _, err := rand.Read(randomID); err != nil {
		return "", fmt.Errorf("generate epoch ID: %w", err)
	}
	epochID := availabilitytypes.EpochID("epoch-" + hex.EncodeToString(randomID))
	uncovered := make([]int, totalChunks)
	for index := range uncovered {
		uncovered[index] = index
	}

	epochStore.Lock()
	defer epochStore.Unlock()
	key := providerID + "\x00" + fileID
	roundNumber := uint64(0)
	if previousID, ok := epochStore.latest[key]; ok {
		if previous, exists := epochStore.epochs[previousID]; exists {
			previous.Status = availabilitytypes.EpochEnded
			roundNumber = previous.RoundNumber + 1
		}
	}
	epochStore.epochs[epochID] = &availabilitytypes.Epoch{
		ID: epochID, ProviderID: providerID, FileID: fileID, TotalChunks: totalChunks,
		Uncovered: uncovered, Status: availabilitytypes.EpochActive, RoundNumber: roundNumber,
		CreatedAt: time.Now().UTC(),
	}
	epochStore.latest[key] = epochID
	return epochID, nil
}

// SelectUncoveredChunk chooses and removes one currently uncovered chunk.
func SelectUncoveredChunk(epochID availabilitytypes.EpochID, chunkHash []byte) (int, error) {
	if epochID == "" || len(chunkHash) == 0 {
		return 0, errors.New("epoch ID and chunk hash are required")
	}
	epochStore.Lock()
	defer epochStore.Unlock()
	epoch, ok := epochStore.epochs[epochID]
	if !ok {
		return 0, fmt.Errorf("unknown epoch: %s", epochID)
	}
	if epoch.Status != availabilitytypes.EpochActive {
		return 0, fmt.Errorf("epoch %s is not active", epochID)
	}
	if len(epoch.Uncovered) == 0 {
		return 0, errors.New("no uncovered chunks remain in epoch")
	}
	index := new(big.Int).SetBytes(chunkHash)
	index.Mod(index, big.NewInt(int64(len(epoch.Uncovered))))
	selectedIndex := int(index.Int64())
	chunkID := epoch.Uncovered[selectedIndex]
	epoch.Uncovered = append(epoch.Uncovered[:selectedIndex], epoch.Uncovered[selectedIndex+1:]...)
	epoch.ChallengeCounter++
	return chunkID, nil
}

// CheckRoundTrigger applies the finalized one-in-three round termination rule.
func CheckRoundTrigger(roundHash []byte) bool {
	if len(roundHash) == 0 {
		return false
	}
	value := new(big.Int).SetBytes(roundHash)
	return value.Mod(value, big.NewInt(3)).Sign() == 0
}

// StartNextEpoch ends the current epoch for a provider/file and resets every
// chunk as eligible in the new epoch.
func StartNextEpoch(providerID, fileID string) (availabilitytypes.EpochID, error) {
	if providerID == "" || fileID == "" {
		return "", errors.New("provider ID and file ID are required")
	}
	epochStore.RLock()
	previousID, ok := epochStore.latest[providerID+"\x00"+fileID]
	if !ok {
		epochStore.RUnlock()
		return "", errors.New("no prior epoch exists for provider and file")
	}
	previous, ok := epochStore.epochs[previousID]
	if !ok {
		epochStore.RUnlock()
		return "", errors.New("latest epoch is unavailable")
	}
	totalChunks := previous.TotalChunks
	epochStore.RUnlock()
	return CreateEpoch(providerID, fileID, totalChunks)
}

// GetEpochState returns the externally relevant state without exposing the
// mutable uncovered-chunk collection.
func GetEpochState(epochID availabilitytypes.EpochID) (availabilitytypes.EpochState, error) {
	if epochID == "" {
		return availabilitytypes.EpochState{}, errors.New("epoch ID is required")
	}
	epochStore.RLock()
	defer epochStore.RUnlock()
	epoch, ok := epochStore.epochs[epochID]
	if !ok {
		return availabilitytypes.EpochState{}, fmt.Errorf("unknown epoch: %s", epochID)
	}
	return availabilitytypes.EpochState{EpochID: epoch.ID, ChallengeCounter: epoch.ChallengeCounter, UncoveredChunkCount: len(epoch.Uncovered), EpochStatus: epoch.Status, RoundNumber: epoch.RoundNumber}, nil
}
