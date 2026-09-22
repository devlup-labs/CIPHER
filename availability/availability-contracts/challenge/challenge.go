// Package challenge implements epoch-based probabilistic chunk selection.
package challenge

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	availabilitytypes "cipher/availability/availability-contracts/types"
)

// CreateChallenge creates an epoch for the file, then generates its first
// challenge. It only chooses a chunk; it does not send or verify a challenge.
func CreateChallenge(contractID, providerID, fileID string, totalChunks int, roundNumber uint64) (availabilitytypes.Challenge, error) {
	if contractID == "" || providerID == "" || fileID == "" {
		return availabilitytypes.Challenge{}, errors.New("contract ID, provider ID, and file ID are required")
	}
	epochID, err := CreateEpoch(providerID, fileID, totalChunks)
	if err != nil {
		return availabilitytypes.Challenge{}, fmt.Errorf("create epoch: %w", err)
	}
	return GenerateNextChallenge(contractID, providerID, fileID, epochID, 0, roundNumber)
}

// GenerateNextChallenge produces one challenge within an existing epoch using
// the finalized nonce, hash, selection, and round-trigger sequence.
func GenerateNextChallenge(contractID, providerID, fileID string, epochID availabilitytypes.EpochID, challengeCounter, roundNumber uint64) (availabilitytypes.Challenge, error) {
	if contractID == "" || providerID == "" || fileID == "" || epochID == "" {
		return availabilitytypes.Challenge{}, errors.New("contract ID, provider ID, file ID, and epoch ID are required")
	}
	epochState, err := GetEpochState(epochID)
	if err != nil {
		return availabilitytypes.Challenge{}, fmt.Errorf("get epoch state: %w", err)
	}
	if epochState.ChallengeCounter != challengeCounter {
		return availabilitytypes.Challenge{}, fmt.Errorf("challenge counter %d does not match epoch counter %d", challengeCounter, epochState.ChallengeCounter)
	}
	nonce, err := GenerateNonce(epochID, challengeCounter)
	if err != nil {
		return availabilitytypes.Challenge{}, fmt.Errorf("generate nonce: %w", err)
	}
	chunkHash := GenerateChunkHash(providerID, fileID, epochID, challengeCounter, nonce)
	chunkID, err := SelectUncoveredChunk(epochID, chunkHash)
	if err != nil {
		return availabilitytypes.Challenge{}, fmt.Errorf("select uncovered chunk: %w", err)
	}
	roundHash := GenerateRoundHash(providerID, fileID, epochID, challengeCounter, nonce, roundNumber)

	return availabilitytypes.Challenge{
		ChallengeID:      challengeID(contractID, epochID, challengeCounter, nonce),
		ContractID:       availabilitytypes.ContractID(contractID),
		ProviderID:       providerID,
		FileID:           fileID,
		EpochID:          epochID,
		ChunkID:          chunkID,
		Nonce:            append([]byte(nil), nonce...),
		ChallengeCounter: challengeCounter,
		RoundNumber:      roundNumber,
		TriggerStatus:    CheckRoundTrigger(roundHash),
		CreatedAt:        time.Now().UTC(),
	}, nil
}

func challengeID(contractID string, epochID availabilitytypes.EpochID, counter uint64, nonce []byte) availabilitytypes.ChallengeID {
	input := fmt.Sprintf("%s|%s|%d|", contractID, epochID, counter)
	sum := sha256.Sum256(append([]byte(input), nonce...))
	return availabilitytypes.ChallengeID("challenge-" + hex.EncodeToString(sum[:]))
}
