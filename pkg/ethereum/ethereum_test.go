package ethereum

import (
	"testing"
	"path/filepath"
	"os"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// newPeerID creates a fresh libp2p PeerID for testing.
func newPeerID(t *testing.T) peer.ID {
	priv, _, err := libp2pcrypto.GenerateEd25519Key(nil)
	if err != nil {
		t.Fatal(err)
	}

	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}

	return id
}

// TestValidBinding verifies the happy path.

// Ensures that a correctly signed binding between a libp2p PeerID and an
// Ethereum wallet is accepted by the verifier. This is the expected behaviour
// when a legitimate provider registers.
func TestValidBinding(t *testing.T) {

	pid := newPeerID(t)

	ethKey, _ := gethcrypto.GenerateKey()

	req, err := SignBinding(pid, ethKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := VerifyBinding(req); err != nil {
		t.Fatal(err)
	}
}

// TestModifiedPeerID ensures that the PeerID cannot be changed after signing.
// Prevents an attacker from taking a valid signature and replacing the PeerID
// with another one to falsely claim ownership of a different libp2p identity.
func TestModifiedPeerID(t *testing.T) {

	pid := newPeerID(t)

	ethKey, _ := gethcrypto.GenerateKey()

	req, _ := SignBinding(pid, ethKey)

	another := newPeerID(t)

	req.PeerID = another.String()

	if VerifyBinding(req) == nil {
		t.Fatal("verification should fail")
	}
}

// TestWrongWallet ensures that the wallet address is part of the signed data.
// Prevents an attacker from replacing the registered Ethereum wallet while
// keeping the original signature.
func TestWrongWallet(t *testing.T) {

	pid := newPeerID(t)

	ethKey, _ := gethcrypto.GenerateKey()

	req, _ := SignBinding(pid, ethKey)

	otherKey, _ := gethcrypto.GenerateKey()

	req.Wallet = gethcrypto.PubkeyToAddress(otherKey.PublicKey).Hex()

	if VerifyBinding(req) == nil {
		t.Fatal("verification should fail")
	}
}

// TestCorruptedSignature ensures that even a small modification to the
// signature invalidates the proof.
func TestCorruptedSignature(t *testing.T) {

	pid := newPeerID(t)

	ethKey, _ := gethcrypto.GenerateKey()

	req, _ := SignBinding(pid, ethKey)

	req.Signature = req.Signature[:10] + "ff" + req.Signature[12:]

	if VerifyBinding(req) == nil {
		t.Fatal("verification should fail")
	}
}

// TestLoadOrGenerateEthereumKey_GeneratesNewKey verifies that a new
// Ethereum private key is generated and persisted when no key file exists.
func TestLoadOrGenerateEthereumKey_GeneratesNewKey(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "provider_eth.key")

	key, err := LoadOrGenerateEthereumKey(keyPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if key == nil {
		t.Fatal("expected generated key, got nil")
	}

	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("expected key file to exist: %v", err)
	}
}

// TestLoadOrGenerateEthereumKey_LoadsExistingKey verifies that an existing
// Ethereum key is loaded instead of generating a new one.
// The provider's Ethereum identity must remain stable across restarts.
// Loading the same private key ensures the wallet address does not change.
func TestLoadOrGenerateEthereumKey_LoadsExistingKey(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "provider_eth.key")

	key1, err := LoadOrGenerateEthereumKey(keyPath)
	if err != nil {
		t.Fatal(err)
	}

	key2, err := LoadOrGenerateEthereumKey(keyPath)
	if err != nil {
		t.Fatal(err)
	}

	addr1 := gethcrypto.PubkeyToAddress(key1.PublicKey)
	addr2 := gethcrypto.PubkeyToAddress(key2.PublicKey)

	if addr1 != addr2 {
		t.Fatalf("expected same wallet, got %s and %s", addr1.Hex(), addr2.Hex())
	}
}

// TestLoadOrGenerateEthereumKey_InvalidKeyFile verifies that loading fails
// when the key file contains invalid data.
//
// Why this test?
// Prevents corrupted or malformed key files from being silently accepted.
func TestLoadOrGenerateEthereumKey_InvalidKeyFile(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "provider_eth.key")

	if err := os.WriteFile(keyPath, []byte("not-a-valid-key"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadOrGenerateEthereumKey(keyPath); err == nil {
		t.Fatal("expected error for invalid key file")
	}
}

//helper
func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	providers = make(map[string]*Provider)
}

// TestRegisterProvider verifies that a provider is successfully added
// to the in-memory registry after registration.
func TestRegisterProvider(t *testing.T) {
	resetRegistry()
	peerID := "peer1"
	wallet := "0x123"

	if err := Register(peerID, wallet); err != nil {
		t.Fatal(err)
	}

	p, err := GetProvider(peerID)
	if err != nil {
		t.Fatal(err)
	}

	if p.PeerID != peerID {
		t.Fatalf("expected peerID %s, got %s", peerID, p.PeerID)
	}

	if p.Wallet != wallet {
		t.Fatalf("expected wallet %s, got %s", wallet, p.Wallet)
	}

	if p.Active {
		t.Fatal("provider should not be active before staking")
	}
}

// TestDuplicateRegister verifies that the same provider cannot be
// registered twice.
func TestDuplicateRegister(t *testing.T) {
	resetRegistry()
	peerID := "peer2"
	wallet := "0xabc"

	if err := Register(peerID, wallet); err != nil {
		t.Fatal(err)
	}

	if err := Register(peerID, wallet); err == nil {
		t.Fatal("expected duplicate registration to fail")
	}
}

// TestStake verifies that staking activates the provider and stores
// the correct stake amount.
func TestStake(t *testing.T) {
	resetRegistry()
	peerID := "peer3"

	if err := Register(peerID, "0x456"); err != nil {
		t.Fatal(err)
	}

	if err := Stake(peerID, DefaultStake); err != nil {
		t.Fatal(err)
	}

	p, err := GetProvider(peerID)
	if err != nil {
		t.Fatal(err)
	}

	if !p.Active {
		t.Fatal("expected provider to become active after staking")
	}

	if p.Stake != DefaultStake {
		t.Fatalf("expected stake %d, got %d", DefaultStake, p.Stake)
	}
}

// TestSlashStake verifies that slashing decreases the provider's stake.
func TestSlashStake(t *testing.T) {

	resetRegistry()
	peerID := "peer4"

	Register(peerID, "0x789")
	Stake(peerID, 100)

	if err := SlashStake(peerID, 10); err != nil {
		t.Fatal(err)
	}

	p, err := GetProvider(peerID)
	if err != nil {
		t.Fatal(err)
	}

	if p.Stake != 90 {
		t.Fatalf("expected stake 90, got %d", p.Stake)
	}
}

// TestSlashStakeToZero verifies that a provider becomes inactive once
// its entire stake has been slashed.
func TestSlashStakeToZero(t *testing.T) {

	resetRegistry()
	peerID := "peer5"

	Register(peerID, "0x999")
	Stake(peerID, 50)

	if err := SlashStake(peerID, 50); err != nil {
		t.Fatal(err)
	}

	p, err := GetProvider(peerID)
	if err != nil {
		t.Fatal(err)
	}

	if p.Stake != 0 {
		t.Fatalf("expected stake 0, got %d", p.Stake)
	}

	if p.Active {
		t.Fatal("provider should become inactive after stake reaches zero")
	}
}