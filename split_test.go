package dist_test

import (
	"sync"
	"sync/atomic"
	"testing"

	dist "github.com/okulik/distributed-go"
)

func TestSplit(t *testing.T) {
	const destinations = 2
	const items = 10

	src := make(chan int)
	dests := dist.Split[int](src, destinations)

	go func() {
		defer close(src)
		for i := 0; i < items; i++ {
			src <- i
		}
	}()

	var wg sync.WaitGroup
	wg.Add(destinations)
	var itemsSent []int32 = make([]int32, destinations)
	for i, dest := range dests {
		go func(i int, dest <-chan int) {
			defer wg.Done()
			for range dest {
				atomic.AddInt32(&itemsSent[i], 1)
			}
		}(i, dest)
	}
	wg.Wait()

	var itemsSentTotal int32
	for _, itemsSent := range itemsSent {
		itemsSentTotal += itemsSent
	}

	if itemsSentTotal != items {
		t.Error("items sent", itemsSent)
	}
}
