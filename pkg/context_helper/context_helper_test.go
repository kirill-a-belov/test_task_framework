package context_helper

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRunWithTimeout(t *testing.T) {
	// Context done
	// Func done with error
	// Context timeout exceeded

	testCaseList := []struct {
		name      string
		args      func() (context.Context, time.Duration, func() error)
		wantError bool
	}{
		{
			name: "Success",
			args: func() (context.Context, time.Duration, func() error) {
				ctx := context.Background()
				timeout := time.Second
				testFunc := func() error {
					return nil
				}

				return ctx, timeout, testFunc
			},
			wantError: false,
		},
		{
			name: "Error in func",
			args: func() (context.Context, time.Duration, func() error) {
				ctx := context.Background()
				timeout := time.Second
				testFunc := func() error {
					return errors.New("example error")
				}

				return ctx, timeout, testFunc
			},
			wantError: true,
		},
		{
			name: "Context done",
			args: func() (context.Context, time.Duration, func() error) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				timeout := time.Second
				testFunc := func() error {
					return nil
				}

				return ctx, timeout, testFunc
			},
			wantError: true,
		},
		{
			name: "Timeout exceeded",
			args: func() (context.Context, time.Duration, func() error) {
				ctx := context.Background()
				timeout := time.Millisecond
				testFunc := func() error {
					time.Sleep(time.Second)
					return nil
				}

				return ctx, timeout, testFunc
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			err := RunWithTimeout(tc.args())
			if tc.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}
