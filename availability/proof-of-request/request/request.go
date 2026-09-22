package request

import (
	"errors"
	"sync"
	"time"

	"proof-of-request/model"
)

// requestRecord contains the request plus its lifecycle timer.
//
// We keep this separate from model.Request because the timer is an
// implementation detail of the request module. The model should only
// represent the actual request data.
type requestRecord struct {
	request model.Request
	timer   *time.Timer
}

// mu protects the request store.
//
// A request can be changed by different goroutines:
//   - the client may resolve/abort it
//   - the timer goroutine may expire it
//
// Therefore, all access to the store and request status must be protected.
var mu sync.Mutex

// We only keep currently active (Pending) requests here.
//
// We deliberately do NOT keep Resolved, Aborted, or Expired requests
// because demand is based only on currently active requests and we do
// not need request history.
var requests []*requestRecord

// CreateRequest creates a new Pending request.
//
// This function only creates the request object.
// It does not store the request or start its lifetime timer.
// Storing and lifecycle management are handled by GetOrCreateRequest.
func CreateRequest(fileID string, clientID string) (model.Request, error) {

	// The actual validation of FileID, ClientID, timestamp, etc.
	// is handled later by the validation package.
	//
	// This function is responsible only for constructing the request.
	return model.Request{
		FileID:    fileID,
		ClientID:  clientID,
		Timestamp: time.Now().UnixNano(),
		Status:    model.Pending,
	}, nil
}

// GetOrCreateRequest returns the currently active request for a
// (FileID, ClientID) pair, or creates a new one if none exists.
//
// We allow only one Pending request for the same client and file.
// This prevents repeated calls from creating multiple demand entries
// for the same active request.
func GetOrCreateRequest(
	fileID string,
	clientID string,
	lifetime time.Duration,
) (model.Request, bool, error) {

	// A request must have a meaningful lifetime.
	// Without this check, a request could have an invalid or
	// immediately meaningless lifecycle.
	if lifetime <= 0 {
		return model.Request{}, false, errors.New(
			"request lifetime must be positive",
		)
	}

	mu.Lock()
	defer mu.Unlock()

	// First check whether this client already has an active request
	// for this file.
	//
	// If it does, we reuse it instead of creating another request.
	for _, record := range requests {

		if record.request.FileID != fileID ||
			record.request.ClientID != clientID ||
			record.request.Status != model.Pending {
			continue
		}

		// The request is still active, so return the existing request.
		return record.request, false, nil
	}

	// No active request exists, so create a new one.
	req, err := CreateRequest(fileID, clientID)
	if err != nil {
		return model.Request{}, false, err
	}

	// Create an internal record so that the timer can be associated
	// with THIS exact request instance.
	//
	// This is important because the same (FileID, ClientID) pair can
	// create another request after the old one expires/resolves.
	record := &requestRecord{
		request: req,
	}

	// Store the new active request.
	requests = append(requests, record)

	// Start the request lifetime timer.
	//
	// We use a timer rather than checking the age only when somebody
	// asks for demand. Otherwise, an abandoned request could remain
	// Pending forever if nobody calls CalculateDemand again.
	//
	// The timer captures the exact record pointer. Therefore, if this
	// request disappears and a new request with the same FileID/ClientID
	// is created later, this old timer still refers only to the old record.
	record.timer = time.AfterFunc(lifetime, func() {
		expireRequest(record)
	})

	return req, true, nil
}

