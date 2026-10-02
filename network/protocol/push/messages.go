package push

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"cipher/network/content/core"
)

const (
	CurrentPushVersion uint16 = 1

	MaxMessageSize  uint32        = 4 * 1024 * 1024 // 4 MB
	MaxChunkSize    uint32        = 2 * 1024 * 1024 // 2 MB
	MaxManifestSize uint32        = 2 * 1024 * 1024 // 2 MB

	ReadTimeout  = 15 * time.Second
	WriteTimeout = 15 * time.Second
	AckTimeout   = 30 * time.Second
)

type PushMessageType uint8

const (
	MsgPushManifest         PushMessageType = 0x10
	MsgPushManifestAck      PushMessageType = 0x11
	MsgPushChunk            PushMessageType = 0x12
	MsgPushChunkAck         PushMessageType = 0x13
	MsgPushBatchComplete    PushMessageType = 0x14
	MsgPushBatchCompleteAck PushMessageType = 0x15
	MsgPushError            PushMessageType = 0x16
)

const (
	PushStatusOK              byte = 0x00
	PushStatusUnauthorized    byte = 0x01
	PushStatusDiskFull        byte = 0x02
	PushStatusMalformed       byte = 0x03
	PushStatusIncomplete      byte = 0x04
	PushStatusHashMismatch    byte = 0x05
	PushStatusNotInAssignedSet byte = 0x06
	PushStatusIOError         byte = 0x07
)

type PushMessage struct {
	Version uint16
	Type    PushMessageType
	Payload []byte
}

func WritePushMessage(w io.Writer, msg *PushMessage) error {
	buf := new(bytes.Buffer)

	// Envelope: Version (2B), Type (1B)
	if err := binary.Write(buf, binary.LittleEndian, msg.Version); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.Type); err != nil {
		return err
	}
	buf.Write(msg.Payload)

	size := uint32(buf.Len())
	if size > MaxMessageSize {
		return fmt.Errorf("message size %d exceeds limit %d", size, MaxMessageSize)
	}

	// 4B size prefix
	if err := binary.Write(w, binary.LittleEndian, size); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

func ReadPushMessage(r io.Reader) (*PushMessage, error) {
	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, err
	}

	if size > MaxMessageSize {
		return nil, fmt.Errorf("message size %d exceeds maximum frame size %d", size, MaxMessageSize)
	}
	if size < 3 { // Must have at least Version (2) + Type (1)
		return nil, errors.New("message frame too short")
	}

	data := make([]byte, size)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)
	msg := &PushMessage{}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Version); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Type); err != nil {
		return nil, err
	}

	msg.Payload = make([]byte, buf.Len())
	if _, err := buf.Read(msg.Payload); err != nil && err != io.EOF {
		return nil, err
	}

	return msg, nil
}

// -- Payload Builders & Parsers --

func BuildPushManifest(contentID core.ContentID, assignedChunkIDs []core.ChunkID, manifestData []byte) *PushMessage {
	buf := new(bytes.Buffer)
	buf.Write(contentID[:])

	count := uint32(len(assignedChunkIDs))
	_ = binary.Write(buf, binary.LittleEndian, count)

	for _, cid := range assignedChunkIDs {
		buf.Write(cid[:])
	}
	buf.Write(manifestData)

	return &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushManifest,
		Payload: buf.Bytes(),
	}
}

func ParsePushManifest(payload []byte) (core.ContentID, []core.ChunkID, []byte, error) {
	var contentID core.ContentID
	if len(payload) < 36 { // 32 bytes ContentID + 4 bytes count
		return contentID, nil, nil, errors.New("invalid payload length for PUSH_MANIFEST")
	}

	copy(contentID[:], payload[:32])
	count := binary.LittleEndian.Uint32(payload[32:36])

	expectedOffset := 36 + int(count)*32
	if len(payload) < expectedOffset {
		return contentID, nil, nil, fmt.Errorf("payload length %d too short for %d assigned chunks", len(payload), count)
	}

	assignedChunkIDs := make([]core.ChunkID, count)
	for i := 0; i < int(count); i++ {
		offset := 36 + i*32
		copy(assignedChunkIDs[i][:], payload[offset:offset+32])
	}

	manifestData := payload[expectedOffset:]
	return contentID, assignedChunkIDs, manifestData, nil
}

