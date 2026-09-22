package request

import (
	"time"

	"proof-of-request/model"
)

// CreateRequest creates a new request for a file.
func CreateRequest(fileID string, clientID string) (model.Request, error) {

	return model.Request{
		FileID:    fileID,
		ClientID:  clientID,
		Timestamp: time.Now().Unix(),
		Status:    model.Pending,
	}, nil
}

var requests []model.Request

// GetOrCreateRequest returns an existing pending request for the given fileID and clientID,
// or creates a new one if none exists.
func GetOrCreateRequest(fileID string, clientID string) (model.Request, bool, error) {
	for _, req := range requests {
		if req.FileID == fileID &&
			req.ClientID == clientID &&
			req.Status == model.Pending {
			return req, false, nil
		}
	}

	req, err := CreateRequest(fileID, clientID)
	if err != nil {
		return model.Request{}, true, err
	}

	requests = append(requests, req)

	return req, true, nil
}

// encapsulate resolution logic
func ResolveRequest(fileID string, clientID string) bool {
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
