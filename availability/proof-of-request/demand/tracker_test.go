package demand

import (
	"proof-of-request/model"
	"testing"
	"time"
)

func resetRecords() {
	records = nil
}

func TestCalculateDemand(t *testing.T) {
	resetRecords()

	now := time.Now()

	RecordRequest(model.Request{
		FileID:    "file-A",
		ClientID:  "client-1",
		Timestamp: now.Add(-10 * time.Minute).UnixNano(),
	})

	RecordRequest(model.Request{
		FileID:    "file-A",
		ClientID:  "client-2",
		Timestamp: now.Add(-20 * time.Minute).UnixNano(),
	})

	RecordRequest(model.Request{
		FileID:    "file-B",
		ClientID:  "client-3",
		Timestamp: now.Add(-10 * time.Minute).UnixNano(),
	})

	if got := CalculateDemand("file-A", time.Hour); got != 2 {
		t.Errorf("expected demand 2, got %d", got)
	}

	if got := CalculateDemand("file-B", time.Hour); got != 1 {
		t.Errorf("expected demand 1, got %d", got)
	}
}

func TestCalculateDemandIgnoresOldRequests(t *testing.T) {
	resetRecords()

	now := time.Now()

	RecordRequest(model.Request{
		FileID:    "file-A",
		ClientID:  "client-old",
		Timestamp: now.Add(-2 * time.Hour).UnixNano(),
	})

	RecordRequest(model.Request{
		FileID:    "file-A",
		ClientID:  "client-new",
		Timestamp: now.Add(-10 * time.Minute).UnixNano(),
	})

	if got := CalculateDemand("file-A", time.Hour); got != 1 {
		t.Errorf("expected demand 1, got %d", got)
	}
}