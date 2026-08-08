package chunk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// Frame header:
	// 4 bytes - envelope length, little-endian
	// 2 bytes - protocol version, little-endian
	// 1 byte  - message type
	frameHeaderSize = 7

	frameLengthSize   = 4
	messageHeaderSize = 3
)

var (
	ErrInvalidProtocolVersion = errors.New("invalid protocol version")
	ErrUnknownMessageType     = errors.New("unknown message type")
	ErrInvalidPayloadSize     = errors.New("invalid payload size")
	ErrPayloadTooLarge        = errors.New("payload exceeds protocol limit")
	ErrTruncatedFrame         = errors.New("truncated protocol frame")
)

// Frame represents one safely decoded chunk-protocol frame.
//
// Payload is only allocated after its declared size has been checked against
// the limit for the corresponding message type.
type Frame struct {
	Version     uint16
	MessageType MessageType
	Payload     []byte
}

// DecodeFrame reads one complete frame from r.
//
// It first reads the small fixed-size header, validates the protocol version,
// message type and declared payload size, and only then allocates memory for
// the payload.
func DecodeFrame(r io.Reader) (*Frame, error) {
	if r == nil {
		return nil, errors.New("chunk decoder: nil reader")
	}

	header := make([]byte, frameHeaderSize)

	if _, err := io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("%w: failed to read frame header: %v", ErrTruncatedFrame, err)
	}

	frameSize := binary.LittleEndian.Uint32(header[:frameLengthSize])
	if frameSize < messageHeaderSize {
		return nil, fmt.Errorf(
			"%w: frame size %d is smaller than message header",
			ErrInvalidPayloadSize,
			frameSize,
		)
	}

	version := binary.LittleEndian.Uint16(header[frameLengthSize : frameLengthSize+2])
	if version != CurrentMessageVersion {
		return nil, fmt.Errorf(
			"%w: received=%d supported=%d",
			ErrInvalidProtocolVersion,
			version,
			CurrentMessageVersion,
		)
	}

	messageType := MessageType(header[frameLengthSize+2])

	maxPayloadSize := MaxPayloadSizeForMessage(messageType)
	if maxPayloadSize <= 0 {
		return nil, fmt.Errorf(
			"%w: %d",
			ErrUnknownMessageType,
			messageType,
		)
	}

	declaredSize := frameSize - messageHeaderSize

	if declaredSize == 0 {
		return nil, fmt.Errorf(
			"%w: message type %d declared an empty payload",
			ErrInvalidPayloadSize,
			messageType,
		)
	}

	if uint64(declaredSize) > uint64(maxPayloadSize) {
		return nil, fmt.Errorf(
			"%w: message type=%d declared=%d maximum=%d",
			ErrPayloadTooLarge,
			messageType,
			declaredSize,
			maxPayloadSize,
		)
	}

	// Allocation happens only after all size checks pass.
	payload := make([]byte, int(declaredSize))

	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf(
			"%w: expected=%d bytes: %v",
			ErrTruncatedFrame,
			declaredSize,
			err,
		)
	}

	return &Frame{
		Version:     version,
		MessageType: messageType,
		Payload:     payload,
	}, nil
}

// DecodeExpectedFrame decodes a frame and verifies that it has the message
// type expected for the current protocol step.
func DecodeExpectedFrame(
	r io.Reader,
	expected MessageType,
) (*Frame, error) {
	frame, err := DecodeFrame(r)
	if err != nil {
		return nil, err
	}

	if frame.MessageType != expected {
		return nil, fmt.Errorf(
			"unexpected message type: received=%d expected=%d",
			frame.MessageType,
			expected,
		)
	}

	return frame, nil
}

// DecodeFrameHeader reads and validates only a frame header.
//
// This can be useful when the caller needs to reserve memory through the
// libp2p Resource Manager before reading the payload.
func DecodeFrameHeader(r io.Reader) (
	version uint16,
	messageType MessageType,
	payloadSize uint32,
	err error,
) {
	if r == nil {
		err = errors.New("chunk decoder: nil reader")
		return
	}

	header := make([]byte, frameHeaderSize)

	if _, readErr := io.ReadFull(r, header); readErr != nil {
		err = fmt.Errorf(
			"%w: failed to read frame header: %v",
			ErrTruncatedFrame,
			readErr,
		)
		return
	}

	frameSize := binary.LittleEndian.Uint32(header[:frameLengthSize])
	if frameSize < messageHeaderSize {
		err = fmt.Errorf(
			"%w: frame size %d is smaller than message header",
			ErrInvalidPayloadSize,
			frameSize,
		)
		return
	}

	version = binary.LittleEndian.Uint16(header[frameLengthSize : frameLengthSize+2])
	if version != CurrentMessageVersion {
		err = fmt.Errorf(
			"%w: received=%d supported=%d",
			ErrInvalidProtocolVersion,
			version,
			CurrentMessageVersion,
		)
		return
	}

	messageType = MessageType(header[frameLengthSize+2])

	maxPayloadSize := MaxPayloadSizeForMessage(messageType)
	if maxPayloadSize <= 0 {
		err = fmt.Errorf(
			"%w: %d",
			ErrUnknownMessageType,
			messageType,
		)
		return
	}

	payloadSize = frameSize - messageHeaderSize

	if payloadSize == 0 {
		err = fmt.Errorf(
			"%w: message type %d declared an empty payload",
			ErrInvalidPayloadSize,
			messageType,
		)
		return
	}

	if uint64(payloadSize) > uint64(maxPayloadSize) {
		err = fmt.Errorf(
			"%w: message type=%d declared=%d maximum=%d",
			ErrPayloadTooLarge,
			messageType,
			payloadSize,
			maxPayloadSize,
		)
		return
	}

	return
}

// ReadFramePayload reads an exact payload after DecodeFrameHeader has already
// checked its declared size.
//
// The size is checked again so this function remains safe if used separately.
func ReadFramePayload(
	r io.Reader,
	messageType MessageType,
	payloadSize uint32,
) ([]byte, error) {
	if r == nil {
		return nil, errors.New("chunk decoder: nil reader")
	}

	maxPayloadSize := MaxPayloadSizeForMessage(messageType)
	if maxPayloadSize <= 0 {
		return nil, fmt.Errorf(
			"%w: %d",
			ErrUnknownMessageType,
			messageType,
		)
	}

	if payloadSize == 0 {
		return nil, ErrInvalidPayloadSize
	}

	if uint64(payloadSize) > uint64(maxPayloadSize) {
		return nil, fmt.Errorf(
			"%w: declared=%d maximum=%d",
			ErrPayloadTooLarge,
			payloadSize,
			maxPayloadSize,
		)
	}

	payload := make([]byte, int(payloadSize))

	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf(
			"%w: expected=%d bytes: %v",
			ErrTruncatedFrame,
			payloadSize,
			err,
		)
	}

	return payload, nil
}
