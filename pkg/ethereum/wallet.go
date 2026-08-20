package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/1amKhush/CIPHER/pkg/logger"
)

// LoadOrGenerateEthereumKey loads an Ethereum private key from a file,
// or generates a new one and saves it if the file doesn't exist.
func LoadOrGenerateEthereumKey(path string) (*ecdsa.PrivateKey, error) {
	if path == "" {
		// Generate an ephemeral key if no path is provided.
		return gethcrypto.GenerateKey()
	}

	keyData, err := os.ReadFile(path)
	if err == nil {
		// Remove newline if present.
		keyHex := strings.TrimSpace(string(keyData))

		keyBytes, err := hex.DecodeString(keyHex)
		if err != nil {
			return nil, fmt.Errorf("failed to decode ethereum key: %w", err)
		}

		privKey, err := gethcrypto.ToECDSA(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ethereum key: %w", err)
		}

		return privKey, nil
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read ethereum key file: %w", err)
	}

	// Generate new key.
	logger.Info().
		Str("path", path).
		Msg("Generating new Ethereum private key")

	privKey, err := gethcrypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ethereum key: %w", err)
	}

	keyBytes := gethcrypto.FromECDSA(privKey)
	keyHex := hex.EncodeToString(keyBytes)

	if err := os.WriteFile(path, []byte(keyHex), 0600); err != nil {
		return nil, fmt.Errorf("failed to save ethereum private key: %w", err)
	}

	return privKey, nil
}

// LoadEthereumKey loads an existing Ethereum private key.
func LoadEthereumKey(path string) (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	logger.Info().
        Str("path", path).
        Msg("Loaded existing Ethereum private key")
	if err != nil {
		return nil, fmt.Errorf("failed to read ethereum key: %w", err)
	}

	keyHex := strings.TrimSpace(string(keyData))

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ethereum key: %w", err)
	}

	privKey, err := gethcrypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ethereum key: %w", err)
	}

	return privKey, nil
}

// SaveEthereumKey saves an Ethereum private key to disk.
func SaveEthereumKey(path string, privKey *ecdsa.PrivateKey) error {
	keyBytes := gethcrypto.FromECDSA(privKey)
	keyHex := hex.EncodeToString(keyBytes)

	if err := os.WriteFile(path, []byte(keyHex), 0600); err != nil {
		return fmt.Errorf("failed to save ethereum key: %w", err)
	}

	return nil
}