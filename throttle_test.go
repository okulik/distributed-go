package dist_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	dist "github.com/okulik/distributed-go"
)

var mockEffectorCalled int64

func mockEffector(ctx context.Context) (*string, error) {
	atomic.AddInt64(&mockEffectorCalled, 1)
	result := "effector called"
	return &result, nil
}

func TestThrottleWithRefill_NoRefill(t *testing.T) {
	effector := dist.ThrottleWithRefill(mockEffector, 10, 1, time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 30; i++ {
				_, _ = effector(context.Background())
				time.Sleep(time.Millisecond * 10)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if mockEffectorCalled != 10 {
		t.Error("mockEffectorCalled", mockEffectorCalled)
	}
	mockEffectorCalled = 0
}

func TestThrottleWithRefill_Refill(t *testing.T) {
	effector := dist.ThrottleWithRefill(mockEffector, 10, 5, time.Millisecond*200)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 30; i++ {
				_, _ = effector(context.Background())
				time.Sleep(time.Millisecond * 10)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if mockEffectorCalled != 15 {
		t.Error("mockEffectorCalled", mockEffectorCalled)
	}
	mockEffectorCalled = 0
}
