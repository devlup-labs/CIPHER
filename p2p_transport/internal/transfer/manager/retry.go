package manager

import (
	"context"
	"time"
)

type RetryPolicy struct {
	MaxRetries int
	RetryDelay time.Duration
	Backoff    float64
}

var DefaultRetryPolicy = RetryPolicy{
	MaxRetries: 3,
	RetryDelay: 2 * time.Second,
	Backoff:    1.5,
}

// ExecuteWithRetry attempts to execute the given operation according to the policy.
func ExecuteWithRetry(ctx context.Context, policy RetryPolicy, op func() error) error {
	var lastErr error
	delay := policy.RetryDelay

	for i := 0; i <= policy.MaxRetries; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = op()
		if lastErr == nil {
			return nil
		}

		if i < policy.MaxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * policy.Backoff)
		}
	}

	return lastErr
}
