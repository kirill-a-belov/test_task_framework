package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type OrderDetails struct {
	HotelID string    `json:"hotel_id"` // TBD: consider numeric value
	RoomID  string    `json:"room_id"`  // TBD: consider numeric value
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}

func (o *OrderDetails) Validate() error {
	// TODO(KB): use global configuration
	const minBookingLengthHours = 24

	if o.From.After(o.To) || (o.To.Sub(o.From) < minBookingLengthHours*time.Hour) {
		return errors.New("invalid time range")
	}

	return nil
}

type Order struct {
	ID int

	OrderDetails

	UserEmail  string `validate:"email" json:"email"`             // TBD: Abstraction leak, should be in Reservation or user ID only if needed
	CreditCard string `validate:"credit_card" json:"credit_card"` // TBD: Abstraction leak - should be in accounting module, Account id in Reservation

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (o *Order) Validate() error {
	if err := o.OrderDetails.Validate(); err != nil {
		return errors.Wrap(err, "order details validate")
	}

	return validator.New(validator.WithRequiredStructEnabled()).Struct(o)
}
