package push

import (
	"bytes"
	"crypto/rand"
	"testing"

	"cipher/network/content/core"
)

func TestPushMessageFraming(t *testing.T) {
	origMsg := BuildPushError(0x05, "disk full error")
	buf := new(bytes.Buffer)

	if err := WritePushMessage(buf, origMsg); err != nil {
		t.Fatalf("WritePushMessage failed: %v", err)
	}

	readMsg, err := ReadPushMessage(buf)
	if err != nil {
		t.Fatalf("ReadPushMessage failed: %v", err)
	}

	if readMsg.Version != origMsg.Version {
		t.Errorf("version mismatch: got %d, want %d", readMsg.Version, origMsg.Version)
	}
	if readMsg.Type != origMsg.Type {
		t.Errorf("type mismatch: got %d, want %d", readMsg.Type, origMsg.Type)
	}

	code, str, err := ParsePushError(readMsg.Payload)
	if err != nil {
		t.Fatalf("ParsePushError failed: %v", err)
	}
	if code != 0x05 || str != "disk full error" {
		t.Errorf("error payload mismatch: got code=%d, msg=%s", code, str)
	}
}

func TestPushManifestSerialization(t *testing.T) {
	var contentID core.ContentID
	_, _ = rand.Read(contentID[:])

	var cid1, cid2, cid3 core.ChunkID
	_, _ = rand.Read(cid1[:])
	_, _ = rand.Read(cid2[:])
	_, _ = rand.Read(cid3[:])

	assigned := []core.ChunkID{cid1, cid2, cid3}
	manifestJSON := []byte(`{"version":1,"descriptor":{"id":"abc"}}`)

	msg := BuildPushManifest(contentID, assigned, manifestJSON)
	parsedCID, parsedAssigned, parsedJSON, err := ParsePushManifest(msg.Payload)
	if err != nil {
		t.Fatalf("ParsePushManifest failed: %v", err)
	}

	if parsedCID != contentID {
		t.Errorf("contentID mismatch")
	}
	if len(parsedAssigned) != len(assigned) {
		t.Fatalf("assigned chunk count mismatch: got %d, want %d", len(parsedAssigned), len(assigned))
	}
	for i := range assigned {
		if parsedAssigned[i] != assigned[i] {
			t.Errorf("chunk %d mismatch", i)
		}
	}
	if string(parsedJSON) != string(manifestJSON) {
		t.Errorf("manifest json mismatch: got %s, want %s", string(parsedJSON), string(manifestJSON))
	}
}

func TestPushChunkSerialization(t *testing.T) {
	var contentID core.ContentID
	var chunkID core.ChunkID
	_, _ = rand.Read(contentID[:])
	_, _ = rand.Read(chunkID[:])

	chunk := &core.Chunk{
		Header: core.ChunkHeader{
			Version:    1,
			ID:         chunkID,
			Index:      4,
			Offset:     1024,
			PlainSize:  512,
			CipherSize: 528,
		},
		Data: []byte("encrypted-ciphertext-data-payload-example"),
	}

	msg, err := BuildPushChunk(contentID, chunk)
	if err != nil {
		t.Fatalf("BuildPushChunk failed: %v", err)
	}

	parsedCID, parsedChunk, err := ParsePushChunk(msg.Payload)
	if err != nil {
		t.Fatalf("ParsePushChunk failed: %v", err)
	}

	if parsedCID != contentID {
		t.Errorf("contentID mismatch")
	}
	if parsedChunk.Header.ID != chunkID {
		t.Errorf("chunk ID mismatch")
	}
	if parsedChunk.Header.Index != 4 {
		t.Errorf("chunk index mismatch")
	}
	if !bytes.Equal(parsedChunk.Data, chunk.Data) {
		t.Errorf("chunk data mismatch")
	}
}

func TestFrameSizeLimits(t *testing.T) {
	// Attempt to create a message larger than MaxMessageSize
	hugePayload := make([]byte, MaxMessageSize+1)
	msg := &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushChunk,
		Payload: hugePayload,
	}

	buf := new(bytes.Buffer)
	err := WritePushMessage(buf, msg)
	if err == nil {
		t.Fatalf("expected error for oversized message, got nil")
	}
}
