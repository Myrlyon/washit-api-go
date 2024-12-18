package orderRepository

import (
	"context"
	"errors"
	"log"
	"strconv"

	orderModel "washit-api/app/order/dto/model"
	dbs "washit-api/db"
)

type OrderRepositoryInterface interface {
	GetOrders(ctx context.Context, userId string) ([]*orderModel.Order, error)
	GetOrderById(ctx context.Context, orderId string, userId string) (*orderModel.Order, error)
	CreateOrder(ctx context.Context, order *orderModel.Order) (*orderModel.Order, error)
}

type OrderRepository struct {
	db dbs.DatabaseInterface
}

func NewOrderRepository(db dbs.DatabaseInterface) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *orderModel.Order) (*orderModel.Order, error) {
	if err := r.db.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetOrders(ctx context.Context, userId string) ([]*orderModel.Order, error) {
	var orders []*orderModel.Order
	query := []dbs.FindOption{
		dbs.WithLimit(10),
		dbs.WithOrder("created_at DESC"),
		dbs.WithPreload([]string{"User"}),
	}

	log.Println("user" + userId)

	if userId != "0" {
		query = append(query, dbs.WithQuery(dbs.NewQuery("user_id = ?", userId)))
	}

	if err := r.db.Find(ctx, &orders, query...); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetOrderById(ctx context.Context, orderId string, userId string) (*orderModel.Order, error) {
	var order orderModel.Order

	if err := r.db.FindById(ctx, orderId, &order); err != nil {
		return nil, err
	}

	if userId != "0" && strconv.Itoa(order.UserID) != userId {
		return nil, errors.New("unauthorized")
	}

	return &order, nil
}
