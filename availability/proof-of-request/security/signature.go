package security

import (
	"crypto/ed25519"
	"encoding/json"

	"proof-of-request/model"
)

func SignRequest(req model.Request, privateKey ed25519.PrivateKey) (model.SignedRequest, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return model.SignedRequest{}, err
	}

	signature := ed25519.Sign(privateKey, data)

	return model.SignedRequest{
		Request:   req,
		Signature: signature,
	}, nil
}

func VerifySignature(
	signedReq model.SignedRequest,
	publicKey ed25519.PublicKey,
) bool {
	data, err := json.Marshal(signedReq.Request)
	if err != nil {
		return false
	}

	return ed25519.Verify(
		publicKey,
		data,
		signedReq.Signature,
	)
}
