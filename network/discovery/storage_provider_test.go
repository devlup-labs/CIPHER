package discovery

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func TestStorageProviderDHTRegistration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 1. Create Bootstrapper DHT Node
	h1, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"))
	if err != nil {
		t.Fatalf("failed to create h1: %v", err)
	}
	defer h1.Close()

	dht1, err := NewDHT(h1, dht.ModeServer)
	if err != nil {
		t.Fatalf("failed to create dht1: %v", err)
	}
	defer dht1.Close()

	// 2. Create Storage Provider DHT Node
	h2, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"))
	if err != nil {
		t.Fatalf("failed to create h2: %v", err)
	}
	defer h2.Close()

	dht2, err := NewDHT(h2, dht.ModeServer)
	if err != nil {
		t.Fatalf("failed to create dht2: %v", err)
	}
	defer dht2.Close()

	h1Info := peer.AddrInfo{ID: h1.ID(), Addrs: h1.Addrs()}

	// Bootstrap h2 with h1
	if err := Bootstrap(ctx, dht2, h2, []peer.AddrInfo{h1Info}); err != nil {
		t.Fatalf("dht2 bootstrap failed: %v", err)
	}

	// 3. Register h2 as Storage Provider
	if err := RegisterStorageProvider(ctx, dht2); err != nil {
		t.Fatalf("RegisterStorageProvider failed: %v", err)
	}

	// 4. Create Publisher DHT Node
	h3, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"))
	if err != nil {
		t.Fatalf("failed to create h3: %v", err)
	}
	defer h3.Close()

	dht3, err := NewDHT(h3, dht.ModeServer)
	if err != nil {
		t.Fatalf("failed to create dht3: %v", err)
	}
	defer dht3.Close()

	// Bootstrap h3 with h1
	if err := Bootstrap(ctx, dht3, h3, []peer.AddrInfo{h1Info}); err != nil {
		t.Fatalf("dht3 bootstrap failed: %v", err)
	}

	// Give DHT record a moment to settle
	time.Sleep(300 * time.Millisecond)

	// 5. Query DHT for Storage Providers from h3
	discovered, err := FindStorageProviders(ctx, dht3, 5)
	if err != nil {
		t.Fatalf("FindStorageProviders failed: %v", err)
	}

	found := false
	for _, p := range discovered {
		if p.ID == h2.ID() {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected to discover storage provider %s, got %v", h2.ID(), discovered)
	}
}
