package order

import (
	"context"

	"github.com/pkg/errors"

	accountModel "github.com/kirill-a-belov/test_task_framework/internal/pkg/account/model"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/order/model"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/order/storage"
	userModel "github.com/kirill-a-belov/test_task_framework/internal/pkg/user/model"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

type userFetcher interface {
	Fetch(ctx context.Context, id int) (*userModel.User, error)
}

type accountFetcher interface {
	Fetch(ctx context.Context, id int) (*accountModel.Account, error)
}

type fetcherStub[A any] struct {
	item *A
}

func (fs *fetcherStub[A]) Fetch(ctx context.Context, id int) (*A, error) {
	//TODO(KB): Implement
	return fs.item, nil
}

type Module struct {
	storage        *storage.Storage
	userFetcher    userFetcher
	accountFetcher accountFetcher
}

func New() *Module {
	return &Module{
		storage: storage.New(),
		userFetcher: &fetcherStub[userModel.User]{
			item: &userModel.User{Email: "example@email.com"},
		},
		accountFetcher: &fetcherStub[accountModel.Account]{
			item: &accountModel.Account{CreditCard: "12345"},
		},
	}
}

type CreateOrderRequest struct {
	UserID int

	model.OrderDetails
}
type CreateOrderResponse struct {
	ID int
}

func (m *Module) CreateOrder(ctx context.Context, request *CreateOrderRequest) (*CreateOrderResponse, error) {
	_, span := tracer.Start(ctx, "internal.pkg.order.Module.CreateOrder")
	defer span.End()

	if err := request.Validate(); err != nil {
		return nil, errors.Wrap(err, "create order request validation")
	}

	// TODO(KB): Should be removed after DB refactoring
	user, err := m.userFetcher.Fetch(ctx, request.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching user")
	}
	account, err := m.accountFetcher.Fetch(ctx, request.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching account")
	}

	newOrder := &model.Order{
		OrderDetails: request.OrderDetails,
		UserEmail:    user.Email,
		CreditCard:   account.CreditCard,
	}

	if err := m.storage.Order.Upsert(ctx, newOrder); err != nil {
		return nil, errors.Wrap(err, "creating order")
	}

	return &CreateOrderResponse{
		ID: newOrder.ID,
	}, nil
}
