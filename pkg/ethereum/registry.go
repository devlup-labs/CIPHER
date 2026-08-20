package ethereum

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
)

func RegisterProvider(h host.Host) error {

	// Load an existing Ethereum wallet or generate a new one.
	privateKey, err := LoadOrGenerateEthereumKey("provider_eth.key")
	if err != nil {
		return err
	}

	// Create a signed binding between the libp2p PeerID and
	// the Ethereum wallet.
	req, err := SignBinding(h.ID(), privateKey)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", req)

	// ------------------------------------------------------------------
	// MVP:
	// Verify the binding locally in Go.
	//
	// TODO:
	// Replace VerifyBinding() with a smart-contract call.
	// The contract will verify the signature and permanently
	// store the PeerID -> Wallet mapping.
	// ------------------------------------------------------------------
	if err := VerifyBinding(req); err != nil {
		return fmt.Errorf("binding verification failed: %w", err)
	}

	fmt.Println("Binding verified successfully!")

	// ------------------------------------------------------------------
	// MVP:
	// Register the provider in the in-memory registry.
	//
	// TODO:
	// Replace this with contract.registerProvider().
	// ------------------------------------------------------------------
	if err := Register(req.PeerID, req.Wallet); err != nil {
		return err
	}

	// ------------------------------------------------------------------
	// MVP:
	// Give the provider a hardcoded stake.
	//
	// TODO:
	// Replace this with:
	//
	//     contract.stake{value: ...}()
	//
	// so that the provider deposits real ETH.
	// ------------------------------------------------------------------
	if err := Stake(req.PeerID, DefaultStake); err != nil {
		return err
	}

	fmt.Printf("Provider registered with stake %d\n", DefaultStake)

	return nil
}