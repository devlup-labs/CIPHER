package chunk

const (
	// ProtocolVersion identifies the current chunk protocol format.
	ProtocolVersion uint8 = 1
)

const (
	// StandardChunkSize is the normal plaintext chunk size.
	StandardChunkSize = 32 * 1024 // 32 KiB

	// ChaCha20NonceSize is the nonce size used by standard ChaCha20-Poly1305.
	ChaCha20NonceSize = 12

	// Poly1305TagSize is the authentication-tag size.
	Poly1305TagSize = 16

	// EncryptionOverhead is added to every encrypted chunk.
	EncryptionOverhead = ChaCha20NonceSize + Poly1305TagSize

	// MaxCiphertextSize is the largest encrypted chunk accepted.
	MaxCiphertextSize = StandardChunkSize + EncryptionOverhead
)

const (
	// Cryptographic fields based on 256-bit hashes and keys.
	HashSize       = 32
	ContentIDSize  = 32
	FileIDSize     = 32
	MerkleRootSize = 32
	ChunkKeySize   = 32

	// RequestIDSize may be changed if your current protocol uses
	// another request-ID format.
	RequestIDSize = 16
)

const (
	// MaxMerkleProofNodes limits Merkle proof depth.
	//
	// A depth of 64 is already enough for an extremely large number
	// of chunks and prevents attackers from supplying huge proofs.
	MaxMerkleProofNodes = 64

	// MaxMerkleProofSize is the maximum encoded proof-hash data.
	MaxMerkleProofSize = MaxMerkleProofNodes * HashSize
)

const (
	// Maximum attacker-controlled text sizes.
	MaxFilenameSize     = 255
	MaxMIMETypeSize     = 128
	MaxMetadataSize     = 2 * 1024
	MaxErrorMessageSize = 512
)

const (
	// ChunkRequest should contain only identifiers, index and
	// small protocol metadata.
	MaxChunkRequestSize = 512

	// MaxFrameSize matches the existing ReadMessage frame cap.
	MaxFrameSize = 2 * 1024 * 1024

	// MaxMessagePayloadSize subtracts the version and type fields that live
	// inside the length-prefixed frame.
	MaxMessagePayloadSize = MaxFrameSize - 3

	// ChunkResponse carries a binary ChunkHeader plus encrypted chunk data.
	MaxChunkResponseSize = MaxMessagePayloadSize

	// KeyReveal normally contains a key and small identifying fields.
	MaxKeyRevealSize = 256

	// MaxProtocolErrorSize limits encoded protocol error responses.
	MaxProtocolErrorSize = MaxErrorMessageSize + 128

	// Manifest responses can be larger than request payloads but remain bound
	// by the frame cap.
	MaxManifestSize = MaxMessagePayloadSize
)

const (
	// The initial hardened protocol allows one chunk transaction
	// per stream.
	MaxChunksPerRequest      = 1
	MaxTransactionsPerStream = 1
	MaxMessagesPerStream     = 4
)

const (
	// MaxSupportedChunkCount prevents absurd chunk indexes and
	// unbounded manifest or proof calculations.
	MaxSupportedChunkCount uint64 = 1 << 32

	// MaxSupportedChunkIndex is inclusive.
	MaxSupportedChunkIndex = MaxSupportedChunkCount - 1

	// MaxSupportedFileSize is approximately 128 TiB when using
	// 32 KiB chunks.
	MaxSupportedFileSize uint64 = MaxSupportedChunkCount * uint64(StandardChunkSize)
)

// MaxPayloadSizeForMessage returns the maximum payload size accepted for a
// particular chunk protocol message type.
//
// Rename these message constants if messages.go currently uses different names.
func MaxPayloadSizeForMessage(messageType MessageType) int {
	switch messageType {
	case MsgRequestManifest:
		return MaxChunkRequestSize

	case MsgManifest:
		return MaxManifestSize

	case MsgRequestChunk:
		return MaxChunkRequestSize

	case MsgChunk:
		return MaxChunkResponseSize

	case MsgAck:
		return HashSize + 1

	case MsgError:
		return MaxProtocolErrorSize

	default:
		return 0
	}
}
