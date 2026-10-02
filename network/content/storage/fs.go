package storage

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cipher/network/content/core"
)

// FSStorage implements core.ChunkSource and core.ChunkSink using local filesystem.
type FSStorage struct {
	baseDir string
}

func NewFSStorage(baseDir string) error {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}
	return nil
}

func NewFSStore(baseDir string) *FSStorage {
	return &FSStorage{baseDir: baseDir}
}

func (s *FSStorage) pathForChunk(id core.ChunkID) string {
	encoded := hex.EncodeToString(id[:])
	if len(encoded) >= 4 {
		// Shard by first 2 bytes (4 hex chars), e.g. ab/cd/abcdef...
		return filepath.Join(s.baseDir, encoded[0:2], encoded[2:4], encoded)
	}
	return filepath.Join(s.baseDir, encoded)
}

func (s *FSStorage) HasChunk(ctx context.Context, id core.ChunkID) (bool, error) {
	_, err := os.Stat(s.pathForChunk(id))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *FSStorage) PutChunk(ctx context.Context, chunk *core.Chunk) error {
	path := s.pathForChunk(chunk.Header.ID)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create shard dir: %w", err)
	}

	tmpPath := fmt.Sprintf("%s.%d.tmp", path, os.Getpid())
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	// Serialize Header
	if err := binary.Write(f, binary.LittleEndian, &chunk.Header); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write chunk header: %w", err)
	}

	// Write Data
	if _, err := f.Write(chunk.Data); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write chunk data: %w", err)
	}

	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to fsync chunk: %w", err)
	}
	f.Close()

	return os.Rename(tmpPath, path)
}

func (s *FSStorage) GetChunk(ctx context.Context, id core.ChunkID) (*core.Chunk, error) {
	path := s.pathForChunk(id)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	chunk := &core.Chunk{}
	if err := binary.Read(f, binary.LittleEndian, &chunk.Header); err != nil {
		return nil, fmt.Errorf("failed to read chunk header: %w", err)
	}

	// Calculate data size from file info minus header size, or use chunk.Header.CipherSize
	// Note: It's either PlainSize or CipherSize depending on if it's encrypted.
	// But actually, we just read the rest of the file.
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk data: %w", err)
	}

	// Validation: length of data should match either CipherSize or PlainSize
	// (usually CipherSize since it's stored encrypted).
	// We won't enforce strictly here since the Engine decryptor will validate it.
	chunk.Data = data

	return chunk, nil
}

func (s *FSStorage) manifestPath(id core.ContentID) string {
	encoded := hex.EncodeToString(id[:])
	return filepath.Join(s.baseDir, "manifests", encoded+".json")
}

func (s *FSStorage) GetManifestBytes(ctx context.Context, id core.ContentID) ([]byte, error) {
	path := s.manifestPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("manifest not found: %w", err)
	}
	return data, nil
}

func (s *FSStorage) PutManifestBytes(ctx context.Context, id core.ContentID, data []byte) error {
	path := s.manifestPath(id)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create manifest dir: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	return nil
}

func (s *FSStorage) ListManifests(ctx context.Context) ([]core.ContentID, error) {
	manifestDir := filepath.Join(s.baseDir, "manifests")
	var ids []core.ContentID

	entries, err := os.ReadDir(manifestDir)
	if err != nil {
		if os.IsNotExist(err) {
			return ids, nil
		}
		return nil, fmt.Errorf("failed to read manifest dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()
		encoded := name[:len(name)-5] // Strip .json
		
		idBytes, err := hex.DecodeString(encoded)
		if err != nil || len(idBytes) != 32 {
			continue
		}

		var id core.ContentID
		copy(id[:], idBytes)
		ids = append(ids, id)
	}

	return ids, nil
}
