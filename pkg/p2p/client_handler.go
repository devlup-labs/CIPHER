package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/1amKhush/CIPHER/pkg/engine"
	"github.com/1amKhush/CIPHER/pkg/wire"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// RequestChunk performs the full 4-message handshake for one chunk.
func RequestChunk(ctx context.Context, h host.Host, providerID peer.ID,
	fileID [32]byte, merkleRoot [32]byte, chunkIndex uint64,
	signingKey p2pcrypto.PrivKey) ([]byte, error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, OperationTimeout)
		defer cancel()
		deadline, _ = ctx.Deadline()
	}

	s, err := h.NewStream(ctx, providerID, ProtocolID)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	defer s.Close()
	if err := s.SetDeadline(deadline); err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to set stream deadline: %w", err)
	}
	defer s.SetDeadline(time.Time{})

	// 1. Send ChunkRequest
	var nonce [32]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to read request nonce: %w", err)
	}

	req := wire.ChunkRequest{
		Version:    wire.Version,
		MsgType:    wire.TypeRequest,
		ChunkIndex: chunkIndex,
		Nonce:      nonce,
		FileID:     fileID,
	}

	if err := wire.WriteMsg(s, req.Marshal()); err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to send ChunkRequest: %w", err)
	}

	// 2. Read ChunkResponse
	respData, err := wire.ReadMsg(s)
	if err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to read ChunkResponse: %w", err)
	}

	var resp wire.ChunkResponse
	if err := resp.Unmarshal(respData); err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to unmarshal ChunkResponse: %w", err)
	}

	// 3. Verify Response bounds
	if err := engine.VerifyResponse(&resp); err != nil {
		s.Reset()
		return nil, fmt.Errorf("response verification failed: %w", err)
	}

	// 4. Send LotteryTicket
	ticket := wire.LotteryTicket{
		Version:     wire.Version,
		MsgType:     wire.TypeTicket,
		TargetBlock: 0,   // Stub for MVP
		WinProb:     100, // Stub for MVP
	}

	copy(ticket.ChannelID[:], nonce[:]) // Use nonce as mock channel ID for MVP
	copy(ticket.HResp[:], resp.HResp[:])

	sig, err := signingKey.Sign(ticket.DataToSign())
	if err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to sign ticket: %w", err)
	}
	ticket.Signature = sig

	if err := wire.WriteMsg(s, ticket.Marshal()); err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to send LotteryTicket: %w", err)
	}

	// 5. Read KeyReveal
	revealData, err := wire.ReadMsg(s)
	if err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to read KeyReveal: %w", err)
	}

	var reveal wire.KeyReveal
	if err := reveal.Unmarshal(revealData); err != nil {
		s.Reset()
		return nil, fmt.Errorf("failed to unmarshal KeyReveal: %w", err)
	}

	// 6. Verify Reveal, Decrypt, and Verify Merkle Proof
	plaintext, err := engine.VerifyReveal(&reveal, &resp, merkleRoot, fileID, chunkIndex)
	if err != nil {
		s.Reset()
		return nil, fmt.Errorf("reveal verification failed: %w", err)
	}

	return plaintext, nil
}
