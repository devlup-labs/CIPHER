package verifier_test

import (
	"crypto/rand"
	"testing"

	"cipher/internal/content/verifier"
)

func BenchmarkSHA256Digest(b *testing.B) {
	dig := verifier.NewSHA256Digest()
	
	chunkSize := 256 * 1024
	data := make([]byte, chunkSize)
	rand.Read(data)

	b.SetBytes(int64(chunkSize))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dig.Sum(data)
	}
}
