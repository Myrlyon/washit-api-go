package orderService

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	historyModel "washit-api/internal/history/dto/model"
	orderModel "washit-api/internal/order/dto/model"
	orderRequest "washit-api/internal/order/dto/request"
	orderRepository "washit-api/internal/order/repository"
	generate "washit-api/pkg/generator"
	"washit-api/pkg/utils"

	"github.com/go-playground/validator"
)

type IOrderService interface {
	GetOrdersMe(ctx context.Context, userId string) ([]*orderModel.Order, error)
	GetOrdersAll(ctx context.Context) ([]*orderModel.Order, error)
	GetOrderById(ctx context.Context, orderId string, userId string) (*orderModel.Order, error)
	GetOrdersByUser(ctx context.Context, userId string) ([]*orderModel.Order, error)
	CreateOrder(ctx context.Context, userId int, req *orderRequest.Order) (*orderModel.Order, error)
	CancelOrder(ctx context.Context, orderId string, userId string) (*orderModel.Order, error)
	UpdateWeight(ctx context.Context, orderId string, weight string) (*orderModel.Order, error)
}

type OrderService struct {
	repository orderRepository.IOrderRepository
	validator  *validator.Validate
}

func NewOrderService(
	repository orderRepository.IOrderRepository, validator *validator.Validate) *OrderService {
	return &OrderService{
		repository: repository,
		validator:  validator,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userId int, req *orderRequest.Order) (*orderModel.Order, error) {
	if err := s.validator.Struct(&req); err != nil {
		log.Println("Failed to validate Order request ", err)
		return nil, err

	}

	order := &orderModel.Order{}
	ordId, err := generate.AlphaNumericId("ORD")
	if err != nil {
		log.Println("Failed to generate Order ID ", err)
		return nil, err
	}

	utils.CopyTo(&req, &order)
	order.ID = ordId
	order.UserID = userId

	order, err = s.repository.CreateOrder(ctx, order)
	if err != nil {
		log.Println("Failed to create Order ", err)
		return nil, fmt.Errorf("failed to create order")
	}

	return order, nil
}

func (s *OrderService) GetOrdersMe(ctx context.Context, userId string) ([]*orderModel.Order, error) {
	order, err := s.repository.GetOrders(ctx, userId)
	if err != nil {
		log.Println("Failed to get Orders me")
		return nil, fmt.Errorf("failed to get orders me")
	}

	return order, nil
}

func (s *OrderService) GetOrdersAll(ctx context.Context) ([]*orderModel.Order, error) {
	order, err := s.repository.GetOrders(ctx, "")
	if err != nil {
		log.Println("Failed to get all Orders ", err)
		return nil, fmt.Errorf("failed to get all orders")
	}

	return order, nil
}

func (s *OrderService) GetOrdersByUser(ctx context.Context, userId string) ([]*orderModel.Order, error) {
	order, err := s.repository.GetOrders(ctx, userId)
	if err != nil {
		log.Println("Failed to get Orders from userID ", err)
		return nil, fmt.Errorf("failed to get orders from userid: %v", userId)
	}

	return order, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, orderId string, userId string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(ctx, orderId, userId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) UpdateWeight(ctx context.Context, orderId string, weight string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(ctx, orderId, "0")
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", orderId)
	}

	weightFloat, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		log.Println("Failed to parse weight", err)
		return nil, fmt.Errorf("failed to parse weight: %v", weight)
	}

	order.Weight = &weightFloat

	if err := s.repository.UpdateOrder(ctx, order); err != nil {
		log.Println("Failed to update order weight by ID:", err)
		return nil, fmt.Errorf("failed to update order weight by ID: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderId string, userId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(ctx, orderId, userId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", orderId)
	}

	utils.CopyTo(&order, &history)
	history.Reason = "cancelled"
	history.DeletedAt = time.Now()

	if err := s.repository.CreateHistory(ctx, &history); err != nil {
		log.Println("Failed to move order to history", err)
		return nil, fmt.Errorf("failed to move order to history: %v", err)
	}

	if err := s.repository.DeleteOrder(ctx, order); err != nil {
		log.Println("Failed to delete order by ID:", err)
		return nil, fmt.Errorf("failed to delete order by ID: %v", orderId)
	}

	return order, nil
}
