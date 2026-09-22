package identity

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
)

// LoadOrCreate loads an existing Ed25519 private key from the user's config directory,
// or generates a new one and saves it if it doesn't exist.
func LoadOrCreate() (crypto.PrivKey, error) {
	configDir := os.Getenv("CIPHER_CONFIG_DIR")
	if configDir == "" {
		var err error
		configDir, err = os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user config dir: %w", err)
		}
	}

	appDir := filepath.Join(configDir, "cipher")
	keyPath := filepath.Join(appDir, "identity.key")
	return LoadOrCreateFromPath(keyPath)
}

// LoadOrCreateFromPath loads an existing Ed25519 key from keyPath or generates and saves one.
func LoadOrCreateFromPath(keyPath string) (crypto.PrivKey, error) {
	// Try to load existing key
	if _, err := os.Stat(keyPath); err == nil {
		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %w", err)
		}

		priv, err := crypto.UnmarshalPrivateKey(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal private key: %w", err)
		}

		return priv, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat key file: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(keyPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Generate a new Ed25519 key
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Marshal and save
	keyBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := os.WriteFile(keyPath, keyBytes, 0600); err != nil {
		return nil, fmt.Errorf("failed to write key file: %w", err)
	}

	return priv, nil
}

// GenerateEphemeral generates a one-time in-memory Ed25519 private key without persisting to disk.
func GenerateEphemeral() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral key pair: %w", err)
	}
	return priv, nil
}
