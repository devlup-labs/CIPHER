package validation

import (
	"testing"
	"time"

	"proof-of-request/model"
)

func validRequest() model.Request {
	return model.Request{
		FileID:    "file-123",
		ClientID:  "client-456",
		Timestamp: time.Now().Unix(),
		Status:    model.Pending,
	}
}

func TestValidateRequest(t *testing.T) {
	req := validRequest()

	if !ValidateRequest(req) {
		t.Error("expected valid request to pass validation")
	}
}

func TestValidateRequestRejectsEmptyFileID(t *testing.T) {
	req := validRequest()
	req.FileID = ""

	if ValidateRequest(req) {
		t.Error("expected request with empty FileID to be rejected")
	}
}

func TestValidateRequestRejectsEmptyClientID(t *testing.T) {
	req := validRequest()
	req.ClientID = ""

	if ValidateRequest(req) {
		t.Error("expected request with empty ClientID to be rejected")
	}
}

func TestValidateRequestRejectsInvalidTimestamp(t *testing.T) {
	req := validRequest()
	req.Timestamp = 0

	if ValidateRequest(req) {
		t.Error("expected request with invalid timestamp to be rejected")
	}
}

func TestValidateRequestRejectsNonPendingRequest(t *testing.T) {
	req := validRequest()
	req.Status = model.Resolved

	if ValidateRequest(req) {
		t.Error("expected resolved request to be rejected")
	}
}
