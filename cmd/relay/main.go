package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"

	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"
)

func loadIdentity() (crypto.PrivKey, error) {
	encoded := os.Getenv("RELAY_PRIVATE_KEY")

	if encoded == "" {
		return nil, fmt.Errorf("RELAY_PRIVATE_KEY environment variable is not set")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode relay private key: %w", err)
	}

	privateKey, err := crypto.UnmarshalPrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal relay private key: %w", err)
	}

	return privateKey, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	privateKey, err := loadIdentity()
	if err != nil {
		panic(err)
	}

	h, err := libp2p.New(
		libp2p.Identity(privateKey),

		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", port),
		),

		libp2p.Transport(ws.New),
	)

	if err != nil {
		panic(err)
	}

	relayService, err := relay.New(h)
	if err != nil {
		panic(err)
	}
	defer relayService.Close()

	fmt.Println("========================================")
	fmt.Println("CIPHER RELAY STARTED")
	fmt.Println("Relay Peer ID:", h.ID())
	fmt.Println("Port:", port)
	fmt.Println("Circuit Relay v2 service: ACTIVE")
	fmt.Println("Persistent Identity: ACTIVE")
	fmt.Println("========================================")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	h.Close()
}