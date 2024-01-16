package reservation

import (
	"context"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/test_task_framework/internal/pkg/calendar"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/order"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/reservation/model"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/reservation/storage"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

type orderCreator interface {
	CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
}

type availabilityChecker interface {
	CheckAvailability(ctx context.Context, request calendar.CheckAvailabilityRequest) (bool, error)
}

type Module struct {
	storage             *storage.Storage
	orderCreator        orderCreator
	availabilityChecker availabilityChecker
}

func New() *Module {
	return &Module{
		storage:             storage.New(),
		orderCreator:        order.New(),
		availabilityChecker: calendar.New(),
	}
}

type CreateReservationRequest struct {
	UserID int
	model.ReservationDetails
}
type CreateReservationResponse struct {
	model.Reservation
}

func (m *Module) CreateReservation(ctx context.Context, request *CreateReservationRequest) (*CreateReservationResponse, error) {
	_, span := tracer.Start(ctx, "internal.pkg.order.Module.CreateReservation")
	defer span.End()

	if err := request.Validate(); err != nil {
		return nil, errors.Wrap(err, "create reservation request validation")
	}

	isAvailible, err := m.availabilityChecker.CheckAvailability(ctx, calendar.CheckAvailabilityRequest{
		ReservationDetails: request.ReservationDetails,
	})
	if err != nil {
		return nil, errors.Wrap(err, "availability checking")
	}
	if !isAvailible {
		return nil, errors.New("dates are not available")
	}

	createOrderResponse, err := m.orderCreator.CreateOrder(ctx, &order.CreateOrderRequest{
		UserID:       request.UserID,
		OrderDetails: request.ReservationDetails.OrderDetails,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating order for reservation")
	}

	newReservation := &model.Reservation{
		OrderID:            createOrderResponse.ID,
		ReservationDetails: request.ReservationDetails,
	}
	if err := m.storage.Reservation.Upsert(ctx, newReservation); err != nil {
		return nil, errors.Wrap(err, "creating reservation")
	}

	return &CreateReservationResponse{
		Reservation: *newReservation,
	}, nil
}
