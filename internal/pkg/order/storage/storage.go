package storage

import (
	"context"

	"github.com/kirill-a-belov/test_task_framework/internal/pkg/order/model"
)

type Storage struct {
	Order OrderRepo
}

func New() *Storage {
	return &Storage{
		Order: NewOrderMem(),
	}
}

type OrderRepo interface {
	Fetch(ctx context.Context, id ...int) ([]*model.Order, error)
	Upsert(ctx context.Context, itemList ...*model.Order) error
	Delete(ctx context.Context, idList ...int) error
}
