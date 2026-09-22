package security

//
import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"proof-of-request/model"
)

func GeneratePoW(req model.Request, difficulty uint8) model.PoW {
	nonce := uint64(0)

	prefix := strings.Repeat("0", int(difficulty))

	for {
		input := req.ClientID +
			req.FileID +
			strconv.FormatInt(req.Timestamp, 10) +
			strconv.FormatUint(nonce, 10)

		hash := sha256.Sum256([]byte(input))
		hashString := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hashString, prefix) {
			return model.PoW{
				Nonce:      nonce,
				Difficulty: difficulty,
			}
		}

		nonce++
	}
}

func VerifyPoW(req model.Request, proof model.PoW) bool {
	input := req.ClientID +
		req.FileID +
		strconv.FormatInt(req.Timestamp, 10) +
		strconv.FormatUint(proof.Nonce, 10)

	hash := sha256.Sum256([]byte(input))
	hashString := hex.EncodeToString(hash[:])

	prefix := strings.Repeat("0", int(proof.Difficulty))

	return strings.HasPrefix(hashString, prefix)
}
