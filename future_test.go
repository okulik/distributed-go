package dist_test

import (
	"context"
	"testing"
	"time"

	dist "github.com/okulik/distributed-go"
)

func TestFuture(t *testing.T) {
	ctx := context.Background()
	future := functionReturningFuture(ctx)

	_, err := future.Result()
	if err != nil {
		t.Error(err)
	}
}

func functionReturningFuture(ctx context.Context) dist.Future[string] {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <-time.After(time.Second):
			resCh <- "sleepy"
			errCh <- nil
		case <-ctx.Done():
			resCh <- ""
			errCh <- ctx.Err()
		}
	}()

	return dist.NewFutureImpl(resCh, errCh)
}
