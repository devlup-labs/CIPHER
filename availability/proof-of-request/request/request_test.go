package request

import (
	"proof-of-request/model"
	"testing"
)

func TestCreateRequest(t *testing.T) {
	fileID := "file-123"
	clientID := "client-456"

	req, err := CreateRequest(fileID, clientID)

	if err != nil {
		t.Fatalf("CreateRequest() returned an error: %v", err)
	}

	if req.FileID != fileID {
		t.Errorf("expected FileID %q, got %q", fileID, req.FileID)
	}

	if req.ClientID != clientID {
		t.Errorf("expected ClientID %q, got %q", clientID, req.ClientID)
	}

	if req.RequestID == "" {
		t.Error("expected RequestID to be generated")
	}

	if req.Timestamp == 0 {
		t.Error("expected Timestamp to be generated")
	}
}

// input preservation
func TestCreateRequestPreservesInput(t *testing.T) {
	req, err := CreateRequest("file-A", "client-B")
	if err != nil {
		t.Fatalf("CreateRequest() failed: %v", err)
	}

	if req.FileID != "file-A" {
		t.Errorf("expected FileID file-A, got %s", req.FileID)
	}

	if req.ClientID != "client-B" {
		t.Errorf("expected ClientID client-B, got %s", req.ClientID)
	}
}

// test GetOrCreateRequest returns existing request if it exists
func TestGetOrCreateRequestReusesPendingRequest(t *testing.T) {
	req1, err := GetOrCreateRequest("file-123", "client-456")
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	req2, err := GetOrCreateRequest("file-123", "client-456")
	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if req1.RequestID != req2.RequestID {
		t.Errorf(
			"expected existing request to be reused, got different RequestIDs: %s and %s",
			req1.RequestID,
			req2.RequestID,
		)
	}
}

func TestGetOrCreateRequestCreatesNewRequestAfterResolution(t *testing.T) {
	req1, err := GetOrCreateRequest("file-789", "client-101")
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	// Mark the existing request as resolved.
	for i := range requests {
		if requests[i].RequestID == req1.RequestID {
			requests[i].Status = model.Resolved
		}
	}

	req2, err := GetOrCreateRequest("file-789", "client-101")
	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if req1.RequestID == req2.RequestID {
		t.Error("expected a new RequestID after the previous request was resolved")
	}

	if req2.Status != model.Pending {
		t.Errorf("expected new request to be pending, got %s", req2.Status)
	}
}
