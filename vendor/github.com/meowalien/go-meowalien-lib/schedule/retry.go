package schedule

import (
	"context"
	"fmt"
	"time"
)

// Retryer is a retryer that will retry the function if it returns an error
// TryCount is the number of times to retry the function
// RetryInterval is the interval between retries
// If TryCount is zero, just call the function once
// If TryCount is -1, it will retry infinitely
// If RetryInterval is zero, it will retry immediately
type Retryer struct {
	TryCount      int
	RetryInterval time.Duration
}

// Try round will start from 1
func (r *Retryer) Try(ctx context.Context, f func(ctx context.Context) error) (err error) {
	if r.TryCount == -1 {
		return r.retryInfinitely(ctx, f)
	}

	// If tryCount is zero, just call the function once
	if r.TryCount == 0 {
		return f(ctx)
	}
	for {
		// Attempt to execute the function and decrement the retry count
		errA := f(ctx)
		if errA == nil {
			return nil // Success, return nil
		}
		// Append the error to the existing error chain
		if err == nil {
			err = errA
		} else {
			err = fmt.Errorf("%w -> %v", err, errA)
		}

		r.TryCount--
		if r.TryCount <= 0 {
			return err // Return the final error if tryCount is exhausted
		}
		// Wait for either the context to be done or the retry interval to elapse
		select {
		case <-ctx.Done():
			return
		case <-time.After(r.RetryInterval):
			continue
		}
	}
}

func (r *Retryer) retryInfinitely(ctx context.Context, f func(ctx context.Context) error) (err error) {
	for {
		// Attempt to execute the function and decrement the retry count
		errA := f(ctx)
		if errA == nil {
			return nil // Success, return nil
		}

		// Wait for either the context to be done or the retry interval to elapse
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.RetryInterval):
			continue
		}
	}
}
