package orderRepository

import (
	"context"

	orderModel "washit-api/app/order/model"
	dbs "washit-api/db"
)

type OrderRepositoryInterface interface {
	GetOrders(ctx context.Context) ([]*orderModel.Order, error)
}

type OrderRepository struct {
	db dbs.DatabaseInterface
}

func NewOrderRepository(db dbs.DatabaseInterface) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrders(ctx context.Context) ([]*orderModel.Order, error) {
	var Orders []*orderModel.Order
	if err := r.db.Find(ctx, &Orders, dbs.WithLimit(10), dbs.WithOrder("id")); err != nil {
		return nil, err
	}

	return Orders, nil
}
