package dist_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	dist "github.com/okulik/distributed-go"
)

var mockCircuitCalled int64

func mockCircuitForDebounce(ctx context.Context) (*string, error) {
	atomic.AddInt64(&mockCircuitCalled, 1)
	result := "debounce called"
	return &result, nil
}

func TestDebounceFirst(t *testing.T) {
	circuit := dist.DebounceFirst[string](mockCircuitForDebounce, time.Millisecond*100)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 50; i++ {
				_, _ = circuit(context.Background())
				time.Sleep(time.Millisecond * 10)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if mockCircuitCalled == 0 {
		t.Error("mockCircuit not called")
	}
	mockCircuitCalled = 0
}

func TestDebounceLast(t *testing.T) {
	circuit := dist.DebounceLast[string](mockCircuitForDebounce, time.Millisecond*100)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 50; i++ {
				_, _ = circuit(context.Background())
				time.Sleep(time.Millisecond * 10)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if mockCircuitCalled == 0 {
		t.Error("mockCircuit not called", mockCircuitCalled)
	}
	mockCircuitCalled = 0
}
