package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
				Port:         1,
				ConnPoolSize: 2,
				ConnTTL:      time.Millisecond,
			},
			wantError: false,
		},
		{
			name: "Port negative value",
			args: Config{
				Port:         -1,
				ConnPoolSize: 2,
				ConnTTL:      time.Millisecond,
			},
			wantError: true,
		},
		{
			name: "Port above maximum",
			args: Config{
				Port:         66000,
				ConnPoolSize: 2,
				ConnTTL:      time.Millisecond,
			},
			wantError: true,
		},
		{
			name: "Connection pool size less than one",
			args: Config{
				Port:         1,
				ConnPoolSize: 0,
				ConnTTL:      time.Millisecond,
			},
			wantError: true,
		},
		{
			name: "Connection pool size above maximum",
			args: Config{
				Port:         1,
				ConnPoolSize: 10000,
				ConnTTL:      time.Millisecond,
			},
			wantError: true,
		},
		{
			name: "TTL too small",
			args: Config{
				Port:         1,
				ConnPoolSize: 2,
				ConnTTL:      time.Microsecond,
			},
			wantError: true,
		},
		{
			name: "TTL above maximum",
			args: Config{
				Port:         1,
				ConnPoolSize: 2,
				ConnTTL:      time.Hour,
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
