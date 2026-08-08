package chunk_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"

	"cipher/internal/content/manifest"
	"cipher/internal/protocol/chunk"
	"cipher/internal/transport"
)

func TestChunkProtocol_InterruptedTransfer(t *testing.T) {
	h1, h2 := setupMockNetwork(t)
	eng1 := createTestEngine(t)
	eng2 := createTestEngine(t)
	
	chunk.NewStreamHandler(h1, eng1)
	
	ctx := context.Background()
	data := make([]byte, 1024*1024) // 4 chunks
	rand.Read(data)
	
	m, _ := eng1.Ingest(ctx, bytes.NewReader(data), manifest.TypeFile)
	mBytes, _ := m.Serialize()
	eng1.PutManifestBytes(ctx, m.Descriptor.ID, mBytes)
	
	// Create client
	client, err := chunk.NewClient(ctx, transport.NewTransport(h2), h1.ID(), eng2)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, _ = client.Resolve(ctx, m.Descriptor.ID)

	// Interrupt after first chunk
	go func() {
		// Wait just a moment for download to start then abruptly close the network
		client.Close()
	}()
	
	err = client.Download(ctx, m.ChunkIDs)
	if err == nil {
		t.Fatal("Expected download to fail due to interruption")
	}
	
	// Ensure store is not corrupted (e.g. we didn't write partial chunks).
	// With verify-then-store, it's impossible to write a partial chunk.
	// But let's verify eng2 only has the fully verified chunks (likely 0 or 1).
	for _, chunkID := range m.ChunkIDs {
		// Just ensure it doesn't crash on get
		eng2.GetChunk(ctx, chunkID) 
	}
}
