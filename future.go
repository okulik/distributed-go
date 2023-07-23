package dist

import "sync"

type Future[T any] interface {
	// Result returns the value of the future. If the future is not ready, the
	// call will block until the future is ready.
	Result() (T, error)
}

type futureImpl[T any] struct {
	once sync.Once
	wg   sync.WaitGroup

	res   T
	err   error
	resCh <-chan T
	errCh <-chan error
}

func NewFutureImpl[T any](resCh <-chan T, errCh <-chan error) Future[T] {
	return &futureImpl[T]{
		resCh: resCh,
		errCh: errCh,
	}
}

func (f *futureImpl[T]) Result() (T, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	f.wg.Wait()

	return f.res, f.err
}
