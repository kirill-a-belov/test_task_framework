package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/config"
)

func TestServer_Start(t *testing.T) {
	testCaseList := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "Success",
			wantError: false,
		},
		// TODO(KB): Check if register route failed
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := New(ctx, &config.Config{})

			err := s.Start(ctx)
			defer s.Stop(ctx)

			if tc.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestServer_Stop(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		s := New(ctx, &config.Config{})

		s.Stop(ctx)

		_, closed := <-s.stopChan
		assert.False(t, closed)
	})
}
