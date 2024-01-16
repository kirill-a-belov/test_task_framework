package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/test_task_framework/internal/pkg/order/model"
)

type ReservationDetails struct {
	model.OrderDetails // TBD: Should be replaced by OrderID
}

type Reservation struct {
	ID      int
	OrderID int

	ReservationDetails

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (o *Reservation) Validate() error {
	if err := o.ReservationDetails.Validate(); err != nil {
		return errors.Wrap(err, "reservation details validate")
	}

	return validator.New(validator.WithRequiredStructEnabled()).Struct(o)
}
