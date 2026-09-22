package request

import (
	"proof-of-request/model"
	"testing"
	"time"
)

func resetRequests() {
	requests = nil
}
func TestCreateRequest(t *testing.T) {
	resetRequests()
	req, err := CreateRequest("file-123", "client-456")

	if err != nil {
		t.Fatalf("CreateRequest() failed: %v", err)
	}

	if req.FileID != "file-123" ||
		req.ClientID != "client-456" ||
		req.Status != model.Pending ||
		req.Timestamp == 0 {
		t.Error("created request has incorrect fields")
	}
}

func TestGetOrCreateRequestLifecycle(t *testing.T) {
	resetRequests()
	lifetime := time.Hour

	// First call should create a new request.
	req1, isNew, err := GetOrCreateRequest("file-123", "client-456", lifetime)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	if !isNew {
		t.Error("expected first request to be new")
	}

	// Second call should reuse the same pending request.
	req2, isNew, err := GetOrCreateRequest("file-123", "client-456", lifetime)
	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if isNew || req1 != req2 {
		t.Error("expected existing pending request to be reused")
	}

	// Resolve the request.
	if !ResolveRequest(req1.FileID, req1.ClientID) {
		t.Fatal("expected pending request to be resolved")
	}

	// After resolution, a new request should be created.
	req3, isNew, err := GetOrCreateRequest("file-123", "client-456", lifetime)
	if err != nil {
		t.Fatalf("third request failed: %v", err)
	}

	if !isNew || req3.Status != model.Pending {
		t.Error("expected a new pending request after resolution")
	}
}

func TestGetOrCreateRequestExpiresOldRequest(t *testing.T) {
	resetRequests()
	// Create a request with a very short lifetime.
	req1, _, err := GetOrCreateRequest("file-expire", "client-expire", time.Nanosecond)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	// Wait until the request is outside its lifetime.
	time.Sleep(time.Millisecond)

	req2, isNew, err := GetOrCreateRequest("file-expire", "client-expire", time.Nanosecond)
	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if !isNew {
		t.Error("expected expired request to produce a new request")
	}

	if req2.Status != model.Pending {
		t.Errorf("expected new request to be pending, got %s", req2.Status)
	}

	// The old request should no longer be pending.
	if requests[0].Status != model.Expired {
		t.Errorf("expected old request to be expired, got %s", requests[0].Status)
	}

	if req1 == req2 {
		t.Error("expected a different request after expiration")
	}
}

func TestAbortRequest(t *testing.T) {
	resetRequests()
	req, _, err := GetOrCreateRequest("file-abort", "client-abort", time.Hour)
	if err != nil {
		t.Fatalf("request creation failed: %v", err)
	}

	if !AbortRequest(req.FileID, req.ClientID) {
		t.Fatal("expected pending request to be aborted")
	}

	// An aborted request should allow a new request to be created.
	newReq, isNew, err := GetOrCreateRequest(req.FileID, req.ClientID, time.Hour)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}

	if !isNew || newReq.Status != model.Pending {
		t.Error("expected new pending request after abort")
	}
}

func TestResolveAndAbortReturnFalseWithoutPendingRequest(t *testing.T) {
	resetRequests()
	if ResolveRequest("missing-file", "missing-client") {
		t.Error("expected ResolveRequest to return false")
	}

	if AbortRequest("missing-file", "missing-client") {
		t.Error("expected AbortRequest to return false")
	}
}

func TestGetOrCreateRequestRejectsInvalidLifetime(t *testing.T) {
	resetRequests()
	_, _, err := GetOrCreateRequest("file-123", "client-456", 0)

	if err == nil {
		t.Error("expected error for non-positive lifetime")
	}
}
