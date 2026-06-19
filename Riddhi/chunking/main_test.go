package main

import "testing"
import"os"

func TestGenerateHash(t *testing.T) {

    data := []byte("hello")

    hash1 := GenerateHash(data)

    hash2 := GenerateHash(data)

    if hash1 != hash2 {
        t.Errorf("hashes should match")
    }

// === RUN   TestGenerateHash
// --- PASS: TestGenerateHash (0.00s)
// PASS
// ok      riddhi/chunking 0.005s
}

func TestCommitmentVerification(t *testing.T) {

    data := []byte("hello world")

    key, commitment, err := GenerateCommitment(data)

    if err != nil {
        t.Fatal(err)
    }

    valid := VerifyCommitment(
        key,
        data,
        commitment,
    )

    if !valid {
        t.Errorf("commitment verification failed")
    }

// 	=== RUN   TestCommitmentVerification
// --- PASS: TestCommitmentVerification (0.00s)
// PASS
// ok      riddhi/chunking 0.005s
}

func TestCommitmentTampering(t *testing.T) {

    data := []byte("hello world")

    key, commitment, err := GenerateCommitment(data)

    if err != nil {
        t.Fatal(err)
    }

    tampered := []byte("hello worlds")

    valid := VerifyCommitment(
        key,
        tampered,
        commitment,
    )

    if valid {
        t.Errorf("tampered data should fail verification")
    }
// 	=== RUN   TestCommitmentTampering
// --- PASS: TestCommitmentTampering (0.00s)
// PASS
// ok      riddhi/chunking 0.004s
}

func TestLeafDeterminism(t *testing.T) {

    data := []byte("hello")

    leaf1 := GenerateLeafHash(
        "file123",
        0,
        len(data),
        data,
    )

    leaf2 := GenerateLeafHash(
        "file123",
        0,
        len(data),
        data,
    )

    if leaf1 != leaf2 {
        t.Errorf("leaf hashes should match")
    }

// 	=== RUN   TestLeafDeterminism
// --- PASS: TestLeafDeterminism (0.00s)
// PASS
// ok      riddhi/chunking 0.005s
}

func TestLeafChanges(t *testing.T) {

    data := []byte("hello")

    leaf1 := GenerateLeafHash(
        "file123",
        0,
        len(data),
        data,
    )

    leaf2 := GenerateLeafHash(
        "file123",
        1,
        len(data),
        data,
    )

    if leaf1 == leaf2 {
        t.Errorf("different index should change leaf")
    }

// 	=== RUN   TestLeafChanges
// --- PASS: TestLeafChanges (0.00s)
// PASS
// ok      riddhi/chunking 0.006s
}

func TestChunkFile(t *testing.T) {

    content := make([]byte, 70000)

    err := os.WriteFile(
        "testfile.bin",
        content,
        0644,
    )

    if err != nil {
        t.Fatal(err)
    }

    defer os.Remove("testfile.bin")

    chunker := DefaultFileChunker{
        chunkSize: 32768,
    }

    chunks, err := chunker.ChunkFile("testfile.bin")

    if err != nil {
        t.Fatal(err)
    }

    if len(chunks) != 3 {
        t.Errorf("expected 3 chunks got %d",
            len(chunks))
    }

// 	=== RUN   TestChunkFile
// --- PASS: TestChunkFile (0.00s)
// PASS
// ok      riddhi/chunking 0.007s
}