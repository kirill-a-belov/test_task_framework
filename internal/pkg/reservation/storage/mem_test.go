package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-a-belov/test_task_framework/internal/pkg/reservation/model"
)

func TestMemStorage_Fetch(t *testing.T) {
	testCaseList := []struct {
		name      string
		args      func() (input []int, result []*model.Reservation, storage *MemStorage)
		wantError bool
	}{
		{
			name: "Success",
			args: func() (input []int, result []*model.Reservation, storage *MemStorage) {
				testList := []*model.Reservation{
					{
						ID: 1,
					},
					{
						ID:        2,
						DeletedAt: time.Now(),
					},
					{
						ID: 3,
					},
				}
				input = []int{testList[0].ID, testList[1].ID, testList[2].ID}
				result = []*model.Reservation{testList[0], testList[2]}

				storage = NewReservationMem()
				for _, item := range testList {
					storage.m.Store(item.ID, item)
				}

				return input, result, storage
			},
		},
		{
			name: "Wrong item",
			args: func() (input []int, result []*model.Reservation, storage *MemStorage) {
				testList := []*model.Reservation{
					{
						ID: 1,
					},
					{
						ID:        2,
						DeletedAt: time.Now(),
					},
					{
						ID: 3,
					},
				}
				input = []int{testList[0].ID, testList[1].ID, testList[2].ID, 4}
				result = []*model.Reservation{testList[0], testList[2]}

				storage = NewReservationMem()
				for _, item := range testList {
					storage.m.Store(item.ID, item)
				}
				storage.m.Store(4, struct{}{})

				return input, result, storage
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			input, result, storage := tc.args()

			r, err := storage.Fetch(context.Background(), input...)
			if tc.wantError {
				assert.Error(t, err)

				return
			}
			require.NoError(t, err)
			assert.Equal(t, result, r)
		})
	}
}

func TestMemStorage_Upsert(t *testing.T) {
	// TODO(KB): Implement
}

func TestMemStorage_Delete(t *testing.T) {
	// TODO(KB): Implement
}
