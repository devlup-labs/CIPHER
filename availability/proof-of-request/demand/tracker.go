package demand

import (
	"proof-of-request/model"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	records []model.Request
)

func RecordRequest(req model.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Store the validated request for future demand calculation.
	records = append(records, req)
}

func CalculateDemand(fileID string, window time.Duration) int {
	mu.Lock()
	defer mu.Unlock()

	// Count only active, unfulfilled requests within the demand window.
	cutoff := time.Now().Add(-window)

	count := 0

	for _, req := range records {
		if req.FileID != fileID || req.Status != model.Pending {
			continue
		}

		requestTime := time.Unix(0, req.Timestamp)

		if requestTime.After(cutoff) {
			count++
		}
	}

	return count
}