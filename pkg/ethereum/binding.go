package ethereum

/*
Package ethereum implements OFF-CHAIN identity binding between a libp2p PeerID
and an Ethereum wallet.

Current status:
- Generates Ethereum signatures.
- Verifies signatures locally in Go.
- Used for development and unit testing.

Future work:
- Implement registerProvider() in the smart contract.
- Move signature verification on-chain.
- Persist PeerID -> Wallet mapping on-chain.
- Replace local verification with blockchain transactions.
*/

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"time"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

type BindingRequest struct {
	PeerID    string `json:"peer_id"`
	Wallet    string `json:"wallet"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// BuildMessage creates the canonical message
// Both the signer and verifier MUST construct exactly the same message.
func BuildMessage(peerID peer.ID, timestamp int64) []byte {
	return []byte(fmt.Sprintf(
		"PeerID:%s\nTimestamp:%d",
		peerID.String(),
		timestamp,
	))
}

// SignBinding signs the PeerID + Timestamp
// Signing always happens OFF-CHAIN because the private key never leaves the provider's machine
func SignBinding(
	peerID peer.ID,
	privateKey *ecdsa.PrivateKey,
) (*BindingRequest, error) {

	timestamp := time.Now().Unix()

	message := BuildMessage(peerID, timestamp)

	hash := gethcrypto.Keccak256Hash(message)

	signature, err := gethcrypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	address := gethcrypto.PubkeyToAddress(privateKey.PublicKey)

	return &BindingRequest{
		PeerID:    peerID.String(),
		Wallet:    address.Hex(),
		Timestamp: timestamp,
		Signature: hex.EncodeToString(signature),
	}, nil
}