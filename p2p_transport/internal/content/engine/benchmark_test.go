package engine_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"os"
	"testing"

	"cipher/internal/content/core"
	"cipher/internal/content/crypto"
	"cipher/internal/content/engine"
	"cipher/internal/content/manifest"
	"cipher/internal/content/storage"
	"cipher/internal/content/verifier"
)

func BenchmarkContentEngine_Ingest(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "content-bench-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config := core.EngineConfig{ChunkSize: 256 * 1024}
	enc := crypto.NewChaCha20Encryptor()
	dig := verifier.NewSHA256Digest()
	keys := engine.NewLocalKeyProvider()

	if err := storage.NewFSStorage(tmpDir); err != nil {
		b.Fatalf("failed to init storage: %v", err)
	}
	store := storage.NewFSStore(tmpDir)

	eng := engine.NewContentEngine(config, enc, dig, store, store, keys)

	// Create 10MB of random data
	dataSize := 10 * 1024 * 1024
	data := make([]byte, dataSize)
	rand.Read(data)

	ctx := context.Background()

	b.SetBytes(int64(dataSize))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		_, err := eng.Ingest(ctx, reader, manifest.TypeFile)
		if err != nil {
			b.Fatalf("Ingest failed: %v", err)
		}
	}
}

func BenchmarkContentEngine_Reassemble(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "content-bench-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config := core.EngineConfig{ChunkSize: 256 * 1024}
	enc := crypto.NewChaCha20Encryptor()
	dig := verifier.NewSHA256Digest()
	keys := engine.NewLocalKeyProvider()
	
	store := storage.NewFSStore(tmpDir)
	eng := engine.NewContentEngine(config, enc, dig, store, store, keys)

	dataSize := 10 * 1024 * 1024
	data := make([]byte, dataSize)
	rand.Read(data)
	ctx := context.Background()

	// Setup: Ingest once
	reader := bytes.NewReader(data)
	m, err := eng.Ingest(ctx, reader, manifest.TypeFile)
	if err != nil {
		b.Fatalf("Ingest failed: %v", err)
	}

	b.SetBytes(int64(dataSize))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var outBuf bytes.Buffer
		if err := eng.Reassemble(ctx, m, &outBuf); err != nil {
			b.Fatalf("Reassemble failed: %v", err)
		}
	}
}
