package chunk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"cipher/internal/content/core"
)

var (
	ErrNilMessage             = errors.New("message is nil")
	ErrInvalidMessageVersion  = errors.New("invalid message version")
	ErrInvalidMessageType     = errors.New("invalid message type")
	ErrInvalidPayloadLength   = errors.New("invalid payload length")
	ErrInvalidContentID       = errors.New("invalid content ID")
	ErrInvalidChunkID         = errors.New("invalid chunk ID")
	ErrInvalidChunkHeader     = errors.New("invalid chunk header")
	ErrInvalidCiphertext      = errors.New("invalid ciphertext")
	ErrInvalidAckStatus       = errors.New("invalid ack status")
	ErrInvalidErrorCode       = errors.New("invalid error code")
	ErrInvalidErrorMessage    = errors.New("invalid error message")
	ErrContentMismatch        = errors.New("content ID does not match request")
	ErrChunkMismatch          = errors.New("chunk ID does not match request")
	ErrChunkIndexMismatch     = errors.New("chunk index does not match request")
	ErrInvalidChunkIndex      = errors.New("invalid chunk index")
	ErrInvalidManifestPayload = errors.New("invalid manifest payload")
)

type ContentInfo struct {
	ContentID   core.ContentID
	TotalChunks uint64
}

// ValidateMessage performs cheap structural validation on the current
// length-prefixed Message envelope and its payload.
func ValidateMessage(msg *Message) error {
	if msg == nil {
		return fmt.Errorf("message: %w", ErrNilMessage)
	}

	if msg.Version != CurrentMessageVersion {
		return fmt.Errorf(
			"message: %w: received=%d supported=%d",
			ErrInvalidMessageVersion,
			msg.Version,
			CurrentMessageVersion,
		)
	}

	return ValidateMessagePayload(msg.Type, msg.Payload)
}

func ValidateMessagePayload(messageType MessageType, payload []byte) error {
	maxPayloadSize := MaxPayloadSizeForMessage(messageType)
	if maxPayloadSize <= 0 {
		return fmt.Errorf("message: %w: %d", ErrInvalidMessageType, messageType)
	}

	if len(payload) > maxPayloadSize {
		return fmt.Errorf(
			"message: %w: type=%d received=%d maximum=%d",
			ErrInvalidPayloadLength,
			messageType,
			len(payload),
			maxPayloadSize,
		)
	}

	switch messageType {
	case MsgRequestManifest:
		return ValidateRequestManifestPayload(payload)
	case MsgManifest:
		return ValidateManifestPayload(payload)
	case MsgRequestChunk:
		return ValidateRequestChunkPayload(payload)
	case MsgChunk:
		return ValidateChunkPayload(payload)
	case MsgAck:
		return ValidateAckPayload(payload)
	case MsgError:
		return ValidateErrorPayload(payload)
	default:
		return fmt.Errorf("message: %w: %d", ErrInvalidMessageType, messageType)
	}
}

func ValidateRequestManifestPayload(payload []byte) error {
	if len(payload) != ContentIDSize {
		return fmt.Errorf(
			"request manifest: %w: received=%d expected=%d",
			ErrInvalidContentID,
			len(payload),
			ContentIDSize,
		)
	}
	return nil
}

func ValidateManifestPayload(payload []byte) error {
	if len(payload) < ContentIDSize {
		return fmt.Errorf(
			"manifest: %w: received=%d minimum=%d",
			ErrInvalidManifestPayload,
			len(payload),
			ContentIDSize,
		)
	}
	return nil
}

func ValidateRequestChunkPayload(payload []byte) error {
	if len(payload) != HashSize {
		return fmt.Errorf(
			"request chunk: %w: received=%d expected=%d",
			ErrInvalidChunkID,
			len(payload),
			HashSize,
		)
	}
	return nil
}

