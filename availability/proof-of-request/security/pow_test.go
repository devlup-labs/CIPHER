package security

import (
    "crypto/sha256"
    "encoding/hex"
    "strconv"
    "strings"
    "testing"

    "proof-of-request/model"
	"proof-of-request/request"
)


func TestGeneratePoW(t *testing.T) {
    req := model.Request{
        FileID:    "file-123",
        ClientID:  "client-1",
        Timestamp: 123456789,
        Status:    model.Pending,
    }

    difficulty := uint8(2) //passing difficulty of 2 leading zeros in the hash

    pow := GeneratePoW(req, difficulty)

    if pow.Difficulty != difficulty {
        t.Fatalf("expected difficulty %d, got %d",
            difficulty, pow.Difficulty)
    }

    input := req.ClientID +
        req.FileID +
        strconv.FormatInt(req.Timestamp, 10) +
        strconv.FormatUint(pow.Nonce, 10)

    hash := sha256.Sum256([]byte(input))
    hashString := hex.EncodeToString(hash[:])

    prefix := strings.Repeat("0", int(difficulty))

    if !strings.HasPrefix(hashString, prefix) {
        t.Fatalf("invalid PoW: hash %s does not have prefix %s",
            hashString, prefix)
    }
}
// test checks that the GeneratePoW function returns a PoW with the correct difficulty and that the hash of the input string (constructed from the request fields and the nonce) has the required number of leading zeros.
func TestGeneratePoWDifficultyZero(t *testing.T) {
    req := model.Request{
        FileID:    "file-123",
        ClientID:  "client-1",
        Timestamp: 123456789,
        Status:    model.Pending,
    }

    pow := GeneratePoW(req, 0)

    if pow.Nonce != 0 {
        t.Fatalf("expected nonce 0 for difficulty 0, got %d",
            pow.Nonce)
    }

    if pow.Difficulty != 0 {
        t.Fatalf("expected difficulty 0, got %d",
            pow.Difficulty)
    }
}

//integration boundary test to check that the GeneratePoW function can be used in conjunction with the CreateRequest function to generate a valid PoW for a new request.
func TestGeneratePoWForNewRequest(t *testing.T) {
	req, err := request.CreateRequest("file-123", "client-456")
	if err != nil {
		t.Fatalf("CreateRequest() failed: %v", err)
	}

	difficulty := uint8(2)

	proof := GeneratePoW(req, difficulty)

	if proof.Difficulty != difficulty {
		t.Fatalf(
			"expected difficulty %d, got %d",
			difficulty,
			proof.Difficulty,
		)
	}
}