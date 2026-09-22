package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"testing"

	"proof-of-request/model"
)

func TestSignRequest(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	req := model.Request{
		FileID:    "file-123",
		ClientID:  "client-456",
		Timestamp: 123456789,
		Status:    model.Pending,
	}

	signedReq, err := SignRequest(req, privateKey)
	if err != nil {
		t.Fatalf("SignRequest() returned an error: %v", err)
	}

	if signedReq.Request != req {
		t.Error("expected signed request to contain original request")
	}

	if len(signedReq.Signature) == 0 {
		t.Error("expected signature to be generated")
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	if !ed25519.Verify(publicKey, data, signedReq.Signature) {
		t.Error("expected generated signature to be valid")
	}
}

// changing the request would break the signature
func TestSignRequestDetectsModifiedRequest(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	req := model.Request{
		FileID:    "file-123",
		ClientID:  "client-456",
		Timestamp: 123456789,
		Status:    model.Pending,
	}

	signedReq, err := SignRequest(req, privateKey)
	if err != nil {
		t.Fatalf("SignRequest() returned an error: %v", err)
	}

	modifiedReq := req
	modifiedReq.FileID = "file-999"

	data, err := json.Marshal(modifiedReq)
	if err != nil {
		t.Fatalf("failed to marshal modified request: %v", err)
	}

	if ed25519.Verify(publicKey, data, signedReq.Signature) {
		t.Error("expected signature verification to fail for modified request")
	}
}

func TestVerifySignature(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	req := model.Request{
		FileID:    "file-123",
		ClientID:  "client-456",
		Timestamp: 123456789,
		Status:    model.Pending,
	}

	signedReq, err := SignRequest(req, privateKey)
	if err != nil {
		t.Fatalf("SignRequest() returned an error: %v", err)
	}

	if !VerifySignature(signedReq, publicKey) {
		t.Error("expected valid signature to be accepted")
	}
}

func TestVerifySignatureRejectsModifiedRequest(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	req := model.Request{
		FileID:    "file-123",
		ClientID:  "client-456",
		Timestamp: 123456789,
		Status:    model.Pending,
	}

	signedReq, err := SignRequest(req, privateKey)
	if err != nil {
		t.Fatalf("SignRequest() returned an error: %v", err)
	}

	signedReq.Request.FileID = "file-modified"

	if VerifySignature(signedReq, publicKey) {
		t.Error("expected modified request to be rejected")
	}
}
