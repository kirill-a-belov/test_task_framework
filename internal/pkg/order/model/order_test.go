package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Validate(t *testing.T) {
	testCaseList := []struct {
		name      string
		args      Order
		wantError bool
	}{
		{
			name: "Success",
			args: Order{
				OrderDetails: OrderDetails{
					HotelID: "123",
					RoomID:  "456",
					From:    time.Now(),
					To:      time.Now().Add(24 * time.Hour),
				},
				UserEmail:  "example@email.com",
				CreditCard: "5105105105105100",
			},
			wantError: false,
		},
		{
			name: "Invalid email",
			args: Order{
				OrderDetails: OrderDetails{
					HotelID: "123",
					RoomID:  "456",
					From:    time.Now(),
					To:      time.Now().Add(24 * time.Hour),
				},
				UserEmail:  "ex@ample@email.com",
				CreditCard: "5105105105105100",
			},
			wantError: true,
		},
		{
			name: "Invalid credit card",
			args: Order{
				OrderDetails: OrderDetails{
					HotelID: "123",
					RoomID:  "456",
					From:    time.Now(),
					To:      time.Now().Add(24 * time.Hour),
				},
				UserEmail:  "example@email.com",
				CreditCard: "5105105fff105105100",
			},
			wantError: true,
		},
		{
			name: "Invalid timeframe: from > to",
			args: Order{
				OrderDetails: OrderDetails{
					HotelID: "123",
					RoomID:  "456",
					From:    time.Now().Add(24 * time.Hour),
					To:      time.Now(),
				},
				UserEmail:  "example@email.com",
				CreditCard: "5105105105105100",
			},
			wantError: true,
		},
		{
			name: "Invalid timeframe: from > to",
			args: Order{
				OrderDetails: OrderDetails{
					HotelID: "123",
					RoomID:  "456",
					From:    time.Now(),
					To:      time.Now().Add(4 * time.Hour),
				},
				UserEmail:  "example@email.com",
				CreditCard: "5105105105105100",
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.args.Validate()
			if tc.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}
