package context_helper

import (
	"context"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
	"time"
)

func RunWithTimeout(timeout time.Duration, runFunc func() error) error {
	ctx, span := tracer.Start(context.Background(), "pkg.context_helper.RunWithTimeout")
	defer span.End()

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
