package challenge

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	availabilitytypes "cipher/availability/availability-contracts/types"
)

// GenerateNonce returns a cryptographically secure nonce for one epoch and
// challenge counter. The identifiers are validated here but do not weaken the
// entropy supplied by crypto/rand.
func GenerateNonce(epochID availabilitytypes.EpochID, challengeCounter uint64) ([]byte, error) {
	if epochID == "" {
		return nil, errors.New("epoch ID is required")
	}
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate nonce for epoch %s counter %d: %w", epochID, challengeCounter, err)
	}
	return nonce, nil
}

// GenerateChunkHash derives the deterministic hash used for modulo chunk
// selection: SHA256(provider || file || epoch || counter || nonce || "CHUNK").
func GenerateChunkHash(providerID, fileID string, epochID availabilitytypes.EpochID, challengeCounter uint64, nonce []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(providerID))
	hash.Write([]byte(fileID))
	hash.Write([]byte(epochID))
	counter := make([]byte, 8)
	binary.BigEndian.PutUint64(counter, challengeCounter)
	hash.Write(counter)
	hash.Write(nonce)
	hash.Write([]byte("CHUNK"))
	return hash.Sum(nil)
}

// GenerateRoundHash derives the independent hash used exclusively for the
// one-in-three round-trigger decision.
func GenerateRoundHash(providerID, fileID string, epochID availabilitytypes.EpochID, challengeCounter uint64, nonce []byte, roundNumber uint64) []byte {
	hash := sha256.New()
	hash.Write([]byte(providerID))
	hash.Write([]byte(fileID))
	hash.Write([]byte(epochID))
	counter := make([]byte, 8)
	binary.BigEndian.PutUint64(counter, challengeCounter)
	hash.Write(counter)
	hash.Write(nonce)
	round := make([]byte, 8)
	binary.BigEndian.PutUint64(round, roundNumber)
	hash.Write(round)
	hash.Write([]byte("ROUND"))
	return hash.Sum(nil)
}
