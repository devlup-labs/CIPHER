package request

import (
	"proof-of-request/model"
	"testing"
	"time"
)

// resetRequests clears the active request store between tests.
//
// We also stop all active timers because a timer from one test must
// not fire later and modify the state of another test.
func resetRequests() {
	mu.Lock()
	defer mu.Unlock()

	for _, record := range requests {
		if record.timer != nil {
			record.timer.Stop()
		}
	}

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

	// The first call should create a new Pending request.
	req1, isNew, err := GetOrCreateRequest(
		"file-123",
		"client-456",
		lifetime,
	)

	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	if !isNew {
		t.Error("expected first request to be new")
	}

	// A second request from the same client for the same file
	// should reuse the existing Pending request.
	req2, isNew, err := GetOrCreateRequest(
		"file-123",
		"client-456",
		lifetime,
	)

	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if isNew || req1 != req2 {
		t.Error("expected existing pending request to be reused")
	}

	// Resolve the active request.
	//
	// Resolution should remove it from the active request store.
	if !ResolveRequest(req1.FileID, req1.ClientID) {
		t.Fatal("expected pending request to be resolved")
	}

	// Because the old request is no longer Pending, the next request
	// should create a completely new request.
	req3, isNew, err := GetOrCreateRequest(
		"file-123",
		"client-456",
		lifetime,
	)

	if err != nil {
		t.Fatalf("third request failed: %v", err)
	}

	if !isNew || req3.Status != model.Pending {
		t.Error("expected a new pending request after resolution")
	}

	// Only the new Pending request should remain in the active store.
	pending := GetPendingRequests()

	if len(pending) != 1 {
		t.Fatalf("expected 1 pending request, got %d", len(pending))
	}

	if pending[0] != req3 {
		t.Error("expected the new request to be the only pending request")
	}
}

func TestGetOrCreateRequestExpiresOldRequest(t *testing.T) {
	resetRequests()

	// Use a short lifetime so that the timer automatically expires
	// the request.
	req1, isNew, err := GetOrCreateRequest(
		"file-expire",
		"client-expire",
		10*time.Millisecond,
	)

	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	if !isNew {
		t.Fatal("expected first request to be new")
	}

	// Wait long enough for the expiration timer to fire.
	time.Sleep(30 * time.Millisecond)

	// The expired request should have been removed from the active store.
	pending := GetPendingRequests()

	if len(pending) != 0 {
		t.Fatalf(
			"expected no pending requests after expiration, got %d",
			len(pending),
		)
	}

	// Since the previous request is no longer active, a new request
	// should be created for the same client/file pair.
	req2, isNew, err := GetOrCreateRequest(
		"file-expire",
		"client-expire",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if !isNew {
		t.Error("expected expired request to allow a new request")
	}

	if req2.Status != model.Pending {
		t.Errorf(
			"expected new request to be pending, got %s",
			req2.Status,
		)
	}

	if req1 == req2 {
		t.Error("expected a different request after expiration")
	}
}

func TestAbortRequest(t *testing.T) {
	resetRequests()

	req, _, err := GetOrCreateRequest(
		"file-abort",
		"client-abort",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("request creation failed: %v", err)
	}

	// Abort the active request.
	if !AbortRequest(req.FileID, req.ClientID) {
		t.Fatal("expected pending request to be aborted")
	}

	// The aborted request should no longer contribute to demand.
	pending := GetPendingRequests()

	if len(pending) != 0 {
		t.Errorf(
			"expected no pending requests after abort, got %d",
			len(pending),
		)
	}

	// Because the old request is no longer Pending, a new request
	// should be allowed.
	newReq, isNew, err := GetOrCreateRequest(
		req.FileID,
		req.ClientID,
		time.Hour,
	)

	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}

	if !isNew || newReq.Status != model.Pending {
		t.Error("expected new pending request after abort")
	}
}

func TestResolveAndAbortReturnFalseWithoutPendingRequest(t *testing.T) {
	resetRequests()

	// There is no active request, so both operations should fail.
	if ResolveRequest("missing-file", "missing-client") {
		t.Error("expected ResolveRequest to return false")
	}

	if AbortRequest("missing-file", "missing-client") {
		t.Error("expected AbortRequest to return false")
	}
}

// Test that GetOrCreateRequest rejects non-positive lifetimes.
func TestGetOrCreateRequestRejectsInvalidLifetime(t *testing.T) {
	resetRequests()

	_, _, err := GetOrCreateRequest(
		"file-123",
		"client-456",
		0,
	)

	if err == nil {
		t.Error("expected error for non-positive lifetime")
	}
}
