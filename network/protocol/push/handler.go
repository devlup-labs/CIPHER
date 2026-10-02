package push

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"cipher/network/content/core"
	"cipher/network/content/engine"
	"cipher/network/content/manifest"
	"cipher/network/content/verifier"
	"cipher/network/discovery"
	"cipher/network/protocol"
)

type AuthPolicy string

const (
	AuthPolicyOpen      AuthPolicy = "open"
	AuthPolicyAllowlist AuthPolicy = "allowlist"
)

type PendingSession struct {
	ContentID       core.ContentID
	Manifest        *manifest.Manifest
	ManifestBytes   []byte
	ExpectedChunks  map[core.ChunkID]struct{}
	CommittedChunks map[core.ChunkID]bool
	StartedAt       time.Time
	UpdatedAt       time.Time
}

type StreamHandler struct {
	host              host.Host
	engine            *engine.ContentEngine
	kdht              *dht.IpfsDHT
	digest            core.Digest
	allowPush         bool
	authPolicy        AuthPolicy
	allowedPublishers map[peer.ID]struct{}

	sessionsMu sync.RWMutex
	sessions   map[core.ContentID]*PendingSession
}

func NewStreamHandler(
	h host.Host,
	eng *engine.ContentEngine,
	kdht *dht.IpfsDHT,
	allowPush bool,
	authPolicy AuthPolicy,
	allowedPublishers []peer.ID,
) *StreamHandler {
	allowedMap := make(map[peer.ID]struct{})
	for _, pid := range allowedPublishers {
		allowedMap[pid] = struct{}{}
	}

	handler := &StreamHandler{
		host:              h,
		engine:            eng,
		kdht:              kdht,
		digest:            verifier.NewSHA256Digest(),
		allowPush:         allowPush,
		authPolicy:        authPolicy,
		allowedPublishers: allowedMap,
		sessions:          make(map[core.ContentID]*PendingSession),
	}

	h.SetStreamHandler(protocol.PushTransportProtocolID, handler.handleStream)
	return handler
}

func (h *StreamHandler) handleStream(s network.Stream) {
	defer s.Close()
	remotePeer := s.Conn().RemotePeer()

	// 1. Authorization check
	if !h.allowPush {
		log.Printf("[Push Protocol] Ingestion rejected from %s: push disabled (-allow-push=false)", remotePeer)
		_ = WritePushMessage(s, BuildPushError(PushStatusUnauthorized, "provider push disabled"))
		return
	}

	if h.authPolicy == AuthPolicyAllowlist {
		if _, ok := h.allowedPublishers[remotePeer]; !ok {
			log.Printf("[Push Protocol] Ingestion rejected from unauthorized publisher: %s", remotePeer)
			_ = WritePushMessage(s, BuildPushError(PushStatusUnauthorized, "publisher not in allowlist"))
			return
		}
	}

	log.Printf("[Push Protocol] Accepted push stream from %s", remotePeer)

	for {
		_ = s.SetReadDeadline(time.Now().Add(ReadTimeout))
		msg, err := ReadPushMessage(s)
		if err != nil {
			if err == io.EOF || err.Error() == "stream reset" {
				log.Printf("[Push Protocol] Push stream closed by %s", remotePeer)
				return
			}
			log.Printf("[Push Protocol] Error reading message: %v", err)
			return
		}
		_ = s.SetReadDeadline(time.Time{})

		if msg.Version != CurrentPushVersion {
			log.Printf("[Push Protocol] Unsupported message version: %d", msg.Version)
			_ = WritePushMessage(s, BuildPushError(PushStatusMalformed, "unsupported version"))
			return
		}

		switch msg.Type {
		case MsgPushManifest:
			h.handlePushManifest(s, msg)
		case MsgPushChunk:
			h.handlePushChunk(s, msg)
		case MsgPushBatchComplete:
			h.handlePushBatchComplete(s, msg)
		default:
			log.Printf("[Push Protocol] Unsupported message type: %d", msg.Type)
			_ = WritePushMessage(s, BuildPushError(PushStatusMalformed, "unknown message type"))
			return
		}
	}
}

