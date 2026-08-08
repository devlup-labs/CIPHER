package engine

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"sort"

	"cipher/internal/content/chunker"
	"cipher/internal/content/core"
	"cipher/internal/content/manifest"
)

type ContentEngine struct {
	config    core.EngineConfig
	chunker   *chunker.Chunker
	encryptor core.Encryptor
	digest    core.Digest
	source    core.ChunkSource
	sink      core.ChunkSink
	keys      core.KeyProvider

	// Simple in-memory manifest store for Milestone 8 transport integration
	manifests map[core.ContentID][]byte
}

func NewContentEngine(
	config core.EngineConfig,
	encryptor core.Encryptor,
	digest core.Digest,
	source core.ChunkSource,
	sink core.ChunkSink,
	keys core.KeyProvider,
) *ContentEngine {
	return &ContentEngine{
		config:    config,
		chunker:   chunker.NewChunker(config),
		encryptor: encryptor,
		digest:    digest,
		source:    source,
		sink:      sink,
		keys:      keys,
		manifests: make(map[core.ContentID][]byte),
	}
}

// Ingest reads a file, chunks it, encrypts it, stores it, and returns the manifest.
func (e *ContentEngine) Ingest(ctx context.Context, r io.Reader, mtype manifest.ContentType) (*manifest.Manifest, error) {
	chunkCh, errCh := e.chunker.Split(r)

	// Generate a unique ContentID for this upload
	var contentID core.ContentID
	if _, err := rand.Read(contentID[:]); err != nil {
		return nil, fmt.Errorf("failed to generate content id: %w", err)
	}

	// Generate a new encryption key
	key := make([]byte, 32) // XChaCha20-Poly1305 takes a 32-byte key
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	// Store key
	if err := e.keys.Put(ctx, contentID, key); err != nil {
		return nil, fmt.Errorf("failed to store key: %w", err)
	}

	var chunkIDs []core.ChunkID
	var totalSize uint64

	// Read all chunks, encrypt, hash, and store
	for chunk := range chunkCh {
		// Encrypt the chunk
		if err := e.encryptor.EncryptChunk(key, chunk); err != nil {
			return nil, fmt.Errorf("failed to encrypt chunk: %w", err)
		}

		// Hash the ciphertext to get the ChunkID (content-addressing)
		chunkHash := e.digest.Sum(chunk.Data)
		var chunkID core.ChunkID
		copy(chunkID[:], chunkHash[:])
		chunk.Header.ID = chunkID

		// Store the chunk
		if err := e.sink.PutChunk(ctx, chunk); err != nil {
			return nil, fmt.Errorf("failed to store chunk: %w", err)
		}

		chunkIDs = append(chunkIDs, chunkID)
		totalSize += uint64(chunk.Header.PlainSize)
	}

	if err := <-errCh; err != nil {
		return nil, fmt.Errorf("error reading stream: %w", err)
	}

	// For Milestone 7, MerkleRoot is just a hash of the concatenated chunk IDs
	// In the future this will be an actual Merkle Tree root.
	var idConcat []byte
	for _, id := range chunkIDs {
		idConcat = append(idConcat, id[:]...)
	}
	wholeHash := e.digest.Sum(idConcat)

	m := &manifest.Manifest{
		Version: 1,
		Descriptor: manifest.ContentDescriptor{
			ID:   contentID,
			Type: mtype,
			Size: totalSize,
		},
		ChunkIDs:   chunkIDs,
		MerkleRoot: wholeHash,
		WholeHash:  wholeHash,
		Crypto: manifest.CryptoDescriptor{
			Algorithm:      "XChaCha20-Poly1305",
			Version:        1,
			ChunkNonceSize: 24,
			KeyID:          "embedded",
		},
	}

	return m, nil
}

// Reassemble reads the manifest, fetches chunks, decrypts them, verifies integrity, and writes to w.
func (e *ContentEngine) Reassemble(ctx context.Context, m *manifest.Manifest, w io.Writer) error {
	// Retrieve key
	key, err := e.keys.Get(ctx, m.Descriptor.ID)
	if err != nil {
		return fmt.Errorf("failed to get content key: %w", err)
	}

	// Fetch all chunks, decrypt and verify
	// For simplicity in Milestone 7, we fetch sequentially.
	// But chunks can be fetched in parallel. We'll store them in a slice and sort by index.

	chunks := make([]*core.Chunk, 0, len(m.ChunkIDs))

	for _, chunkID := range m.ChunkIDs {
		chunk, err := e.source.GetChunk(ctx, chunkID)
		if err != nil {
			return fmt.Errorf("failed to get chunk %x: %w", chunkID, err)
		}

		// Verify chunk hash matches ID
		hash := e.digest.Sum(chunk.Data)
		if hash != core.Hash(chunkID) {
			return fmt.Errorf("corrupted chunk %x: hash mismatch", chunkID)
		}

		// Decrypt
		if err := e.encryptor.DecryptChunk(key, chunk); err != nil {
			return fmt.Errorf("failed to decrypt chunk %x: %w", chunkID, err)
		}

		chunks = append(chunks, chunk)
	}

	// Sort by index just in case they were fetched out of order
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Header.Index < chunks[j].Header.Index
	})

	// Write out
	for _, chunk := range chunks {
		if _, err := w.Write(chunk.Data); err != nil {
			return fmt.Errorf("failed to write decrypted chunk: %w", err)
		}
	}

	return nil
}

// -- Transport Layer APIs --

func (e *ContentEngine) GetChunk(ctx context.Context, id core.ChunkID) (*core.Chunk, error) {
	return e.source.GetChunk(ctx, id)
}

func (e *ContentEngine) PutChunk(ctx context.Context, chunk *core.Chunk) error {
	return e.sink.PutChunk(ctx, chunk)
}

func (e *ContentEngine) GetManifestBytes(ctx context.Context, id core.ContentID) ([]byte, error) {
	data, ok := e.manifests[id]
	if !ok {
		return nil, fmt.Errorf("manifest not found")
	}
	return data, nil
}

func (e *ContentEngine) PutManifestBytes(ctx context.Context, id core.ContentID, data []byte) error {
	e.manifests[id] = data
	return nil
}
