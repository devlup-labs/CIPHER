package protocol

import "github.com/libp2p/go-libp2p/core/protocol"

const (
	// FileTransferProtocolID is the legacy protocol for direct file transfer (deprecated in Milestone 8).
	FileTransferProtocolID protocol.ID = "/cipher/filetransfer/1.0.0"

	// ChunkTransportProtocolID is the content-addressed transport protocol (Milestone 8).
	ChunkTransportProtocolID protocol.ID = "/cipher/chunk/1.0.0"
)
