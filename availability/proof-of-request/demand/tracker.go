package demand

import (
	"errors"

	"proof-of-request/request"
)

// CalculateDemand returns the number of currently active requests
// for a given file.
//
// We do not pass a time window here because the request lifetime and
// demand window are the same concept in our current design.
//
// The request module automatically expires a request when its lifetime
// ends. Therefore, if a request is returned by GetPendingRequests(),
// it is already considered part of current demand.
func CalculateDemand(fileID string) (int, error) {

	// A demand calculation without a file ID has no meaningful target.
	if fileID == "" {
		return 0, errors.New("file ID cannot be empty")
	}

	// Ask the request module for the current active requests.
	//
	// The demand module does not maintain its own copy of requests.
	// This prevents the request lifecycle from being duplicated in
	// two different places.
	pendingRequests := request.GetPendingRequests()

	count := 0

	// Count only Pending requests belonging to this file.
	//
	// Since GetPendingRequests() already returns active requests,
	// there is no need to check timestamps or expiration here.
	for _, req := range pendingRequests {

		if req.FileID == fileID {
			count++
		}
	}

	return count, nil
}