func (h *StreamHandler) handlePushManifest(s network.Stream, msg *PushMessage) {
	contentID, assignedChunkIDs, manifestData, err := ParsePushManifest(msg.Payload)
	if err != nil {
		log.Printf("[Push Protocol] Failed to parse PUSH_MANIFEST: %v", err)
		_ = WritePushMessage(s, BuildPushError(PushStatusMalformed, "malformed manifest payload"))
		return
	}

	m, err := manifest.Deserialize(manifestData)
	if err != nil {
		log.Printf("[Push Protocol] Failed to deserialize manifest JSON: %v", err)
		_ = WritePushMessage(s, BuildPushError(PushStatusMalformed, "invalid manifest JSON"))
		return
	}

	if m.Descriptor.ID != contentID {
		log.Printf("[Push Protocol] Manifest contentID mismatch: %x vs %x", m.Descriptor.ID, contentID)
		_ = WritePushMessage(s, BuildPushError(PushStatusMalformed, "contentID mismatch"))
		return
	}

	expectedMap := make(map[core.ChunkID]struct{})
	for _, cid := range assignedChunkIDs {
		expectedMap[cid] = struct{}{}
	}

	h.sessionsMu.Lock()
	h.sessions[contentID] = &PendingSession{
		ContentID:       contentID,
		Manifest:        m,
		ManifestBytes:   manifestData,
		ExpectedChunks:  expectedMap,
		CommittedChunks: make(map[core.ChunkID]bool),
		StartedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	h.sessionsMu.Unlock()

	log.Printf("[Push Protocol] Initialized push session for ContentID %x (expecting %d chunks)", contentID, len(expectedMap))

	ack := BuildPushManifestAck(contentID, PushStatusOK)
	if err := WritePushMessage(s, ack); err != nil {
		log.Printf("[Push Protocol] Failed to write PUSH_MANIFEST_ACK: %v", err)
	}
}

func (h *StreamHandler) handlePushChunk(s network.Stream, msg *PushMessage) {
	contentID, chunk, err := ParsePushChunk(msg.Payload)
	if err != nil {
		log.Printf("[Push Protocol] Failed to parse PUSH_CHUNK: %v", err)
		_ = WritePushMessage(s, BuildPushError(PushStatusMalformed, "malformed chunk payload"))
		return
	}

	chunkID := chunk.Header.ID

	h.sessionsMu.RLock()
	session, exists := h.sessions[contentID]
	h.sessionsMu.RUnlock()

	if !exists {
		log.Printf("[Push Protocol] Chunk %x received for unknown/non-pending content %x", chunkID, contentID)
		_ = WritePushMessage(s, BuildPushChunkAck(chunkID, PushStatusNotInAssignedSet))
		return
	}

	// 1. Validate chunk is in assigned set
	if _, ok := session.ExpectedChunks[chunkID]; !ok {
		log.Printf("[Push Protocol] Chunk %x not in assigned set for ContentID %x", chunkID, contentID)
		_ = WritePushMessage(s, BuildPushChunkAck(chunkID, PushStatusNotInAssignedSet))
		return
	}

	// 2. Validate chunk belongs to manifest
	foundInManifest := false
	for _, cid := range session.Manifest.ChunkIDs {
		if cid == chunkID {
			foundInManifest = true
			break
		}
	}
	if !foundInManifest {
		log.Printf("[Push Protocol] Chunk %x not listed in manifest for ContentID %x", chunkID, contentID)
		_ = WritePushMessage(s, BuildPushChunkAck(chunkID, PushStatusMalformed))
		return
	}

	// 3. Verify SHA-256(ciphertext) matches ChunkID
	computedHash := h.digest.Sum(chunk.Data)
	if computedHash != core.Hash(chunkID) {
		log.Printf("[Push Protocol] Checksum mismatch for chunk %x", chunkID)
		_ = WritePushMessage(s, BuildPushChunkAck(chunkID, PushStatusHashMismatch))
		return
	}

	// 4. Idempotent commit: check if already exists in CAS
	ctx := context.Background()
	has, _ := h.engine.HasChunk(ctx, chunkID)
	if !has {
		if err := h.engine.PutChunk(ctx, chunk); err != nil {
			log.Printf("[Push Protocol] Failed to store chunk %x: %v", chunkID, err)
			_ = WritePushMessage(s, BuildPushChunkAck(chunkID, PushStatusIOError))
			return
		}
	}

	h.sessionsMu.Lock()
	session.CommittedChunks[chunkID] = true
	session.UpdatedAt = time.Now()
	h.sessionsMu.Unlock()

	ack := BuildPushChunkAck(chunkID, PushStatusOK)
	if err := WritePushMessage(s, ack); err != nil {
		log.Printf("[Push Protocol] Failed to write PUSH_CHUNK_ACK: %v", err)
	}
}

func (h *StreamHandler) handlePushBatchComplete(s network.Stream, msg *PushMessage) {
	contentID, err := ParsePushBatchComplete(msg.Payload)
	if err != nil {
		log.Printf("[Push Protocol] Failed to parse PUSH_BATCH_COMPLETE: %v", err)
		_ = WritePushMessage(s, BuildPushError(PushStatusMalformed, "malformed batch complete payload"))
		return
	}

	h.sessionsMu.Lock()
	session, exists := h.sessions[contentID]
	if !exists {
		h.sessionsMu.Unlock()
		log.Printf("[Push Protocol] BatchComplete requested for non-pending content %x", contentID)
		_ = WritePushMessage(s, BuildPushBatchCompleteAck(contentID, PushStatusIncomplete))
		return
	}

	// Invariant check: Did we commit ALL assigned chunks?
	allCommitted := true
	for cid := range session.ExpectedChunks {
		if !session.CommittedChunks[cid] {
			allCommitted = false
			break
		}
	}

	if !allCommitted || len(session.CommittedChunks) < len(session.ExpectedChunks) {
		h.sessionsMu.Unlock()
		log.Printf("[Push Protocol] BatchComplete rejected for %x: committed %d/%d assigned chunks",
			contentID, len(session.CommittedChunks), len(session.ExpectedChunks))
		_ = WritePushMessage(s, BuildPushBatchCompleteAck(contentID, PushStatusIncomplete))
		return
	}

	// Commit manifest to local CAS
	ctx := context.Background()
	if err := h.engine.PutManifestBytes(ctx, contentID, session.ManifestBytes); err != nil {
		h.sessionsMu.Unlock()
		log.Printf("[Push Protocol] Failed to store manifest for %x: %v", contentID, err)
		_ = WritePushMessage(s, BuildPushBatchCompleteAck(contentID, PushStatusIOError))
		return
	}

	// Remove session from pending
	delete(h.sessions, contentID)
	h.sessionsMu.Unlock()

	log.Printf("[Push Protocol] Content %x successfully committed to CAS (all %d assigned chunks verified)",
		contentID, len(session.ExpectedChunks))

	// Announce to DHT
	if h.kdht != nil {
		go func() {
			dhtCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := discovery.Provide(dhtCtx, h.kdht, contentID); err != nil {
				log.Printf("[Push Protocol] Warning: Failed to announce %x on DHT: %v", contentID, err)
			} else {
				log.Printf("[Push Protocol] [✓] Successfully announced ContentID %x on DHT", contentID)
			}
		}()
	}

	ack := BuildPushBatchCompleteAck(contentID, PushStatusOK)
	if err := WritePushMessage(s, ack); err != nil {
		log.Printf("[Push Protocol] Failed to write PUSH_BATCH_COMPLETE_ACK: %v", err)
	}
}
