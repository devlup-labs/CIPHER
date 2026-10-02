package identity

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestEthereumIdentity_LoadOrCreate(t *testing.T) {
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "ethereum.key")

	// 1. Initial generation
	priv1, err := LoadOrCreateEthereumKey(keyPath)
	if err != nil {
		t.Fatalf("Failed to generate first Ethereum key: %v", err)
	}
	if priv1 == nil {
		t.Fatal("Generated key is nil")
	}

	addr1 := EthereumAddress(priv1)
	if addr1 == (common.Address{}) {
		t.Fatal("Derived Ethereum address is zero")
	}

	// Verify file permissions
	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("Failed to stat key file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}

	// 2. Reloading existing key
	priv2, err := LoadOrCreateEthereumKey(keyPath)
	if err != nil {
		t.Fatalf("Failed to reload Ethereum key: %v", err)
	}

	addr2 := EthereumAddress(priv2)
	if addr1 != addr2 {
		t.Fatalf("Address mismatch on reload: %s vs %s", addr1.Hex(), addr2.Hex())
	}

	// Check private key bytes match
	if fmt.Sprintf("%x", crypto.FromECDSA(priv1)) != fmt.Sprintf("%x", crypto.FromECDSA(priv2)) {
		t.Fatal("Private key scalar mismatch on reload")
	}
}

func TestEthereumIdentity_ParseKey(t *testing.T) {
	// Generate reference key
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	hexRaw := fmt.Sprintf("%x", crypto.FromECDSA(priv))
	hexWith0x := "0x" + hexRaw

	// Parse without 0x
	parsed1, err := ParseEthereumKey(hexRaw)
	if err != nil {
		t.Fatalf("Failed to parse raw hex key: %v", err)
	}
	if EthereumAddress(parsed1) != crypto.PubkeyToAddress(priv.PublicKey) {
		t.Fatal("Address mismatch for raw hex")
	}

	// Parse with 0x
	parsed2, err := ParseEthereumKey(hexWith0x)
	if err != nil {
		t.Fatalf("Failed to parse 0x-prefixed hex key: %v", err)
	}
	if EthereumAddress(parsed2) != crypto.PubkeyToAddress(priv.PublicKey) {
		t.Fatal("Address mismatch for 0x-prefixed hex")
	}

	// Invalid hex
	_, err = ParseEthereumKey("not-a-valid-hex-string")
	if err == nil {
		t.Fatal("Expected error parsing invalid hex, got nil")
	}
}

func TestEthereumIdentity_NilKeyHandling(t *testing.T) {
	var nilPriv *ecdsa.PrivateKey
	addr := EthereumAddress(nilPriv)
	if addr != (common.Address{}) {
		t.Fatalf("Expected zero address for nil key, got %s", addr.Hex())
	}
}
