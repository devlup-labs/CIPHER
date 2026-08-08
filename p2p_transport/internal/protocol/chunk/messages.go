package chunk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"cipher/internal/content/core"
)

const (
	CurrentMessageVersion uint16 = 1
)

type MessageType uint8

const (
	MsgRequestManifest MessageType = 0x01
	MsgManifest        MessageType = 0x02
	MsgRequestChunk    MessageType = 0x03
	MsgChunk           MessageType = 0x04
	MsgAck             MessageType = 0x05
	MsgError           MessageType = 0x06
)

type ErrorCode uint8

const (
	ErrContentNotFound  ErrorCode = 0x01
	ErrChunkNotFound    ErrorCode = 0x02
	ErrInvalidManifest  ErrorCode = 0x03
	ErrPermissionDenied ErrorCode = 0x04
	ErrInternal         ErrorCode = 0x05
	ErrIntegrityMismatch ErrorCode = 0x06
	ErrBadRequest       ErrorCode = 0x07
	ErrUnsupportedMessage ErrorCode = 0x08
)

// Message is the symmetric envelope for all protocol communications.
type Message struct {
	Version uint16
	Type    MessageType
	Payload []byte
}

func WriteMessage(w io.Writer, msg *Message) error {
	buf := new(bytes.Buffer)
	
	// Envelope: Version (2), Type (1)
	if err := binary.Write(buf, binary.LittleEndian, msg.Version); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.Type); err != nil {
		return err
	}
	// Payload
	buf.Write(msg.Payload)

	// Frame: Size prefix (4 bytes)
	size := uint32(buf.Len())
	if err := binary.Write(w, binary.LittleEndian, size); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

func ReadMessage(r io.Reader) (*Message, error) {
	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, err
	}

	if size > 2*1024*1024 { // 2MB max frame size
		return nil, errors.New("message exceeds maximum frame size")
	}

	data := make([]byte, size)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Version); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Type); err != nil {
		return nil, err
	}

	msg.Payload = make([]byte, buf.Len())
	buf.Read(msg.Payload)

	return msg, nil
}

// -- Payload Builders & Parsers --

func BuildRequestManifest(id core.ContentID) *Message {
	return &Message{
		Version: CurrentMessageVersion,
		Type:    MsgRequestManifest,
		Payload: id[:],
	}
}

func ParseRequestManifest(payload []byte) (core.ContentID, error) {
	var id core.ContentID
	if len(payload) != 32 {
		return id, fmt.Errorf("invalid payload length for REQUEST_MANIFEST: %d", len(payload))
	}
	copy(id[:], payload)
	return id, nil
}

func BuildManifest(id core.ContentID, data []byte) *Message {
	payload := append(id[:], data...)
	return &Message{
		Version: CurrentMessageVersion,
		Type:    MsgManifest,
		Payload: payload,
	}
}

func ParseManifest(payload []byte) (core.ContentID, []byte, error) {
	var id core.ContentID
	if len(payload) < 32 {
		return id, nil, fmt.Errorf("invalid payload length for MANIFEST: %d", len(payload))
	}
	copy(id[:], payload[:32])
	return id, payload[32:], nil
}

func BuildRequestChunk(id core.ChunkID) *Message {
	return &Message{
		Version: CurrentMessageVersion,
		Type:    MsgRequestChunk,
		Payload: id[:],
	}
}

func ParseRequestChunk(payload []byte) (core.ChunkID, error) {
	var id core.ChunkID
	if len(payload) != 32 {
		return id, fmt.Errorf("invalid payload length for REQUEST_CHUNK: %d", len(payload))
	}
	copy(id[:], payload)
	return id, nil
}

func BuildChunk(chunk *core.Chunk) (*Message, error) {
	buf := new(bytes.Buffer)
	// Write Header
	if err := binary.Write(buf, binary.LittleEndian, &chunk.Header); err != nil {
		return nil, err
	}
	// Write Ciphertext
	buf.Write(chunk.Data)
	return &Message{
		Version: CurrentMessageVersion,
		Type:    MsgChunk,
		Payload: buf.Bytes(),
	}, nil
}

func ParseChunk(payload []byte) (*core.Chunk, error) {
	chunk := &core.Chunk{}
	buf := bytes.NewReader(payload)
	if err := binary.Read(buf, binary.LittleEndian, &chunk.Header); err != nil {
		return nil, err
	}
	chunk.Data = make([]byte, buf.Len())
	buf.Read(chunk.Data)
	return chunk, nil
}

func BuildAck(id core.ChunkID, status uint8) *Message {
	payload := append(id[:], status)
	return &Message{
		Version: CurrentMessageVersion,
		Type:    MsgAck,
		Payload: payload,
	}
}

func ParseAck(payload []byte) (core.ChunkID, uint8, error) {
	var id core.ChunkID
	if len(payload) != 33 {
		return id, 0, fmt.Errorf("invalid payload length for ACK: %d", len(payload))
	}
	copy(id[:], payload[:32])
	return id, payload[32], nil
}

func BuildError(code ErrorCode, msg string) *Message {
	payload := append([]byte{byte(code)}, []byte(msg)...)
	return &Message{
		Version: CurrentMessageVersion,
		Type:    MsgError,
		Payload: payload,
	}
}

func ParseError(payload []byte) (ErrorCode, string, error) {
	if len(payload) < 1 {
		return 0, "", errors.New("invalid payload length for ERROR")
	}
	return ErrorCode(payload[0]), string(payload[1:]), nil
}
