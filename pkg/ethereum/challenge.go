package ethereum

import (
	"fmt"
	"sync"
	"time"
	"github.com/1amKhush/CIPHER/pkg/logger"
)

// Challenge represents an active challenge against a provider.
//
// MVP:
//   - Stored in Go memory.
//
// Future:
//   - This state will be stored in the Ethereum smart contract.
type Challenge struct {
	ProviderPeerID string
	ClientPeerID   string
	ChunkIndex     uint64

	// Time when the challenge was created.
	StartTime time.Time

	// Deadline before which the provider must respond.
	Deadline time.Time

	// Whether the challenge has been resolved.
	Resolved bool
}

var (
	challenges = make(map[string]*Challenge)
	challengeMu sync.RWMutex
)

// CreateChallenge creates a new challenge for a provider.
//
// MVP:
//   - Stores the challenge in memory.
//   - Starts a countdown timer.
//
// Future:
//   - Replace this with a smart contract transaction:
//
//       challenge(provider)
//
//   - The contract will emit a ChallengeCreated event
//     and store the deadline on-chain.
func CreateChallenge(providerPeerID string, timeout time.Duration) error {
	challengeMu.Lock()
	defer challengeMu.Unlock()

	if _, exists := challenges[providerPeerID]; exists {
		return fmt.Errorf("challenge already exists")
	}

	now := time.Now()

	logger.Info().
    Str("provider", providerPeerID).
    Dur("timeout", timeout).
    Msg("Challenge created")

	// MVP:
	// Start a local timer to expire the challenge.
	//
	// Future:
	// The Ethereum smart contract will enforce deadlines.
	// This goroutine will be removed.
	go func(provider string, timeout time.Duration) {
		time.Sleep(timeout)
		ExpireChallenges()
	}(providerPeerID, timeout)

	challenges[providerPeerID] = &Challenge{
		ProviderPeerID: providerPeerID,
		StartTime:      now,
		Deadline:       now.Add(timeout),
		Resolved:       false,
	}

	return nil
}

// ResolveChallenge marks a challenge as successfully resolved.
//
// MVP:
//   - Checks that the provider responded before the deadline.
//   - Removes the challenge from memory.
//
// Future:
//   - Replace this with:
//
//       submitKey(...)
//
//
//   - The smart contract will verify the revealed key
//     and close the challenge.
func ResolveChallenge(providerPeerID string) error {
	challengeMu.Lock()
	defer challengeMu.Unlock()

	challenge, ok := challenges[providerPeerID]
	if !ok {
		return fmt.Errorf("challenge not found")
	}

	if time.Now().After(challenge.Deadline) {
		return fmt.Errorf("challenge deadline already expired")
	}

	challenge.Resolved = true

	logger.Info().
    Str("provider", providerPeerID).
    Msg("Challenge resolved")

	delete(challenges, providerPeerID)

	return nil
}

// ExpireChallenges checks all active challenges.
//
// MVP:
//   - Any unresolved challenge past its deadline
//     causes the provider's stake to be slashed.
//
// Future:
//   - This logic will be executed by the smart contract,
//     where anyone can trigger expiration after the deadline.
func ExpireChallenges() {
	challengeMu.Lock()
	defer challengeMu.Unlock()

	now := time.Now()

	for peerID, challenge := range challenges {

		if challenge.Resolved {
			delete(challenges, peerID)
			continue
		}

		if now.After(challenge.Deadline) {

			// Slash 10 stake units for MVP.
			// Future:
			// Replace with:
			//
			//     contract.slashStake(peerID)
			//
			// executed by the Ethereum smart contract.
			logger.Warn().
			Str("provider", peerID).
			Uint64("slashed", 10).
			Msg("Stake slashed")
			_ = SlashStake(peerID, 10)

			delete(challenges, peerID)
		}
	}
}

// GetChallenge returns an active challenge.
//
// MVP:
//   - Reads from the in-memory registry.
//
// Future:
//   - Replace with a smart contract query or event lookup.
func GetChallenge(providerPeerID string) (*Challenge, error) {
	challengeMu.RLock()
	defer challengeMu.RUnlock()

	challenge, ok := challenges[providerPeerID]
	if !ok {
		return nil, fmt.Errorf("challenge not found")
	}

	copy := *challenge
	return &copy, nil
}