// transitionRecord changes a request from Pending to a terminal state.
//
// The mutex + status check are what provide correctness here.
//
// Example race:
//
//   Goroutine A: ResolveRequest()
//   Goroutine B: expiration timer fires
//
// Both may try to change the request at almost the same time.
// The first one that obtains the mutex changes Pending -> terminal.
// The second one sees that the request is no longer Pending and does nothing.
//
// Timer.Stop() alone cannot provide this guarantee.
func transitionRecord(
	record *requestRecord,
	newStatus model.RequestStatus,
) bool {

	mu.Lock()

	// Only a Pending request is allowed to transition.
	//
	// This is the critical correctness check. It guarantees that
	// Resolved, Aborted, and Expired cannot overwrite one another.
	if record.request.Status != model.Pending {
		mu.Unlock()
		return false
	}

	// Perform the lifecycle transition.
	record.request.Status = newStatus

	// Save the timer before removing the record.
	//
	// For Resolve/Abort, we will stop this timer because the request
	// has already reached its final state.
	timer := record.timer
	record.timer = nil

	// We don't need this request anymore because demand only cares
	// about currently Pending requests.
	removeRecordLocked(record)

	mu.Unlock()

	// Stop the timer after releasing the mutex.
	//
	// Stop() is an optimization here: it prevents an unnecessary timer
	// callback when possible.
	//
	// It is NOT what guarantees correctness. The status check above
	// guarantees that an already-running timer cannot overwrite the
	// new state.
	if timer != nil && newStatus != model.Expired {
		timer.Stop()
	}

	return true
}

// expireRequest is called automatically when a request's lifetime ends.
//
// The record pointer is important here: the timer can only expire the
// exact request that created that timer.
func expireRequest(record *requestRecord) {

	// Try to transition this exact request to Expired.
	//
	// If ResolveRequest or AbortRequest already changed the request,
	// transitionRecord will see that it is no longer Pending and return
	// false.
	transitionRecord(record, model.Expired)
}

// transitionRequest finds an active request using FileID + ClientID
// and attempts to move it to the requested terminal state.
func transitionRequest(
	fileID string,
	clientID string,
	newStatus model.RequestStatus,
) bool {

	mu.Lock()

	var record *requestRecord

	// Find the currently active request for this client/file pair.
	for _, candidate := range requests {

		if candidate.request.FileID == fileID &&
			candidate.request.ClientID == clientID &&
			candidate.request.Status == model.Pending {

			record = candidate
			break
		}
	}

	mu.Unlock()

	// No active request exists.
	if record == nil {
		return false
	}

	// Re-check the status under the mutex inside transitionRecord.
	//
	// Another goroutine may have changed the request between the
	// search above and this call.
	return transitionRecord(record, newStatus)
}

// ResolveRequest marks the active request as resolved.
//
// Once resolved, it is removed from the active store because it no
// longer contributes to demand.
func ResolveRequest(fileID string, clientID string) bool {

	return transitionRequest(
		fileID,
		clientID,
		model.Resolved,
	)
}

// AbortRequest marks the active request as explicitly aborted.
//
// An aborted request also stops contributing to demand, so it is
// removed from the active store.
func AbortRequest(fileID string, clientID string) bool {

	return transitionRequest(
		fileID,
		clientID,
		model.Aborted,
	)
}

// GetPendingRequests returns a snapshot of all currently active requests.
//
// This function is intentionally a READ-ONLY operation.
//
// Expiration is NOT performed here because request lifecycle management
// belongs to the request module's timers. The demand module should only
// observe the current Pending state.
func GetPendingRequests() []model.Request {

	mu.Lock()
	defer mu.Unlock()

	// Create a new slice so callers cannot modify our internal store.
	pending := make([]model.Request, 0, len(requests))

	for _, record := range requests {

		if record.request.Status == model.Pending {
			pending = append(pending, record.request)
		}
	}

	return pending
}

// removeRecordLocked removes a request from the active store.
//
// The caller must already hold mu.
//
// We remove terminal requests immediately because we explicitly decided
// that request history is not part of the demand system.
func removeRecordLocked(target *requestRecord) {

	for i, record := range requests {

		if record != target {
			continue
		}

		// Remove the record while preserving the remaining requests.
		requests = append(
			requests[:i],
			requests[i+1:]...,
		)

		return
	}
}