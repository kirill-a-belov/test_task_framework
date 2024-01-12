package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_validate(t *testing.T) {
	testCaseList := []struct {
		name      string
		args      Config
		wantError bool
	}{
		{
			name: "Success",
			args: Config{
				Address: "localhost:1234",
				Delay:   time.Second,
				ConnTTL: time.Second,
			},
			wantError: false,
		},
		{
			name: "Invalid address",
			args: Config{
				Address: "local:host:1234",
				Delay:   time.Second,
				ConnTTL: time.Second,
			},
			wantError: true,
		},
		{
			name: "Invalid delay value",
			args: Config{
				Address: "localhost:1234",
				Delay:   time.Hour,
				ConnTTL: time.Second,
			},
			wantError: true,
		},
		{
			name: "Invalid TTL",
			args: Config{
				Address: "localhost:1234",
				Delay:   time.Second,
				ConnTTL: time.Hour,
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.args.validate()
			if tc.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}
