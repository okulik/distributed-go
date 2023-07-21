package dist

import "sync"

// Funnel takes a variadic number of channels and returns a single channel that
// will receive all values from all channels. The returned channel will be
// closed when all source channels are closed. The order of values received on
// the returned channel is not guaranteed.
func Funnel[T any](sources ...<-chan T) <-chan T {
	dest := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(sources))
	for _, src := range sources {
		go func(c <-chan T) {
			defer wg.Done()
			for v := range c {
				dest <- v
			}
		}(src)
	}

	go func() {
		defer close(dest)
		wg.Wait()
	}()

	return dest
}