func BuildPushManifestAck(contentID core.ContentID, status byte) *PushMessage {
	payload := append(contentID[:], status)
	return &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushManifestAck,
		Payload: payload,
	}
}

func ParsePushManifestAck(payload []byte) (core.ContentID, byte, error) {
	var id core.ContentID
	if len(payload) != 33 {
		return id, 0, fmt.Errorf("invalid payload length for PUSH_MANIFEST_ACK: %d", len(payload))
	}
	copy(id[:], payload[:32])
	return id, payload[32], nil
}

func BuildPushChunk(contentID core.ContentID, chunk *core.Chunk) (*PushMessage, error) {
	buf := new(bytes.Buffer)
	buf.Write(contentID[:])

	if err := binary.Write(buf, binary.LittleEndian, &chunk.Header); err != nil {
		return nil, err
	}
	buf.Write(chunk.Data)

	return &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushChunk,
		Payload: buf.Bytes(),
	}, nil
}

func ParsePushChunk(payload []byte) (core.ContentID, *core.Chunk, error) {
	var contentID core.ContentID
	var dummyHeader core.ChunkHeader
	headerSize := binary.Size(dummyHeader)

	if len(payload) < 32+headerSize {
		return contentID, nil, errors.New("invalid payload length for PUSH_CHUNK")
	}

	copy(contentID[:], payload[:32])

	buf := bytes.NewReader(payload[32:])
	chunk := &core.Chunk{}
	if err := binary.Read(buf, binary.LittleEndian, &chunk.Header); err != nil {
		return contentID, nil, err
	}

	chunk.Data = make([]byte, buf.Len())
	if _, err := buf.Read(chunk.Data); err != nil && err != io.EOF {
		return contentID, nil, err
	}

	return contentID, chunk, nil
}

func BuildPushChunkAck(chunkID core.ChunkID, status byte) *PushMessage {
	payload := append(chunkID[:], status)
	return &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushChunkAck,
		Payload: payload,
	}
}

func ParsePushChunkAck(payload []byte) (core.ChunkID, byte, error) {
	var id core.ChunkID
	if len(payload) != 33 {
		return id, 0, fmt.Errorf("invalid payload length for PUSH_CHUNK_ACK: %d", len(payload))
	}
	copy(id[:], payload[:32])
	return id, payload[32], nil
}

func BuildPushBatchComplete(contentID core.ContentID) *PushMessage {
	return &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushBatchComplete,
		Payload: contentID[:],
	}
}

func ParsePushBatchComplete(payload []byte) (core.ContentID, error) {
	var id core.ContentID
	if len(payload) != 32 {
		return id, fmt.Errorf("invalid payload length for PUSH_BATCH_COMPLETE: %d", len(payload))
	}
	copy(id[:], payload)
	return id, nil
}

func BuildPushBatchCompleteAck(contentID core.ContentID, status byte) *PushMessage {
	payload := append(contentID[:], status)
	return &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushBatchCompleteAck,
		Payload: payload,
	}
}

func ParsePushBatchCompleteAck(payload []byte) (core.ContentID, byte, error) {
	var id core.ContentID
	if len(payload) != 33 {
		return id, 0, fmt.Errorf("invalid payload length for PUSH_BATCH_COMPLETE_ACK: %d", len(payload))
	}
	copy(id[:], payload[:32])
	return id, payload[32], nil
}

func BuildPushError(code byte, msg string) *PushMessage {
	payload := append([]byte{code}, []byte(msg)...)
	return &PushMessage{
		Version: CurrentPushVersion,
		Type:    MsgPushError,
		Payload: payload,
	}
}

func ParsePushError(payload []byte) (byte, string, error) {
	if len(payload) < 1 {
		return 0, "", errors.New("invalid payload length for PUSH_ERROR")
	}
	return payload[0], string(payload[1:]), nil
}
