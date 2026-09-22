package demand

import (
	"testing"
	"time"

	"proof-of-request/request"
)

// resetDemandState resets the request module between tests.
//
// Demand no longer owns request records, so there is no records slice
// to reset here. The active request store belongs to the request package.
//
// This helper is kept here only if we later add demand-specific state.
// Currently it intentionally does nothing.
func resetDemandState() {
	// No demand state to reset.
}

func TestCalculateDemand(t *testing.T) {
	resetDemandState()

	// Create two active requests for file-A.
	_, _, err := request.GetOrCreateRequest(
		"file-A",
		"client-1",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create first request: %v", err)
	}

	_, _, err = request.GetOrCreateRequest(
		"file-A",
		"client-2",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create second request: %v", err)
	}

	// Create one active request for file-B.
	_, _, err = request.GetOrCreateRequest(
		"file-B",
		"client-3",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create third request: %v", err)
	}

	// file-A currently has two Pending requests.
	demand, err := CalculateDemand("file-A")

	if err != nil {
		t.Fatalf("CalculateDemand() failed: %v", err)
	}

	if demand != 2 {
		t.Errorf("expected demand 2, got %d", demand)
	}

	// file-B currently has one Pending request.
	demand, err = CalculateDemand("file-B")

	if err != nil {
		t.Fatalf("CalculateDemand() failed: %v", err)
	}

	if demand != 1 {
		t.Errorf("expected demand 1, got %d", demand)
	}
}

func TestCalculateDemandIgnoresResolvedRequests(t *testing.T) {
	resetDemandState()

	// Create two requests for the same file.
	_, _, err := request.GetOrCreateRequest(
		"file-A",
		"client-1",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create first request: %v", err)
	}

	req2, _, err := request.GetOrCreateRequest(
		"file-A",
		"client-2",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create second request: %v", err)
	}

	// Resolve one request.
	//
	// The resolved request should be removed from the active store
	// and therefore must no longer contribute to demand.
	if !request.ResolveRequest(req2.FileID, req2.ClientID) {
		t.Fatal("expected request to be resolved")
	}

	demand, err := CalculateDemand("file-A")

	if err != nil {
		t.Fatalf("CalculateDemand() failed: %v", err)
	}

	// Only client-1 still has an active Pending request.
	if demand != 1 {
		t.Errorf("expected demand 1, got %d", demand)
	}
}

func TestCalculateDemandIgnoresAbortedRequests(t *testing.T) {
	resetDemandState()

	// Create two requests for the same file.
	req1, _, err := request.GetOrCreateRequest(
		"file-A",
		"client-1",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create first request: %v", err)
	}

	_, _, err = request.GetOrCreateRequest(
		"file-A",
		"client-2",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create second request: %v", err)
	}

	// Abort one request.
	if !request.AbortRequest(req1.FileID, req1.ClientID) {
		t.Fatal("expected request to be aborted")
	}

	demand, err := CalculateDemand("file-A")

	if err != nil {
		t.Fatalf("CalculateDemand() failed: %v", err)
	}

	// Only client-2 remains Pending.
	if demand != 1 {
		t.Errorf("expected demand 1, got %d", demand)
	}
}

func TestCalculateDemandIgnoresExpiredRequests(t *testing.T) {

	// Use a unique FileID so requests created by other tests cannot
	// affect this test.
	fileID := "file-expiration-test"

	// Create a request with a short lifetime.
	//
	// The request module owns expiration. Demand simply observes
	// whether the request is still active.
	_, _, err := request.GetOrCreateRequest(
		fileID,
		"client-old",
		10*time.Millisecond,
	)

	if err != nil {
		t.Fatalf("failed to create expiring request: %v", err)
	}

	// Create another request for the same file that should remain active.
	_, _, err = request.GetOrCreateRequest(
		fileID,
		"client-active",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("failed to create active request: %v", err)
	}

	// Wait until the expiration timer has actually removed the old
	// request from the request module.
	//
	// We don't use a fixed Sleep because timer callbacks are
	// asynchronous. Instead, we wait for the actual state we care about.
	deadline := time.Now().Add(1 * time.Second)

	for time.Now().Before(deadline) {

		pending := request.GetPendingRequests()

		count := 0

		for _, req := range pending {
			if req.FileID == fileID {
				count++
			}
		}

		// Once only the active request remains, expiration has completed.
		if count == 1 {
			break
		}

		time.Sleep(5 * time.Millisecond)
	}

	// Calculate demand for this test's unique file.
	demand, err := CalculateDemand(fileID)

	if err != nil {
		t.Fatalf("CalculateDemand() failed: %v", err)
	}

	// The expired request should no longer contribute to demand.
	// Only client-active should remain.
	if demand != 1 {
		t.Errorf("expected demand 1, got %d", demand)
	}
}

func TestCalculateDemandEmptyFileID(t *testing.T) {
	resetDemandState()

	// An empty FileID does not identify a file, so demand calculation
	// should reject it instead of silently returning a misleading value.
	_, err := CalculateDemand("")

	if err == nil {
		t.Error("expected error for empty file ID")
	}
}
