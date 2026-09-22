package model

type SignedRequest struct {
	Request   Request
	Signature []byte
}