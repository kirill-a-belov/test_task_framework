package calendar

import (
	"context"

	"github.com/kirill-a-belov/test_task_framework/internal/pkg/reservation/model"
)

type Module struct{}

func New() *Module {
	return &Module{}
}

type CheckAvailabilityRequest struct {
	model.ReservationDetails
}

func (*Module) CheckAvailability(ctx context.Context, request CheckAvailabilityRequest) (bool, error) {
	//TODO(KB): Implement
	return true, nil
}
