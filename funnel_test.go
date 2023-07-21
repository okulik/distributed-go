package dist_test

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"

	dist "github.com/okulik/distributed-go"
)

func TestFunnel(t *testing.T) {
	const items = 10
	const channels = 3

	chans := make([]<-chan int, 0, channels)
	for i := 0; i < channels; i++ {
		ch := make(chan int)
		chans = append(chans, ch)
		go func(i int) {
			for j := 1; j < items; j++ {
				ch <- (j * int(math.Pow(10, float64(i))))
			}
			close(ch)
		}(i)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var itemsFunneled int32
	go func() {
		defer wg.Done()
		for range dist.Funnel[int](chans[0]) {
			atomic.AddInt32(&itemsFunneled, 1)
		}
	}()
	go func() {
		defer wg.Done()
		for range dist.Funnel[int](chans[1:]...) {
			atomic.AddInt32(&itemsFunneled, 1)
		}
	}()
	wg.Wait()

	if itemsFunneled != (items-1)*channels {
		t.Error("items funneled", itemsFunneled)
	}
}
