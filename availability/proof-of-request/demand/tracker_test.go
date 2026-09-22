package demand

import (
	"testing"

	"proof-of-request/model"
)

func TestRecordRequest(t *testing.T) {
	records = nil

	req := model.Request{
		FileID:    "file-123",
		ClientID:  "client-456",
		Timestamp: 123456789,
		Status:    model.Pending,
	}

	RecordRequest(req)

	if len(records) != 1 {
		t.Fatalf("expected 1 recorded request, got %d", len(records))
	}

	if records[0] != req {
		t.Error("expected recorded request to match original request")
	}
}
