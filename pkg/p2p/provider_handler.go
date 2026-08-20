package p2p

import (
	"time"

	"github.com/1amKhush/CIPHER/pkg/engine"
	"github.com/1amKhush/CIPHER/pkg/logger"
	"github.com/1amKhush/CIPHER/pkg/wire"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/1amKhush/CIPHER/pkg/ethereum"
)

// ProviderStreamHandler adapts the engine to a libp2p stream.
func ProviderStreamHandler(store *engine.ChunkStore, h host.Host) network.StreamHandler {
	return func(s network.Stream) {
		defer s.Close()
		if err := s.SetDeadline(time.Now().Add(OperationTimeout)); err != nil {
			logger.Error().Err(err).Msg("Failed to set provider stream deadline")
			s.Reset()
			return
		}
		defer s.SetDeadline(time.Time{})

		remotePeer := s.Conn().RemotePeer()
		pubKey := s.Conn().RemotePublicKey()

		// 1. Read ChunkRequest
		reqData, err := wire.ReadMsg(s)
		if err != nil {
			logger.Error().Err(err).Str("peer", remotePeer.String()).Msg("Failed to read ChunkRequest")
			s.Reset()
			return
		}

		var req wire.ChunkRequest
		if err := req.Unmarshal(reqData); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal ChunkRequest")
			s.Reset()
			return
		}

		// 2. Handle Request
		resp, key, err := store.HandleRequest(&req)
		if err != nil {
			logger.Error().Err(err).Msg("Engine failed to handle request")
			s.Reset()
			return
		}

		// 3. Send ChunkResponse
		if err := wire.WriteMsg(s, resp.Marshal()); err != nil {
			logger.Error().Err(err).Msg("Failed to write ChunkResponse")
			s.Reset()
			return
		}

		// 4. Read LotteryTicket
		ticketData, err := wire.ReadMsg(s)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read LotteryTicket")
			s.Reset()
			return
		}

		var ticket wire.LotteryTicket
		if err := ticket.Unmarshal(ticketData); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal LotteryTicket")
			s.Reset()
			return
		}

		// 5. Handle Ticket
		reveal, err := store.HandleTicket(&ticket, key, pubKey)
		if err != nil {
			logger.Error().Err(err).Msg("Engine failed to handle ticket")
			s.Reset()
			return
		}

		// MVP testing only:
		// Simulate a slow provider so the client timeout path can be exercised.
		//
		// Uncomment during local testing to trigger challenge creation.
		// time.Sleep(35 * time.Second)

		// 6. Send KeyReveal
		if err := wire.WriteMsg(s, reveal.Marshal()); err != nil {
			logger.Error().Err(err).Msg("Failed to write KeyReveal")
			s.Reset()
			return
		}

		// MVP:
		// The provider successfully revealed the decryption key.
		//
		// If the client previously opened a challenge because the provider
		// appeared offline, resolve it now.
		//
		// Future:
		// Replace this local function with a smart-contract transaction
		// that marks the challenge as resolved on-chain.
		if err := ethereum.ResolveChallenge(h.ID().String()); err != nil {
			logger.Warn().
				Err(err).
				Msg("No active challenge to resolve")
		}
	}
}
