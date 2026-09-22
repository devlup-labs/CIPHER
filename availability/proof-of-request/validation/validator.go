package validation

import "proof-of-request/model"

func ValidateRequest(req model.Request) bool {
	if req.FileID == "" {
		return false
	}

	if req.ClientID == "" {
		return false
	}

	if req.Timestamp <= 0 {
		return false
	}

	if req.Status != model.Pending {
		return false
	}

	return true
}
