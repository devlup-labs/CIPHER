package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/1amKhush/CIPHER/pkg/chunker"
	"github.com/1amKhush/CIPHER/pkg/crypto"
	"github.com/1amKhush/CIPHER/pkg/engine"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func TestP2PLoopback(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
	defer cancel()

	// 1. Setup provider
	providerOpts := HostOptions{ListenPort: 0, PrivKeyPath: "", EnableMDNS: false}
	providerHost, err := NewHost(ctx, providerOpts)
	if err != nil {
		t.Fatalf("failed to create provider host: %v", err)
	}
	defer providerHost.Close()

	// Prepare data
	dummyData := make([]byte, 100*1024)
	rand.Read(dummyData)
	chunks := chunker.ChunkBytes(dummyData)
	var leaves [][32]byte
	for _, c := range chunks {
		var fileID [32]byte
		leaf := crypto.MerkleLeaf(fileID, c.Index, uint32(len(c.Data)), c.Data)
		leaves = append(leaves, leaf)
	}
	tree := chunker.NewMerkleTree(leaves)

	var fileID [32]byte
	store := &engine.ChunkStore{FileID: fileID, Chunks: chunks, MerkleTree: tree}
	providerHost.SetStreamHandler(ProtocolID, ProviderStreamHandler(store,providerHost))

	// 2. Setup client
	clientOpts := HostOptions{ListenPort: 0, PrivKeyPath: "", EnableMDNS: false}
	clientHost, err := NewHost(ctx, clientOpts)
	if err != nil {
		t.Fatalf("failed to create client host: %v", err)
	}
	defer clientHost.Close()

	providerInfo := peer.AddrInfo{
		ID:    providerHost.ID(),
		Addrs: providerHost.Addrs(),
	}

	if err := clientHost.Connect(ctx, providerInfo); err != nil {
		t.Fatalf("client failed to connect: %v", err)
	}

	// 3. Test loopback for all chunks
	privKey := GetHostPrivateKey(clientHost)
	var downloadedData []byte

	for i := uint64(0); i < uint64(len(chunks)); i++ {
		plaintext, err := RequestChunk(ctx, clientHost, providerHost.ID(), fileID, tree.Root, i, privKey)
		if err != nil {
			t.Fatalf("RequestChunk %d failed: %v", i, err)
		}

		if len(plaintext) != len(chunks[i].Data) {
			t.Fatalf("Length mismatch for chunk %d: got %d, expected %d", i, len(plaintext), len(chunks[i].Data))
		}
		downloadedData = append(downloadedData, plaintext...)
	}

	if len(downloadedData) != len(dummyData) {
		t.Fatalf("Final length mismatch: got %d, expected %d", len(downloadedData), len(dummyData))
	}
}

func TestP2PRelay(t *testing.T) {
	t.Skip("requires external public relay; run manually as an integration test")
	ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
	defer cancel()

	relayAddr := "/dns4/relay-torrentium-3zok.onrender.com/tcp/443/wss/p2p/12D3KooWEBxhvkASAJtmdeKWiWWhdXCzwXEVvSMpjuY8YrDAi68Z"

	// 1. Setup provider
	providerOpts := HostOptions{ListenPort: 0, PrivKeyPath: "", EnableMDNS: false, RelayAddr: relayAddr}
	providerHost, err := NewHost(ctx, providerOpts)
	if err != nil {
		t.Fatalf("failed to create provider host: %v", err)
	}
	defer providerHost.Close()

	// Prepare data
	dummyData := make([]byte, 100*1024)
	rand.Read(dummyData)
	chunks := chunker.ChunkBytes(dummyData)
	var leaves [][32]byte
	for _, c := range chunks {
		var fileID [32]byte
		leaf := crypto.MerkleLeaf(fileID, c.Index, uint32(len(c.Data)), c.Data)
		leaves = append(leaves, leaf)
	}
	tree := chunker.NewMerkleTree(leaves)

	var fileID [32]byte
	store := &engine.ChunkStore{FileID: fileID, Chunks: chunks, MerkleTree: tree}
	providerHost.SetStreamHandler(ProtocolID, ProviderStreamHandler(store,providerHost))

	// Wait a bit for relay reservation to settle
	importTime := time.Millisecond * 500
	<-time.After(importTime)

	// 2. Setup client
	clientOpts := HostOptions{ListenPort: 0, PrivKeyPath: "", EnableMDNS: false, RelayAddr: relayAddr}
	clientHost, err := NewHost(ctx, clientOpts)
	if err != nil {
		t.Fatalf("failed to create client host: %v", err)
	}
	defer clientHost.Close()

	// Connect to provider via relay circuit
	circuitAddrStr := fmt.Sprintf("%s/p2p-circuit/p2p/%s", relayAddr, providerHost.ID().String())
	maddr, _ := multiaddr.NewMultiaddr(circuitAddrStr)
	providerInfo, _ := peer.AddrInfoFromP2pAddr(maddr)

	if err := clientHost.Connect(ctx, *providerInfo); err != nil {
		t.Fatalf("client failed to connect to provider via relay: %v", err)
	}

	// 3. Test chunk 0 fetch via relay
	privKey := GetHostPrivateKey(clientHost)
	plaintext, err := RequestChunk(ctx, clientHost, providerHost.ID(), fileID, tree.Root, 0, privKey)
	if err != nil {
		t.Fatalf("RequestChunk via relay failed: %v", err)
	}

	if len(plaintext) != len(chunks[0].Data) {
		t.Fatalf("Length mismatch: got %d, expected %d", len(plaintext), len(chunks[0].Data))
	}
}
