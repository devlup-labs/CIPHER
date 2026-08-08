package verifier

import (
	"bytes"
	"crypto/sha256"

	"cipher/internal/content/core"
)

// SHA256Digest implements core.Digest using SHA-256.
type SHA256Digest struct{}

func NewSHA256Digest() *SHA256Digest {
	return &SHA256Digest{}
}

func (d *SHA256Digest) Sum(data []byte) core.Hash {
	h := sha256.Sum256(data)
	var hash core.Hash
	copy(hash[:], h[:])
	return hash
}

func (d *SHA256Digest) Verify(data []byte, hash core.Hash) bool {
	h := sha256.Sum256(data)
	return bytes.Equal(h[:], hash[:])
}

func (d *SHA256Digest) Algorithm() string {
	return "sha256"
}
