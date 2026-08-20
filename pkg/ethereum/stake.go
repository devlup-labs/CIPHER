package ethereum

// MVP:
// Challenges and stakes are stored in local memory for demonstration.
//
// Future:
// The challenge, deadline, and stake will be maintained by an Ethereum
// smart contract so both client and provider observe the same state.

import (
	"fmt"
	"sync"
)
// DefaultStake is the amount assigned to a provider during MVP registration.
//
// MVP:
//   - This is a hardcoded value because we are not yet interacting with
//     an Ethereum smart contract.
//
// Future:
//   - Remove this constant.
//   - The provider will deposit real ETH by calling:
//
//       stake() payable
//
//     on the smart contract, and the deposited amount will become the
//     provider's actual stake.
const DefaultStake uint64 = 100

// Provider represents a registered storage provider.
//
// MVP:
//   - This information is stored only in Go memory.
//
// Future:
//   - This state will live inside the Ethereum smart contract.
type Provider struct {
	PeerID string
	Wallet string

	// Amount currently staked by the provider.
	Stake uint64

	// Whether the provider is allowed to serve data.
	Active bool
}

var (
	providers = make(map[string]*Provider)
	mu        sync.RWMutex
)

// RegisterProvider registers a provider after its PeerID has been
// successfully bound to an Ethereum wallet.
//
// MVP:
//   - Stores the provider in an in-memory registry.
//
// Future:
//   - Replace this with a smart contract transaction that permanently
//     registers the provider on-chain.
func Register(peerID, wallet string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := providers[peerID]; exists {
		return fmt.Errorf("provider already registered")
	}

	providers[peerID] = &Provider{
		PeerID: peerID,
		Wallet: wallet,
		Stake:  0,
		Active: false,
	}

	return nil
}

// Stake locks a provider's stake.
//
// MVP:
//   - Simply stores the stake amount in memory.
//   - No Ethereum transaction is performed.
//
// Future:
//   - Replace this implementation with:
//
//       stake() payable
//
//     on the Ethereum smart contract.
//
//   - The smart contract will hold the provider's deposited ETH and
//     use it as collateral against malicious behaviour.
func Stake(peerID string, amount uint64) error {
	mu.Lock()
	defer mu.Unlock()

	provider, ok := providers[peerID]
	if !ok {
		return fmt.Errorf("provider not registered")
	}

	provider.Stake = amount
	provider.Active = true

	return nil
}

// SlashStake penalizes a provider by reducing its stake.
//
// MVP:
//   - Reduces the in-memory stake.
//   - Used only to simulate the slashing flow.
//
// Future:
//   - Replace this with a smart contract transaction.
//   - The contract will slash a percentage of the provider's locked ETH.
func SlashStake(peerID string, amount uint64) error {
	mu.Lock()
	defer mu.Unlock()

	provider, ok := providers[peerID]
	if !ok {
		return fmt.Errorf("provider not registered")
	}

	if amount > provider.Stake {
		amount = provider.Stake
	}

	provider.Stake -= amount

	if provider.Stake == 0 {
		provider.Active = false
	}

	return nil
}

// GetProvider returns a copy of a registered provider.
//
// MVP:
//   - Reads from the in-memory registry.
//
// Future:
//   - Replace this with a smart contract query or event lookup.
func GetProvider(peerID string) (*Provider, error) {
	mu.RLock()
	defer mu.RUnlock()

	provider, ok := providers[peerID]
	if !ok {
		return nil, fmt.Errorf("provider not found")
	}

	// Return a copy so callers cannot modify internal state.
	copy := *provider
	return &copy, nil
}