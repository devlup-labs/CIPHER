package main

import (
	"cipher/network/identity"
	"cipher/shared/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

func main() {
	golog.SetLogLevel("relay", "warn")
	golog.SetLogLevel("p2p-circuit", "warn")

	log := logger.Relay

	// Load persistent identity for the relay
	priv, err := identity.LoadOrCreate()
	if err != nil {
		log.Fatalf("Failed to load or create identity: %v", err)
	}

	// Listen on TCP 4001, UDP 4002 (QUIC), and TCP 4004 (WebSocket)
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/4001",
			"/ip4/0.0.0.0/udp/4002/quic-v1",
			"/ip4/0.0.0.0/tcp/4004/ws",
		),
		libp2p.Identity(priv),
		libp2p.EnableNATService(),
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		log.Fatalf("Failed to create libp2p relay node: %v", err)
	}

	rc := relay.DefaultResources()
	rc.Limit.Data = 512 * 1024 * 1024 // 512 MB data limit per connection
	rc.Limit.Duration = 15 * time.Minute // 15 minute duration limit
	rc.MaxReservations = 100

	_, err = relay.New(h, relay.WithResources(rc))
	if err != nil {
		log.Fatalf("Failed to instantiate relay service: %v", err)
	}

	fields := []logger.Field{
		{Key: "Node Role", Value: "Circuit Relay v2 Service"},
		{Key: "Relay Peer ID", Value: h.ID().String()},
		{Key: "Max Reservations", Value: "100"},
		{Key: "Data Limit", Value: "512 MB / connection"},
		{Key: "", Value: "Relay Multiaddresses (for other peers):"},
	}
	for _, addr := range h.Addrs() {
		fields = append(fields, logger.Field{Key: "", Value: fmt.Sprintf("  %s/p2p/%s", addr.String(), h.ID().String())})
	}

	log.Banner("CIPHER CIRCUIT RELAY V2", fields...)
	log.Success("Relay node is operational and accepting peer reservations")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Warn("Shutting down relay...")
	if err := h.Close(); err != nil {
		log.Fatalf("Failed to close host: %v", err)
	}
}

