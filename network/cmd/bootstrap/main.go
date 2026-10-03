package main

import (
	"cipher/network/identity"
	"cipher/network/transport"
	"cipher/shared/logger"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func main() {
	port := flag.Int("p", 4003, "Port for the bootstrap node (TCP)")
	wsPort := flag.Int("ws-port", 0, "Port for the bootstrap node (WebSocket, 0 to disable)")
	identityPath := flag.String("identity", "", "Custom path to identity key file (optional)")
	flag.Parse()

	log := logger.Bootstrap

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var priv libp2pcrypto.PrivKey
	var err error
	if *identityPath != "" {
		priv, err = identity.LoadOrCreateFromPath(*identityPath)
	} else {
		priv, err = identity.LoadOrCreate()
	}
	if err != nil {
		log.Fatalf("Failed to load or create identity: %v", err)
	}

	h, kdht, err := transport.NewNode(ctx, *port, *wsPort, priv, "", false)
	if err != nil {
		log.Fatalf("Failed to create bootstrap node: %v", err)
	}
	defer h.Close()
	defer kdht.Close()

	var addrStrings []string
	for _, addr := range h.Addrs() {
		addrStrings = append(addrStrings, fmt.Sprintf("%s/p2p/%s", addr, h.ID()))
	}

	fields := []logger.Field{
		{Key: "Node Role", Value: "Kademlia DHT Bootstrap Coordinator"},
		{Key: "Peer ID", Value: h.ID().String()},
		{Key: "TCP Port", Value: fmt.Sprintf("%d", *port)},
	}
	if *wsPort > 0 {
		fields = append(fields, logger.Field{Key: "WebSocket Port", Value: fmt.Sprintf("%d", *wsPort)})
	}
	fields = append(fields, logger.Field{Key: "", Value: "Multiaddresses (copy for other nodes):"})
	for _, a := range addrStrings {
		fields = append(fields, logger.Field{Key: "", Value: "  " + a})
	}

	log.Banner("CIPHER DHT BOOTSTRAP NODE", fields...)
	log.Success("Bootstrap node is ready and actively routing Kademlia DHT queries")
	log.Info("Waiting for provider, publisher, and consumer peers to connect...")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Warn("Shutting down bootstrap node...")
}

