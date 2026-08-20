 package main

import (
	"context"
	"encoding/hex"
	"flag"
	"os"
	//"time"

	"github.com/1amKhush/CIPHER/pkg/logger"
	"github.com/1amKhush/CIPHER/pkg/p2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	
	 
)

func main() {
	//used to parse command-line argument
	providerAddr := flag.String("provider", "", "Provider multiaddr")
	rootHex := flag.String("root", "", "Merkle root in hex")
	chunksCount := flag.Uint64("chunks", 4, "Number of chunks to download")
	relayAddr := flag.String("relay", "", "Relay multiaddr to connect to (optional)")
	verbose := flag.Bool("verbose", false, "Enable verbose debug logging")
	enableQUIC := flag.Bool("quic", false, "Enable QUIC transport")
	enableMDNS := flag.Bool("mdns", true, "Enable mDNS discovery")
	flag.Parse()
	
	cfg := logger.DefaultConfig()
	if *verbose {
		cfg.Level = "debug"
	}
	logger.Init(cfg)

	if *providerAddr == "" || *rootHex == "" {
		logger.Fatal().Msg("-provider and -root flags are required")
	}

	rootBytes, err := hex.DecodeString(*rootHex)
	if err != nil || len(rootBytes) != 32 {
		logger.Fatal().Err(err).Msg("Invalid merkle root")
	}
	var merkleRoot [32]byte
	copy(merkleRoot[:], rootBytes)

	maddr, err := multiaddr.NewMultiaddr(*providerAddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Invalid multiaddr")
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to extract peer info")
	}

	opts := p2p.HostOptions{
		ListenPort:  0,
		PrivKeyPath: "client_key.key",
		EnableMDNS:  *enableMDNS,
		RelayAddr:   *relayAddr,
		EnableQUIC:  *enableQUIC,
	}
	startupCtx, cancelStartup := context.WithTimeout(context.Background(), p2p.OperationTimeout)
	defer cancelStartup()
	h, err := p2p.NewHost(startupCtx, opts)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start host")
	}
	defer h.Close()

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), p2p.OperationTimeout)
	defer cancelConnect()
	if err := h.Connect(connectCtx, *info); err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to provider")
	}
	logger.Info().Msgf("Connected to provider %s", info.ID)

	privKey := p2p.GetHostPrivateKey(h)

	var fileID [32]byte // zeroed

	outFile, err := os.Create("downloaded_file.txt")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create output file")
	}
	defer outFile.Close()

	for i := uint64(0); i < *chunksCount; i++ {
		requestCtx, cancelRequest := context.WithTimeout(context.Background(), p2p.OperationTimeout)
		plaintext, err := p2p.RequestChunk(requestCtx, h, info.ID, fileID, merkleRoot, i, privKey)
		cancelRequest()
		if err != nil {
		logger.Error().
			Err(err).
			Msgf("Failed to request chunk %d", i)

		//time.Sleep(40 * time.Second)

		return
		}

		if _, err := outFile.Write(plaintext); err != nil {
			logger.Fatal().Err(err).Msgf("Failed to write chunk %d to file", i)
		}

		logger.Info().Msgf("Successfully downloaded chunk %d (%d bytes)", i, len(plaintext))
	}

	logger.Info().Msg("File download complete!")
}