func ValidateChunkPayload(payload []byte) error {
	headerSize := binary.Size(core.ChunkHeader{})
	if headerSize <= 0 {
		return fmt.Errorf("chunk: %w", ErrInvalidChunkHeader)
	}

	if len(payload) < headerSize {
		return fmt.Errorf(
			"chunk: %w: received=%d minimum=%d",
			ErrInvalidPayloadLength,
			len(payload),
			headerSize,
		)
	}

	header, ciphertext, err := splitChunkPayload(payload)
	if err != nil {
		return err
	}

	if header.Version != CurrentMessageVersion {
		return fmt.Errorf(
			"chunk: %w: received=%d supported=%d",
			ErrInvalidMessageVersion,
			header.Version,
			CurrentMessageVersion,
		)
	}

	if uint64(header.Index) > MaxSupportedChunkIndex {
		return fmt.Errorf(
			"chunk: %w: index=%d maximum=%d",
			ErrInvalidChunkIndex,
			header.Index,
			MaxSupportedChunkIndex,
		)
	}

	if header.CipherSize != uint32(len(ciphertext)) {
		return fmt.Errorf(
			"chunk: %w: header=%d payload=%d",
			ErrInvalidCiphertext,
			header.CipherSize,
			len(ciphertext),
		)
	}

	if len(ciphertext) == 0 {
		return fmt.Errorf("chunk: %w: empty ciphertext", ErrInvalidCiphertext)
	}

	return nil
}

func ValidateAckPayload(payload []byte) error {
	if len(payload) != HashSize+1 {
		return fmt.Errorf(
			"ack: %w: received=%d expected=%d",
			ErrInvalidPayloadLength,
			len(payload),
			HashSize+1,
		)
	}

	if payload[HashSize] != 0 {
		return fmt.Errorf("ack: %w: %d", ErrInvalidAckStatus, payload[HashSize])
	}

	return nil
}

func ValidateErrorPayload(payload []byte) error {
	if len(payload) < 1 {
		return fmt.Errorf(
			"protocol error: %w: received=%d minimum=1",
			ErrInvalidPayloadLength,
			len(payload),
		)
	}

	code := ErrorCode(payload[0])
	if !isKnownErrorCode(code) {
		return fmt.Errorf("protocol error: %w: %d", ErrInvalidErrorCode, code)
	}

	messageSize := len(payload) - 1
	if messageSize > MaxErrorMessageSize {
		return fmt.Errorf(
			"protocol error: %w: received=%d maximum=%d",
			ErrInvalidErrorMessage,
			messageSize,
			MaxErrorMessageSize,
		)
	}

	return nil
}

func ValidateManifestForRequest(requested core.ContentID, payload []byte) error {
	if err := ValidateManifestPayload(payload); err != nil {
		return err
	}

	received, _, err := ParseManifest(payload)
	if err != nil {
		return err
	}

	if received != requested {
		return fmt.Errorf("manifest: %w", ErrContentMismatch)
	}

	return nil
}

func ValidateChunkForRequest(requested core.ChunkID, payload []byte) error {
	if err := ValidateChunkPayload(payload); err != nil {
		return err
	}

	header, _, err := splitChunkPayload(payload)
	if err != nil {
		return err
	}

	if header.ID != requested {
		return fmt.Errorf("chunk: %w", ErrChunkMismatch)
	}

	return nil
}

func ValidateChunkRequestForContent(chunkIndex uint64, info ContentInfo) error {
	if info.TotalChunks == 0 || info.TotalChunks > MaxSupportedChunkCount {
		return fmt.Errorf("content information has invalid total chunk count: %d", info.TotalChunks)
	}

	if chunkIndex >= info.TotalChunks {
		return fmt.Errorf(
			"chunk request: %w: index=%d total_chunks=%d",
			ErrInvalidChunkIndex,
			chunkIndex,
			info.TotalChunks,
		)
	}

	return nil
}

func splitChunkPayload(payload []byte) (core.ChunkHeader, []byte, error) {
	var header core.ChunkHeader
	headerSize := binary.Size(header)
	if len(payload) < headerSize {
		return header, nil, fmt.Errorf(
			"chunk: %w: received=%d minimum=%d",
			ErrInvalidPayloadLength,
			len(payload),
			headerSize,
		)
	}

	reader := bytes.NewReader(payload[:headerSize])
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return header, nil, fmt.Errorf("chunk: %w: %v", ErrInvalidChunkHeader, err)
	}

	return header, payload[headerSize:], nil
}

func isKnownErrorCode(code ErrorCode) bool {
	switch code {
	case ErrContentNotFound,
		ErrChunkNotFound,
		ErrInvalidManifest,
		ErrPermissionDenied,
		ErrInternal,
		ErrIntegrityMismatch,
		ErrBadRequest,
		ErrUnsupportedMessage:
		return true
	default:
		return false
	}
}
