package storage

import (
	"context"

	"github.com/kirill-a-belov/test_task_framework/internal/pkg/reservation/model"
)

type Storage struct {
	Reservation ReservationRepo
}

func New() *Storage {
	return &Storage{
		Reservation: NewReservationMem(),
	}
}

type ReservationRepo interface {
	Fetch(ctx context.Context, id ...int) ([]*model.Reservation, error)
	Upsert(ctx context.Context, itemList ...*model.Reservation) error
	Delete(ctx context.Context, idList ...int) error
}
