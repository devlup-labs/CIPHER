package model

// PoW represents a proof of work solution for a specific request.
// The difficulty means the number of leading zero hexadecimal characters in the SHA-256 hash.
type PoW struct {
	Nonce      uint64
	Difficulty uint8
}