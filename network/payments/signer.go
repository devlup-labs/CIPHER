package payments

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	EIP712DomainTypeHash = crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	RoundTicketTypeHash  = crypto.Keccak256Hash([]byte("RoundTicket(address sender,address recipient,uint256 roundId,uint256 localIndex,uint256 faceValue,uint256 winProb,uint256 senderNonce)"))

	ProtocolNameHash    = crypto.Keccak256Hash([]byte("CIPHER Payment Protocol"))
	ProtocolVersionHash = crypto.Keccak256Hash([]byte("1.0.0"))
)

// RoundTicket matches the Solidity IEntropySource.RoundTicket struct.
type RoundTicket struct {
	Sender      common.Address `json:"sender"`
	Recipient   common.Address `json:"recipient"`
	RoundID     *big.Int       `json:"roundId"`
	LocalIndex  *big.Int       `json:"localIndex"`
	FaceValue   *big.Int       `json:"faceValue"`
	WinProb     *big.Int       `json:"winProb"`
	SenderNonce *big.Int       `json:"senderNonce"`
}

// SignedTicket bundles a RoundTicket with its 65-byte cryptographic signature.
type SignedTicket struct {
	Ticket    RoundTicket `json:"ticket"`
	Signature []byte      `json:"signature"`
}

// ComputeDomainSeparator calculates the EIP-712 domain separator matching OpenZeppelin EIP712.
func ComputeDomainSeparator(chainID *big.Int, verifyingContract common.Address) common.Hash {
	var encoded []byte
	encoded = append(encoded, EIP712DomainTypeHash.Bytes()...)
	encoded = append(encoded, ProtocolNameHash.Bytes()...)
	encoded = append(encoded, ProtocolVersionHash.Bytes()...)
	encoded = append(encoded, math.PaddedBigBytes(chainID, 32)...)
	encoded = append(encoded, common.LeftPadBytes(verifyingContract.Bytes(), 32)...)

	return crypto.Keccak256Hash(encoded)
}

// HashRoundTicket calculates the EIP-712 struct hash for a RoundTicket.
func HashRoundTicket(ticket RoundTicket) common.Hash {
	var encoded []byte
	encoded = append(encoded, RoundTicketTypeHash.Bytes()...)
	encoded = append(encoded, common.LeftPadBytes(ticket.Sender.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(ticket.Recipient.Bytes(), 32)...)
	encoded = append(encoded, math.PaddedBigBytes(ticket.RoundID, 32)...)
	encoded = append(encoded, math.PaddedBigBytes(ticket.LocalIndex, 32)...)
	encoded = append(encoded, math.PaddedBigBytes(ticket.FaceValue, 32)...)
	encoded = append(encoded, math.PaddedBigBytes(ticket.WinProb, 32)...)
	encoded = append(encoded, math.PaddedBigBytes(ticket.SenderNonce, 32)...)

	return crypto.Keccak256Hash(encoded)
}

// HashTypedData creates the EIP-712 \x19\x01 digest from domain separator and struct hash.
func HashTypedData(domainSeparator common.Hash, structHash common.Hash) common.Hash {
	prefix := []byte("\x19\x01")
	data := append(prefix, domainSeparator.Bytes()...)
	data = append(data, structHash.Bytes()...)
	return crypto.Keccak256Hash(data)
}

// TicketSigner manages EIP-712 signing and verification for payment tickets.
type TicketSigner struct {
	ChainID           *big.Int
	VerifyingContract common.Address
	DomainSeparator   common.Hash
	PrivateKey        *ecdsa.PrivateKey
	Address           common.Address
}

// NewTicketSigner creates a new TicketSigner instance.
func NewTicketSigner(privKey *ecdsa.PrivateKey, chainID *big.Int, verifyingContract common.Address) *TicketSigner {
	domainSep := ComputeDomainSeparator(chainID, verifyingContract)
	var addr common.Address
	if privKey != nil {
		addr = crypto.PubkeyToAddress(privKey.PublicKey)
	}
	return &TicketSigner{
		ChainID:           chainID,
		VerifyingContract: verifyingContract,
		DomainSeparator:   domainSep,
		PrivateKey:        privKey,
		Address:           addr,
	}
}

// Digest computes the EIP-712 digest for a given ticket.
func (s *TicketSigner) Digest(ticket RoundTicket) common.Hash {
	structHash := HashRoundTicket(ticket)
	return HashTypedData(s.DomainSeparator, structHash)
}

// SignTicket signs a RoundTicket using the configured private key.
// Returns 65-byte signature in [R || S || V] format with V in {27, 28} for EVM compatibility.
func (s *TicketSigner) SignTicket(ticket RoundTicket) ([]byte, error) {
	if s.PrivateKey == nil {
		return nil, errors.New("cannot sign ticket: private key is nil")
	}

	digest := s.Digest(ticket)
	sig, err := crypto.Sign(digest.Bytes(), s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign ticket digest: %w", err)
	}

	// crypto.Sign returns V in {0, 1}; EVM ecrecover requires V in {27, 28}
	if sig[64] < 27 {
		sig[64] += 27
	}

	return sig, nil
}

// RecoverSigner recovers the Ethereum address that signed the ticket.
func (s *TicketSigner) RecoverSigner(ticket RoundTicket, sig []byte) (common.Address, error) {
	if len(sig) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length: %d (expected 65)", len(sig))
	}

	// Copy signature to avoid mutating caller's slice
	sigCopy := make([]byte, 65)
	copy(sigCopy, sig)

	// Transform V from {27, 28} back to {0, 1} for crypto.SigToPub
	if sigCopy[64] >= 27 {
		sigCopy[64] -= 27
	}

	digest := s.Digest(ticket)
	pubKey, err := crypto.SigToPub(digest.Bytes(), sigCopy)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover public key: %w", err)
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

// Verify verifies that the ticket was signed by expectedSigner.
func (s *TicketSigner) Verify(ticket RoundTicket, sig []byte, expectedSigner common.Address) bool {
	recovered, err := s.RecoverSigner(ticket, sig)
	if err != nil {
		return false
	}
	return recovered == expectedSigner
}
