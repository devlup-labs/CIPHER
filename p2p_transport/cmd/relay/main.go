package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cipher/internal/identity"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	golog "github.com/ipfs/go-log/v2"
)

func main() {
	// Enable libp2p debug logging for circuit v2
	golog.SetLogLevel("relay", "debug")
	golog.SetLogLevel("p2p-circuit", "debug")

	// Load persistent identity for the relay
	priv, err := identity.LoadOrCreate()
	if err != nil {
		log.Fatalf("Failed to load or create identity: %v", err)
	}

	// Listen on TCP 4001 and UDP 4001 (QUIC)
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/4001",
			"/ip4/0.0.0.0/udp/4002/quic-v1",
		),
		libp2p.Identity(priv),
		libp2p.EnableNATService(),
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		log.Fatalf("Failed to create libp2p relay node: %v", err)
	}

	// Configure custom relay resources for development/testing.
	// Production or public relays should stick to relay.DefaultResources() to prevent bandwidth abuse,
	// as relay fallback connections are typically only intended for lightweight protocol coordination.
	rc := relay.DefaultResources()
	rc.Limit.Data = 512 * 1024 * 1024 // 512 MB data limit per connection
	rc.Limit.Duration = 15 * time.Minute // 15 minute duration limit
	rc.MaxReservations = 100

	// Instantiate the circuit v2 relay service
	_, err = relay.New(h, relay.WithResources(rc))
	if err != nil {
		log.Fatalf("Failed to instantiate relay service: %v", err)
	}

	log.Printf("Relay Service Started!")
	log.Printf("Relay Peer ID: %s", h.ID().String())
	
	fmt.Println("\nRelay Multiaddresses (for other peers to connect):")
	for _, addr := range h.Addrs() {
		fmt.Printf("%s/p2p/%s\n", addr.String(), h.ID().String())
	}
	fmt.Println("\nTo deploy this relay publicly, replace the local IP (e.g., 127.0.0.1 or 192.168.x.x) above with your server's PUBLIC IP.")

	// Wait for termination signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Println("Shutting down relay...")
	if err := h.Close(); err != nil {
		log.Fatalf("Failed to close host: %v", err)
	}
}
