package ethereum

import (
	"encoding/hex"
	"fmt"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// VerifyBinding verifies the cryptographic binding between a libp2p PeerID
// and an Ethereum wallet.
//
// NOTE:
// This is an OFF-CHAIN verifier used for development and testing.
// In the production protocol, this verification should be performed by the
// Ethereum smart contract (using ecrecover or an equivalent mechanism)
// during provider registration.
//
// Once the smart contract is integrated, this function can either:
//   1. remain as a local pre-check for testing/debugging, or
//   2. be replaced by a transaction that calls registerProvider()
//
// smart contract verifies,validates,stores PeerID <--> Wallet mapping and read it when needed

func VerifyBinding(req *BindingRequest) error {

	peerID, err := peer.Decode(req.PeerID)
	if err != nil {
		return fmt.Errorf("invalid peer id")
	}

	message := BuildMessage(peerID, req.Timestamp)

	hash := gethcrypto.Keccak256Hash(message)

	//Recover the signer from the signature.

	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature")
	}

	pubKey, err := gethcrypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return fmt.Errorf("signature verification failed")
	}

	address := gethcrypto.PubkeyToAddress(*pubKey)

	if address.Hex() != req.Wallet {
		return fmt.Errorf("wallet mismatch")
	}

	return nil
}