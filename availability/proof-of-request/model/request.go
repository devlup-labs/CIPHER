package model

// Request represents a client's request for a specific file.
// RequestID uniquely identifies this request and will later be used
// for replay protection and for binding the PoW proof to this request.
type RequestStatus string

const (
	Pending  RequestStatus = "pending"
	Resolved RequestStatus = "resolved"
)

type Request struct {
	FileID    string
	ClientID  string
	Timestamp int64
	Status    RequestStatus
}
