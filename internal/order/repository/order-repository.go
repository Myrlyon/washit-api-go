package orderRepository

import (
	"context"

	historyModel "washit-api/internal/history/dto/model"
	orderModel "washit-api/internal/order/dto/model"
	"washit-api/pkg/db/dbs"
)

type IOrderRepository interface {
	GetAllOrders(ctx context.Context) ([]*orderModel.Order, error)
	GetOrdersByUser(ctx context.Context, userID int64) ([]*orderModel.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*orderModel.Order, error)
	CreateOrder(ctx context.Context, order *orderModel.Order) (*orderModel.Order, error)
	CreateHistory(ctx context.Context, history *historyModel.History) error
	DeleteOrder(ctx context.Context, order *orderModel.Order) error
	UpdateOrder(ctx context.Context, order *orderModel.Order) error
}

type OrderRepository struct {
	db dbs.IDatabase
}

func NewOrderRepository(db dbs.IDatabase) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *orderModel.Order) (*orderModel.Order, error) {
	if err := r.db.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]*orderModel.Order, error) {
	var orders []*orderModel.Order
	query := []dbs.FindOption{
		dbs.WithLimit(10),
		dbs.WithOrder("created_at DESC"),
		dbs.WithPreload([]string{"User"}),
	}

	if err := r.db.Find(ctx, &orders, query...); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetOrdersByUser(ctx context.Context, userID int64) ([]*orderModel.Order, error) {
	var orders []*orderModel.Order
	query := []dbs.FindOption{
		dbs.WithLimit(10),
		dbs.WithOrder("created_at DESC"),
		dbs.WithPreload([]string{"User"}),
	}

	if userID != 0 {
		query = append(query, dbs.WithQuery(dbs.NewQuery("user_id = ?", userID)))
	}

	if err := r.db.Find(ctx, &orders, query...); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*orderModel.Order, error) {
	var order orderModel.Order
	if err := r.db.FindByID(ctx, orderID, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) CreateHistory(ctx context.Context, history *historyModel.History) error {
	if err := r.db.Create(ctx, history); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, order *orderModel.Order) error {
	if err := r.db.Delete(ctx, order); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *orderModel.Order) error {
	if err := r.db.Update(ctx, order); err != nil {
		return err
	}

	return nil
}
