package dist

import (
	"context"
	"sync"
	"time"
)

// DebounceFirst returns a circuit that wraps the given circuit function and debounces its output.
// The returned circuit waits for a duration of `d` before calling the wrapped circuit.
// If the wrapped circuit is called again before the duration has elapsed, the previous result
// is returned. The returned circuit returns the first output of the wrapped circuit after the
// duration has elapsed.
// The function uses a mutex to protect the result and err variables, which are used to store the
// last result and error from the wrapped circuit. The mutex is locked and unlocked around the
// critical sections of the function to ensure that they are thread-safe.
func DebounceFirst[T any](circuit Circuit[T], d time.Duration) Circuit[T] {
	var threshold time.Time
	var result *T
	var err error
	var mu sync.Mutex

	return func(ctx context.Context) (*T, error) {
		if time.Now().Before(threshold) {
			return result, err
		}

		mu.Lock()
		defer mu.Unlock()

		result, err = circuit(ctx)
		threshold = time.Now().Add(d)

		return result, err
	}
}

// DebounceLast returns a circuit that wraps the given circuit function  and debounces its output.
// The returned circuit waits for a duration of `d` before calling the wrapped circuit.
// If the wrapped circuit is called again before the duration has elapsed, the timer is reset.
// The returned circuit returns the last output of the wrapped circuit after the duration has
// elapsed.
// The function uses a ticker to handle the debounce timer. The ticker is created using sync.Once
// to ensure that it is only created once. When the ticker fires, the function checks if the
// threshold time has elapsed. If it has, the wrapped circuit is called and the result is
// returned. If the context is cancelled, an error is returned.
// The function uses a mutex to protect the result and err variables, which are used to store the
// last result and error from the wrapped circuit. The mutex is locked and unlocked around the
// critical sections of the function to ensure that they are thread-safe.
func DebounceLast[T any](circuit Circuit[T], d time.Duration) Circuit[T] {
	var ticker *time.Ticker
	var result *T
	var err error
	var once sync.Once
	var mu sync.Mutex

	return func(ctx context.Context) (*T, error) {
		mu.Lock()
		defer mu.Unlock()

		threshold := time.Now().Add(d)

		once.Do(func() {
			ticker = time.NewTicker(time.Millisecond * 100)

			go func() {
				defer func() {
					mu.Lock()
					ticker.Stop()
					once = sync.Once{}
					mu.Unlock()
				}()
				for {
					select {
					case <-ticker.C:
						if time.Now().After(threshold) {
							mu.Lock()
							result, err = circuit(ctx)
							mu.Unlock()
							return
						}
					case <-ctx.Done():
						mu.Lock()
						result, err = nil, ctx.Err()
						mu.Unlock()
						return
					}
				}
			}()
		})

		return result, err
	}
}
