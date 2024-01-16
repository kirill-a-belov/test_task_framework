package storage

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/test_task_framework/internal/pkg/order/model"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func NewOrderMem() *MemStorage {
	return &MemStorage{
		m: sync.Map{},
	}
}

type MemStorage struct {
	m sync.Map
}

func (s *MemStorage) Fetch(ctx context.Context, idList ...int) ([]*model.Order, error) {
	_, span := tracer.Start(ctx, "internal.pkg.order.storage.Fetch")
	defer span.End()

	result := make([]*model.Order, 0, len(idList))

	for _, id := range idList {
		if item, ok := s.m.Load(id); ok {
			order, ok := item.(*model.Order)
			if !ok {
				return nil, errors.Errorf("wrong element type by key (%v)", id)
			}
			if !order.DeletedAt.Equal(time.Time{}) {
				continue
			}

			result = append(result, order)
		}
	}

	return result, nil
}

func (s *MemStorage) Upsert(ctx context.Context, itemList ...*model.Order) error {
	_, span := tracer.Start(ctx, "internal.pkg.order.storage.Upsert")
	defer span.End()

	for _, item := range itemList {
		if item.ID == 0 {
			item.ID = rand.Intn(100)
		}
		if item.CreatedAt.Equal(time.Time{}) {
			item.CreatedAt = time.Now()
		}
		item.UpdatedAt = time.Now()

		s.m.Store(item.ID, item)
	}

	return nil
}

func (s *MemStorage) Delete(ctx context.Context, idList ...int) error {
	_, span := tracer.Start(ctx, "internal.pkg.order.storage.Delete")
	defer span.End()

	itemList, err := s.Fetch(ctx, idList...)
	if err != nil {
		return errors.Wrap(err, "delete items")
	}

	for i := range itemList {
		itemList[i].DeletedAt = time.Now()
	}

	return s.Upsert(ctx, itemList...)
}
