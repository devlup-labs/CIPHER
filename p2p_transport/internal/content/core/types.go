package core

import "context"

type ChunkID [32]byte
type ContentID [32]byte
type Hash [32]byte

type ChunkHeader struct {
	Version    uint16
	ID         ChunkID
	Index      uint32
	Offset     int64
	PlainSize  uint32
	CipherSize uint32
	Nonce      [12]byte
}

type Chunk struct {
	Header ChunkHeader
	Data   []byte // Represents current payload (plaintext or ciphertext)
}

type ChunkSource interface {
	HasChunk(ctx context.Context, id ChunkID) (bool, error)
	GetChunk(ctx context.Context, id ChunkID) (*Chunk, error)
}

type ChunkSink interface {
	PutChunk(ctx context.Context, chunk *Chunk) error
}

type Digest interface {
	Sum(data []byte) Hash
	Verify(data []byte, hash Hash) bool
	Algorithm() string
}

type Encryptor interface {
	EncryptChunk(key []byte, chunk *Chunk) error
	DecryptChunk(key []byte, chunk *Chunk) error
}

type KeyProvider interface {
	Get(ctx context.Context, id ContentID) ([]byte, error)
	Put(ctx context.Context, id ContentID, key []byte) error
	Delete(ctx context.Context, id ContentID) error
}

type Scheduler interface{}

type EngineConfig struct {
	ChunkSize uint32
}
