package availability

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	availabilitytypes "cipher/availability/availability-contracts/types"
	verification "cipher/availability/availability-contracts/verification"
	"cipher/network/content/core"
	"cipher/network/content/manifest"

)

// ProviderProofEngine generates cryptographic chunk possession proofs in response
// to Availability challenges issued by publishers.
type ProviderProofEngine struct {
	mu           sync.RWMutex
	providerID   string
	privKey      ed25519.PrivateKey
	pubKey       ed25519.PublicKey
	chunkSource  core.ChunkSource
	trees        map[string]*MerkleTree // fileID -> MerkleTree
	chunkPayload map[string]map[int][]byte // fileID -> chunkIndex -> rawBytes
}

// NewProviderProofEngine initializes a proof engine for a provider node.
func NewProviderProofEngine(providerID string, privKey ed25519.PrivateKey, chunkSource core.ChunkSource) (*ProviderProofEngine, error) {
	if providerID == "" {
		return nil, errors.New("providerID is required")
	}
	if len(privKey) != ed25519.PrivateKeySize {
		return nil, errors.New("valid Ed25519 private key is required")
	}

	pubKey := privKey.Public().(ed25519.PublicKey)

	return &ProviderProofEngine{
		providerID:   providerID,
		privKey:      privKey,
		pubKey:       pubKey,
		chunkSource:  chunkSource,
		trees:        make(map[string]*MerkleTree),
		chunkPayload: make(map[string]map[int][]byte),
	}, nil
}

// PublicKey returns the provider's Ed25519 public key.
func (e *ProviderProofEngine) PublicKey() ed25519.PublicKey {
	return e.pubKey
}

// ProviderID returns the provider identifier.
func (e *ProviderProofEngine) ProviderID() string {
	return e.providerID
}

// RegisterFileChunks registers raw chunks for a file, builds its Merkle tree, and stores it.
func (e *ProviderProofEngine) RegisterFileChunks(fileID string, chunks [][]byte) (*MerkleTree, error) {
	if fileID == "" {
		return nil, errors.New("fileID cannot be empty")
	}
	if len(chunks) == 0 {
		return nil, errors.New("no chunks provided")
	}

	tree, err := BuildMerkleTreeFromChunks(chunks)
	if err != nil {
		return nil, fmt.Errorf("build Merkle tree: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.trees[fileID] = tree
	if _, ok := e.chunkPayload[fileID]; !ok {
		e.chunkPayload[fileID] = make(map[int][]byte)
	}
	for i, c := range chunks {
		chunkCopy := make([]byte, len(c))
		copy(chunkCopy, c)
		e.chunkPayload[fileID][i] = chunkCopy
	}

	return tree, nil
}

// RegisterTree registers an already constructed MerkleTree for a file.
func (e *ProviderProofEngine) RegisterTree(fileID string, tree *MerkleTree) error {
	if fileID == "" || tree == nil {
		return errors.New("fileID and tree cannot be empty")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.trees[fileID] = tree
	return nil
}

// MerkleRoot returns the 32-byte Merkle root for a registered file.
func (e *ProviderProofEngine) MerkleRoot(fileID string) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	tree, ok := e.trees[fileID]
	if !ok {
		return nil, fmt.Errorf("no tree registered for fileID %q", fileID)
	}
	return tree.Root(), nil
}

// HandleChallenge generates a complete, cryptographically verified ChallengeResponse
// for an issued Availability challenge.
func (e *ProviderProofEngine) HandleChallenge(challenge availabilitytypes.Challenge) (verification.ChallengeResponse, error) {
	return e.generateResponse(challenge, false, false)
}

// HandleChallengeAdversarial allows producing corrupted responses for adversarial testing.
func (e *ProviderProofEngine) HandleChallengeAdversarial(challenge availabilitytypes.Challenge, corruptChunkHash, corruptSig bool) (verification.ChallengeResponse, error) {
	return e.generateResponse(challenge, corruptChunkHash, corruptSig)
}

func (e *ProviderProofEngine) generateResponse(challenge availabilitytypes.Challenge, corruptHash, corruptSig bool) (verification.ChallengeResponse, error) {
	if challenge.ChallengeID == "" {
		return verification.ChallengeResponse{}, errors.New("challenge ID is empty")
	}

	e.mu.RLock()
	tree, ok := e.trees[challenge.FileID]
	e.mu.RUnlock()

	if !ok {
		// Attempt dynamic on-demand loading from chunkSource / ManifestStore
		if e.chunkSource != nil {
			if manifestStore, isStore := e.chunkSource.(core.ManifestStore); isStore {
				ctx := context.Background()
				idBytes, err := hex.DecodeString(challenge.FileID)
				if err == nil && len(idBytes) == 32 {
					var contentID core.ContentID
					copy(contentID[:], idBytes)
					if mData, err := manifestStore.GetManifestBytes(ctx, contentID); err == nil {
						if m, err := manifest.Deserialize(mData); err == nil {
							leafHashes := make([][]byte, len(m.ChunkIDs))
							for i, cid := range m.ChunkIDs {
								leafHashes[i] = make([]byte, 32)
								copy(leafHashes[i], cid[:])
							}
							newTree, err := NewMerkleTree(leafHashes)
							if err == nil {
								_ = e.RegisterTree(challenge.FileID, newTree)
								tree = newTree
								ok = true
							}
						}
					}
				}
			}
		}
	}

	if !ok || tree == nil {
		return verification.ChallengeResponse{}, fmt.Errorf("file %q not found on provider", challenge.FileID)
	}


	if challenge.ChunkID < 0 || challenge.ChunkID >= tree.LeafCount() {
		return verification.ChallengeResponse{}, fmt.Errorf("chunk index %d out of bounds [0, %d)", challenge.ChunkID, tree.LeafCount())
	}

	// 1. Generate Merkle audit proof
	proof, err := tree.GenerateProof(challenge.ChunkID)
	if err != nil {
		return verification.ChallengeResponse{}, fmt.Errorf("generate Merkle proof: %w", err)
	}

	// 2. Fetch leaf chunk hash
	chunkHash := make([]byte, 32)
	copy(chunkHash, tree.leaves[challenge.ChunkID])

	if corruptHash {
		chunkHash[0] ^= 0xFF
	}

	// 3. Compute deterministic binding digest: SHA256(ChunkHash || Nonce || ChallengeID)
	digest := verification.ComputeProofDigest(chunkHash, challenge.Nonce, challenge.ChallengeID)

	// 4. Sign digest with provider's Ed25519 private key
	sig := ed25519.Sign(e.privKey, digest)
	if corruptSig {
		sig[0] ^= 0xFF
	}

	return verification.ChallengeResponse{
		ChallengeID: challenge.ChallengeID,
		ContractID:  challenge.ContractID,
		ProviderID:  e.providerID,
		FileID:      challenge.FileID,
		ChunkID:     challenge.ChunkID,
		ChunkHash:   chunkHash,
		MerkleProof: proof,
		Signature:   sig,
		RespondedAt: time.Now().UTC(),
	}, nil
}
