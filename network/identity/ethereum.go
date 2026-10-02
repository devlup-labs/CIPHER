package identity

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// LoadOrCreateEthereumKey loads an existing Secp256k1 private key from keyPath or generates and saves one.
func LoadOrCreateEthereumKey(keyPath string) (*ecdsa.PrivateKey, error) {
	if keyPath == "" {
		configDir := os.Getenv("CIPHER_CONFIG_DIR")
		if configDir == "" {
			var err error
			configDir, err = os.UserConfigDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user config dir: %w", err)
			}
		}
		keyPath = filepath.Join(configDir, "cipher", "ethereum.key")
	}

	// Try to load existing key
	if _, err := os.Stat(keyPath); err == nil {
		hexBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read ethereum key file: %w", err)
		}

		keyHex := strings.TrimSpace(string(hexBytes))
		return ParseEthereumKey(keyHex)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat ethereum key file: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(keyPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate a new Secp256k1 key
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ethereum key: %w", err)
	}

	hexKey := crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	_ = hexKey
	keyBytes := crypto.FromECDSA(privKey)
	keyHex := fmt.Sprintf("%x\n", keyBytes)

	if err := os.WriteFile(keyPath, []byte(keyHex), 0600); err != nil {
		return nil, fmt.Errorf("failed to write ethereum key file: %w", err)
	}

	return privKey, nil
}

// ParseEthereumKey parses a hex-encoded private key (with or without 0x prefix).
func ParseEthereumKey(keyHex string) (*ecdsa.PrivateKey, error) {
	keyHex = strings.TrimSpace(keyHex)
	if strings.HasPrefix(keyHex, "0x") || strings.HasPrefix(keyHex, "0X") {
		keyHex = keyHex[2:]
	}
	privKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid ethereum private key hex: %w", err)
	}
	return privKey, nil
}

// EthereumAddress returns the 20-byte Ethereum address corresponding to the private key.
func EthereumAddress(privKey *ecdsa.PrivateKey) common.Address {
	if privKey == nil {
		return common.Address{}
	}
	return crypto.PubkeyToAddress(privKey.PublicKey)
}
