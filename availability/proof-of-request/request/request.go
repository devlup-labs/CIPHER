package request

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"proof-of-request/model"
)

// CreateRequest creates a new request for a file.
// A cryptographically random RequestID is generated so that
// every request can be uniquely identified and replay protection
// can later distinguish separate requests.
func CreateRequest(fileID string, clientID string) (model.Request, error) {
	requestID, err := generateRequestID()
	if err != nil {
		return model.Request{}, err
	}

	return model.Request{
		RequestID: requestID,
		FileID:    fileID,
		ClientID:  clientID,
		Timestamp: time.Now().Unix(),
		Status:    model.Pending,
	}, nil
}

// generateRequestID generates a cryptographically random identifier.
// We keep this separate from CreateRequest because later the ID
// generation strategy can be changed without changing the request API.
func generateRequestID() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

var requests []model.Request

func GetOrCreateRequest(fileID string, clientID string) (model.Request, error) {
	for _, req := range requests {
		if req.FileID == fileID &&
			req.ClientID == clientID &&
			req.Status == model.Pending {
			return req, nil
		}
	}

	req, err := CreateRequest(fileID, clientID)
	if err != nil {
		return model.Request{}, err
	}

	requests = append(requests, req)

	return req, nil
}
