package chunk_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"cipher/network/content/core"
	"cipher/network/content/manifest"
	"cipher/network/payments"
	"cipher/network/protocol/chunk"
	"cipher/network/transport"
)

func TestChunkProtocol_TicketPaymentIntegration(t *testing.T) {
	h1, h2 := setupMockNetwork(t)

	eng1 := createTestEngine(t)
	eng2 := createTestEngine(t)

	// Ethereum identities
	clientEthKey, err := ethcrypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	clientEthAddr := ethcrypto.PubkeyToAddress(clientEthKey.PublicKey)

	providerEthKey, err := ethcrypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	providerEthAddr := ethcrypto.PubkeyToAddress(providerEthKey.PublicKey)

	chainID := big.NewInt(31337)
	verifyingContract := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

	clientSigner := payments.NewTicketSigner(clientEthKey, chainID, verifyingContract)
	providerSigner := payments.NewTicketSigner(nil, chainID, verifyingContract)

	// Setup Provider handler with ticket verification
	handler1 := chunk.NewStreamHandler(h1, eng1)
	chunk.NewStreamHandler(h2, eng2)

	var receivedTickets []*payments.SignedTicket
	var mu sync.Mutex

	handler1.SetTicketHandler(func(signedTicket *payments.SignedTicket) error {
		mu.Lock()
		defer mu.Unlock()

		// Verify ticket signature
		if !providerSigner.Verify(signedTicket.Ticket, signedTicket.Signature, clientEthAddr) {
			return errors.New("signature verification failed")
		}
		if signedTicket.Ticket.Recipient != providerEthAddr {
			return errors.New("recipient address mismatch")
		}

		receivedTickets = append(receivedTickets, signedTicket)
		return nil
	})

	// Put test content on eng1 (Provider)
	ctx := context.Background()
	payload := make([]byte, 512*1024) // 512 KiB = 2 chunks of 256 KiB
	rand.Read(payload)

	m, err := eng1.Ingest(ctx, bytes.NewReader(payload), manifest.TypeFile)
	if err != nil {
		t.Fatal(err)
	}
	mBytes, _ := m.Serialize()
	eng1.PutManifestBytes(ctx, m.Descriptor.ID, mBytes)
	t.Logf("Total chunks in manifest: %d", len(m.ChunkIDs))

	// Client dials Provider
	t2 := transport.NewTransport(h2)
	client, err := chunk.NewClient(ctx, t2, h1.ID(), eng2)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Configure client ticket generator
	chunkIdx := uint64(0)
	client.SetTicketGenerator(func(chunkID core.ChunkID) (*payments.SignedTicket, error) {
		ticket := payments.RoundTicket{
			Sender:      clientEthAddr,
			Recipient:   providerEthAddr,
			RoundID:     big.NewInt(1),
			LocalIndex:  new(big.Int).SetUint64(chunkIdx),
			FaceValue:   big.NewInt(1000000000000000000), // 1 ETH
			WinProb:     big.NewInt(100000000000000000),
			SenderNonce: big.NewInt(12345),
		}
		sig, err := clientSigner.SignTicket(ticket)
		if err != nil {
			return nil, err
		}
		chunkIdx++
		return &payments.SignedTicket{
			Ticket:    ticket,
			Signature: sig,
		}, nil
	})

	// Fetch all chunks
	for _, chunkID := range m.ChunkIDs {
		c, err := client.FetchChunk(ctx, chunkID)
		if err != nil {
			t.Fatalf("FetchChunk failed: %v", err)
		}
		if c == nil {
			t.Fatal("Expected chunk data, got nil")
		}
	}
	client.Close()

	// Wait up to 500ms for handler goroutine to finish reading final ticket
	for i := 0; i < 50; i++ {
		mu.Lock()
		count := len(receivedTickets)
		mu.Unlock()
		if count == len(m.ChunkIDs) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(receivedTickets) != len(m.ChunkIDs) {
		t.Fatalf("Expected %d received tickets, got %d", len(m.ChunkIDs), len(receivedTickets))
	}

	for i, st := range receivedTickets {
		if st.Ticket.LocalIndex.Uint64() != uint64(i) {
			t.Errorf("Ticket %d has localIndex %d", i, st.Ticket.LocalIndex.Uint64())
		}
		if st.Ticket.Sender != clientEthAddr {
			t.Errorf("Ticket %d has invalid sender: %s", i, st.Ticket.Sender.Hex())
		}
	}
}

func TestChunkProtocol_TicketPayment_RejectedInvalidSignature(t *testing.T) {
	h1, h2 := setupMockNetwork(t)

	eng1 := createTestEngine(t)
	eng2 := createTestEngine(t)

	// Legitimate client address configured on provider
	legitClientEthKey, _ := ethcrypto.GenerateKey()
	legitClientEthAddr := ethcrypto.PubkeyToAddress(legitClientEthKey.PublicKey)

	// Attacker client key
	attackerEthKey, _ := ethcrypto.GenerateKey()

	providerEthKey, _ := ethcrypto.GenerateKey()
	providerEthAddr := ethcrypto.PubkeyToAddress(providerEthKey.PublicKey)

	chainID := big.NewInt(31337)
	verifyingContract := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

	attackerSigner := payments.NewTicketSigner(attackerEthKey, chainID, verifyingContract)
	providerSigner := payments.NewTicketSigner(nil, chainID, verifyingContract)

	handler1 := chunk.NewStreamHandler(h1, eng1)
	chunk.NewStreamHandler(h2, eng2)

	rejectedCount := 0
	var mu sync.Mutex

	handler1.SetTicketHandler(func(signedTicket *payments.SignedTicket) error {
		mu.Lock()
		defer mu.Unlock()

		// Provider verifies against legitimate client address
		if !providerSigner.Verify(signedTicket.Ticket, signedTicket.Signature, legitClientEthAddr) {
			rejectedCount++
			return errors.New("unauthorized ticket signer")
		}
		return nil
	})

	// Put test content on eng1 (Provider)
	ctx := context.Background()
	payload := make([]byte, 512*1024)
	rand.Read(payload)

	m, err := eng1.Ingest(ctx, bytes.NewReader(payload), manifest.TypeFile)
	if err != nil {
		t.Fatal(err)
	}
	mBytes, _ := m.Serialize()
	eng1.PutManifestBytes(ctx, m.Descriptor.ID, mBytes)

	// Attacker dials Provider
	t2 := transport.NewTransport(h2)
	client, err := chunk.NewClient(ctx, t2, h1.ID(), eng2)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Attacker signs ticket claiming to be legitClient, but signed by attacker key
	client.SetTicketGenerator(func(chunkID core.ChunkID) (*payments.SignedTicket, error) {
		ticket := payments.RoundTicket{
			Sender:      legitClientEthAddr,
			Recipient:   providerEthAddr,
			RoundID:     big.NewInt(1),
			LocalIndex:  big.NewInt(0),
			FaceValue:   big.NewInt(1000000000000000000),
			WinProb:     big.NewInt(100000000000000000),
			SenderNonce: big.NewInt(999),
		}
		// Signed by attacker private key!
		sig, err := attackerSigner.SignTicket(ticket)
		if err != nil {
			return nil, err
		}
		return &payments.SignedTicket{
			Ticket:    ticket,
			Signature: sig,
		}, nil
	})

	// First chunk will be served, but ticket will be rejected by provider, terminating the stream
	_, err = client.FetchChunk(ctx, m.ChunkIDs[0])
	if err != nil {
		t.Fatalf("First chunk fetch failed before ticket rejection: %v", err)
	}

	// Give handler a moment to process the rejected ticket and terminate stream
	time.Sleep(30 * time.Millisecond)

	// Subsequent chunk fetch must fail because provider closed the stream on ticket rejection
	_, err = client.FetchChunk(ctx, m.ChunkIDs[1])
	if err == nil {
		t.Fatal("Expected error on subsequent chunk fetch after forged ticket, got nil")
	}

	mu.Lock()
	defer mu.Unlock()
	if rejectedCount != 1 {
		t.Fatalf("Expected 1 rejected ticket, got %d", rejectedCount)
	}
}

func TestChunkProtocol_TicketPayment_TamperedChunkIndex(t *testing.T) {
	chainID := big.NewInt(31337)
	verifyingContract := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

	privKey, _ := ethcrypto.GenerateKey()
	clientAddr := ethcrypto.PubkeyToAddress(privKey.PublicKey)
	signer := payments.NewTicketSigner(privKey, chainID, verifyingContract)

	ticket := payments.RoundTicket{
		Sender:      clientAddr,
		Recipient:   common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
		RoundID:     big.NewInt(1),
		LocalIndex:  big.NewInt(5),
		FaceValue:   big.NewInt(1000000000000000000),
		WinProb:     big.NewInt(100000000000000000),
		SenderNonce: big.NewInt(1),
	}

	sig, err := signer.SignTicket(ticket)
	if err != nil {
		t.Fatal(err)
	}

	// Tamper with local index
	tamperedTicket := ticket
	tamperedTicket.LocalIndex = big.NewInt(6)

	if signer.Verify(tamperedTicket, sig, clientAddr) {
		t.Fatal("Expected signature verification to FAIL for tampered chunk index, but it passed")
	}

	// Tamper with face value
	tamperedFaceValue := ticket
	tamperedFaceValue.FaceValue = big.NewInt(2000000000000000000)

	if signer.Verify(tamperedFaceValue, sig, clientAddr) {
		t.Fatal("Expected signature verification to FAIL for tampered face value, but it passed")
	}
}

