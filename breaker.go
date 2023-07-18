package dist

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Breaker function represents a circuit breaker. It wraps a circuit function
// and returns a new circuit function that will return an error if the circuit
// function fails more than failureThreshold times in a row.
func Breaker[T any](circuit Circuit[T], failureThreshold uint) Circuit[T] {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (*T, error) {
		m.RLock()
		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return nil, errors.New("service unreachable")
			}
		}
		m.RUnlock()

		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()
		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0
		return response, nil
	}
}
