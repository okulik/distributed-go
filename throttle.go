package dist

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Effector[T any] func(ctx context.Context) (*T, error)

func ThrottleWithRefill[T any](e Effector[T], maxTokens uint, refillTokens uint, refillDuration time.Duration) Effector[T] {
	var currentTokens = maxTokens
	var once sync.Once

	return func(ctx context.Context) (*T, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(refillDuration)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						t := currentTokens + refillTokens
						if t > maxTokens {
							t = maxTokens
						}
						currentTokens = t
					case <-ctx.Done():
						return
					}
				}
			}()
		})

		if currentTokens <= 0 {
			return nil, fmt.Errorf("throttle: too many calls")
		}

		currentTokens--

		return e(ctx)
	}
}
