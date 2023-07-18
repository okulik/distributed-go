package dist_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	dist "github.com/okulik/distributed-go"
)

func mockCircuitForBreaker(ctx context.Context) (*string, error) {
	atomic.AddInt64(&mockCircuitCalled, 1)
	return nil, errors.New("mockCircuitForBreaker")
}

func TestBreaker(t *testing.T) {
	circuit := dist.Breaker[string](mockCircuitForBreaker, 5)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				_, _ = circuit(context.Background())
				time.Sleep(time.Millisecond * 100)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if mockCircuitCalled != 5 {
		t.Error("mockCircuitCalled", mockCircuitCalled)
	}
	mockCircuitCalled = 0
}
