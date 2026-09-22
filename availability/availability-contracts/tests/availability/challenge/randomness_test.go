package challenge_test

import (
	"bytes"
	"testing"

	challenge "cipher/availability/availability-contracts/challenge"
	availabilitytypes "cipher/availability/availability-contracts/types"
)

func TestGenerateNonce(t *testing.T) {
	epochID := availabilitytypes.EpochID("epoch-randomness")
	first, err := challenge.GenerateNonce(epochID, 0)
	if err != nil {
		t.Fatalf("GenerateNonce returned error: %v", err)
	}
	second, err := challenge.GenerateNonce(epochID, 1)
	if err != nil {
		t.Fatalf("GenerateNonce returned error: %v", err)
	}
	if len(first) != 32 || len(second) != 32 {
		t.Fatalf("nonce length = %d and %d, want 32", len(first), len(second))
	}
	if bytes.Equal(first, second) {
		t.Fatal("two generated nonces are equal")
	}
	if _, err := challenge.GenerateNonce("", 0); err == nil {
		t.Fatal("GenerateNonce accepted an empty epoch ID")
	}
}

func TestHashesAreDeterministicAndDomainSeparated(t *testing.T) {
	epochID := availabilitytypes.EpochID("epoch-hash")
	nonce := []byte("test nonce")
	chunkHash := challenge.GenerateChunkHash("provider-a", "file-a", epochID, 2, nonce)
	if got := challenge.GenerateChunkHash("provider-a", "file-a", epochID, 2, nonce); !bytes.Equal(chunkHash, got) {
		t.Fatal("chunk hash is not deterministic")
	}
	chunkChanges := [][]byte{
		challenge.GenerateChunkHash("provider-b", "file-a", epochID, 2, nonce),
		challenge.GenerateChunkHash("provider-a", "file-b", epochID, 2, nonce),
		challenge.GenerateChunkHash("provider-a", "file-a", "epoch-other", 2, nonce),
		challenge.GenerateChunkHash("provider-a", "file-a", epochID, 3, nonce),
		challenge.GenerateChunkHash("provider-a", "file-a", epochID, 2, []byte("other nonce")),
	}
	for _, changed := range chunkChanges {
		if bytes.Equal(chunkHash, changed) {
			t.Fatal("changing a chunk-hash input did not change the hash")
		}
	}
	roundHash := challenge.GenerateRoundHash("provider-a", "file-a", epochID, 2, nonce, 3)
	if got := challenge.GenerateRoundHash("provider-a", "file-a", epochID, 2, nonce, 3); !bytes.Equal(roundHash, got) {
		t.Fatal("round hash is not deterministic")
	}
	if bytes.Equal(roundHash, challenge.GenerateRoundHash("provider-a", "file-a", epochID, 2, nonce, 4)) {
		t.Fatal("changing round number did not change round hash")
	}
	if bytes.Equal(roundHash, challenge.GenerateRoundHash("provider-b", "file-a", epochID, 2, nonce, 3)) {
		t.Fatal("changing provider ID did not change round hash")
	}
	if bytes.Equal(chunkHash, roundHash) {
		t.Fatal("CHUNK and ROUND hashes are not domain separated")
	}
}
