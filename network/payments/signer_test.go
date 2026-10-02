package payments

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSigner_TicketSignAndRecover(t *testing.T) {
	// Generate random ECDSA private key
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	expectedAddr := crypto.PubkeyToAddress(privKey.PublicKey)

	chainID := big.NewInt(31337) // Anvil default chain ID
	verifyingContract := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

	signer := NewTicketSigner(privKey, chainID, verifyingContract)

	ticket := RoundTicket{
		Sender:      expectedAddr,
		Recipient:   common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
		RoundID:     big.NewInt(1),
		LocalIndex:  big.NewInt(8),
		FaceValue:   big.NewInt(1000000000000000000), // 1 ETH
		WinProb:     big.NewInt(100000000000000000),  // 0.1
		SenderNonce: big.NewInt(42),
	}

	sig, err := signer.SignTicket(ticket)
	if err != nil {
		t.Fatalf("Failed to sign ticket: %v", err)
	}

	if len(sig) != 65 {
		t.Fatalf("Invalid signature length: got %d, expected 65", len(sig))
	}

	if sig[64] != 27 && sig[64] != 28 {
		t.Fatalf("Invalid V value: got %d, expected 27 or 28", sig[64])
	}

	recovered, err := signer.RecoverSigner(ticket, sig)
	if err != nil {
		t.Fatalf("Failed to recover signer: %v", err)
	}

	if recovered != expectedAddr {
		t.Fatalf("Recovered address mismatch: got %s, expected %s", recovered.Hex(), expectedAddr.Hex())
	}

	if !signer.Verify(ticket, sig, expectedAddr) {
		t.Fatalf("Verify returned false for valid signature")
	}

	// Verify tampering detection
	tamperedTicket := ticket
	tamperedTicket.LocalIndex = big.NewInt(9)
	if signer.Verify(tamperedTicket, sig, expectedAddr) {
		t.Fatalf("Verify should return false for tampered ticket")
	}
}
