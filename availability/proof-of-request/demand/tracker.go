package demand

import (
	"proof-of-request/model"
	"sync"
)

var (
	mu      sync.Mutex
	records []model.Request
)

// RecordRequest stores an accepted request for demand tracking.
func RecordRequest(req model.Request) {
	mu.Lock()
	defer mu.Unlock()

	records = append(records, req)
}

func CalculateDemand(fileID string, startTime int64, endTime int64) int {
	mu.Lock()
	defer mu.Unlock()

	count := 0

	for _, req := range records {
		if req.FileID == fileID &&
			req.Timestamp >= startTime &&
			req.Timestamp < endTime {
			count++
		}
	}

	return count
}
