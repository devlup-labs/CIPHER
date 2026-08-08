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

func BenchmarkChunkTransport_Sequential(b *testing.B) {
	h1, h2 := setupMockNetwork(b)
	eng1 := createTestEngine(b)
	eng2 := createTestEngine(b)
	
	chunk.NewStreamHandler(h1, eng1)
	
	ctx := context.Background()
	
	// Create 10MB of data to transfer per benchmark iteration
	dataSize := 10 * 1024 * 1024
	data := make([]byte, dataSize)
	rand.Read(data)
	
	m, _ := eng1.Ingest(ctx, bytes.NewReader(data), manifest.TypeFile)
	mBytes, _ := m.Serialize()
	eng1.PutManifestBytes(ctx, m.Descriptor.ID, mBytes)
	
	b.SetBytes(int64(dataSize))
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		client, err := chunk.NewClient(ctx, transport.NewTransport(h2), h1.ID(), eng2)
		if err != nil {
			b.Fatalf("Failed to create client: %v", err)
		}
		
		_, err = client.Resolve(ctx, m.Descriptor.ID)
		if err != nil {
			b.Fatalf("Resolve failed: %v", err)
		}
		
		err = client.Download(ctx, m.ChunkIDs)
		if err != nil {
			b.Fatalf("Download failed: %v", err)
		}
		
		client.Close()
	}
}
