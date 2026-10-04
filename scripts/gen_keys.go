package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func generateAndSave(name, path string) (peer.ID, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	priv, pub, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return "", err
	}
	bytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return "", err
	}
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return "", err
	}
	fmt.Printf("%s Key saved to: %s\n", name, path)
	fmt.Printf("%s Peer ID: %s\n\n", name, pid.String())
	return pid, nil
}

func main() {
	_, err := generateAndSave("BOOTSTRAP", "config/identity/bootstrap.key")
	if err != nil {
		panic(err)
	}
	_, err = generateAndSave("RELAY", "config/identity/relay.key")
	if err != nil {
		panic(err)
	}
}
