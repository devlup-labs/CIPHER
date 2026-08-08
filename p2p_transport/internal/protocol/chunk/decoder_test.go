package chunk_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"cipher/internal/content/core"
	"cipher/internal/protocol/chunk"
)

func TestDecodeFrame_ReadsWriteMessageFrame(t *testing.T) {
	var contentID core.ContentID
	contentID[0] = 0xAA
	contentID[31] = 0xBB

	msg := chunk.BuildRequestManifest(contentID)

	var buf bytes.Buffer
	if err := chunk.WriteMessage(&buf, msg); err != nil {
		t.Fatalf("WriteMessage failed: %v", err)
	}

	frame, err := chunk.DecodeFrame(&buf)
	if err != nil {
		t.Fatalf("DecodeFrame failed: %v", err)
	}

	if frame.Version != chunk.CurrentMessageVersion {
		t.Fatalf("expected version %d, got %d", chunk.CurrentMessageVersion, frame.Version)
	}
	if frame.MessageType != chunk.MsgRequestManifest {
		t.Fatalf("expected message type %d, got %d", chunk.MsgRequestManifest, frame.MessageType)
	}

	parsedID, err := chunk.ParseRequestManifest(frame.Payload)
	if err != nil {
		t.Fatalf("ParseRequestManifest failed: %v", err)
	}
	if parsedID != contentID {
		t.Fatalf("expected content ID %x, got %x", contentID, parsedID)
	}
}

func TestDecodeFrame_RejectsOversizedPayloadBeforeAllocation(t *testing.T) {
	var buf bytes.Buffer

	frameSize := uint32(chunk.MaxChunkRequestSize + 1 + 3)
	if err := binary.Write(&buf, binary.LittleEndian, frameSize); err != nil {
		t.Fatalf("failed to write frame size: %v", err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, chunk.CurrentMessageVersion); err != nil {
		t.Fatalf("failed to write version: %v", err)
	}
	if err := buf.WriteByte(byte(chunk.MsgRequestChunk)); err != nil {
		t.Fatalf("failed to write message type: %v", err)
	}

	_, err := chunk.DecodeFrame(&buf)
	if !errors.Is(err, chunk.ErrPayloadTooLarge) {
		t.Fatalf("expected ErrPayloadTooLarge, got %v", err)
	}
}
