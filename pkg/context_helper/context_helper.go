package context_helper

import (
	"context"
	"time"
)

func RunWithTimeout(ctx context.Context, timeout time.Duration, runFunc func() error) error {
	ctx, _ = context.WithTimeout(ctx, timeout)

	errChan := make(chan error)
	go func() {
		errChan <- runFunc()
	}()

	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case err := <-errChan:
		return err
	}
}
