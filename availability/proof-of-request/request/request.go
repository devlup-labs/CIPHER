package request

import (
	"errors"
	"time"

	"proof-of-request/model"
)

var requests []model.Request

func CreateRequest(fileID string, clientID string) (model.Request, error) {
	// Create a new pending request with the current timestamp.
	return model.Request{
		FileID:    fileID,
		ClientID:  clientID,
		Timestamp: time.Now().UnixNano(),// Use UnixNano for higher precision timestamps
		Status:    model.Pending,
	}, nil
}

func GetOrCreateRequest(fileID string, clientID string, lifetime time.Duration) (model.Request, bool, error) {
	// A request lifetime must be positive.
	if lifetime <= 0 {
		return model.Request{}, false, errors.New("request lifetime must be positive")
	}

	// Check whether this client already has an active request for this file.
	for i := range requests {
		if requests[i].FileID != fileID ||
			requests[i].ClientID != clientID ||
			requests[i].Status != model.Pending {
			continue
		}

		// If the pending request is still within its lifetime, reuse it.
		requestAge := time.Since(time.Unix(0, requests[i].Timestamp)) // Use Unix(0, timestamp) to convert nanoseconds to time.Time

		if requestAge <= lifetime {
			return requests[i], false, nil
		}

		// The request is too old, so automatically expire it.
		requests[i].Status = model.Expired
	}

	// No active request exists, so create a new one.
	req, err := CreateRequest(fileID, clientID)
	if err != nil {
		return model.Request{}, false, err
	}

	requests = append(requests, req)

	return req, true, nil
}

func ResolveRequest(fileID string, clientID string) bool {
	// Find the active request and mark it as resolved.
	for i := range requests {
		if requests[i].FileID == fileID &&
			requests[i].ClientID == clientID &&
			requests[i].Status == model.Pending {
			requests[i].Status = model.Resolved
			return true
		}
	}

	return false
}

func AbortRequest(fileID string, clientID string) bool {
	// Find the active request and mark it as explicitly aborted by the client.
	for i := range requests {
		if requests[i].FileID == fileID &&
			requests[i].ClientID == clientID &&
			requests[i].Status == model.Pending {
			requests[i].Status = model.Aborted
			return true
		}
	}

	return false
}