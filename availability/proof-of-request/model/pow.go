package model

// PoW represents a proof of work solution for a specific request.
type PoW struct {
    RequestID string
    Solution  string
    Timestamp int64
}