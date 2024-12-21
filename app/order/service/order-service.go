package orderService

import (
	"context"
	"fmt"
	"log"

	historyModel "washit-api/app/history/dto/model"
	orderModel "washit-api/app/order/dto/model"
	orderRequest "washit-api/app/order/dto/request"
	orderRepository "washit-api/app/order/repository"
	"washit-api/utils"
)

type OrderServiceInterface interface {
	GetOrders(ctx context.Context, userId string) ([]*orderModel.Order, error)
	GetOrdersAll(ctx context.Context) ([]*orderModel.Order, error)
	GetOrderById(ctx context.Context, orderId string, userId string) (*orderModel.Order, error)
	GetOrdersByUser(ctx context.Context, userId string) ([]*orderModel.Order, error)
	CreateOrder(ctx context.Context, userId int, req *orderRequest.Order) (*orderModel.Order, error)
	CancelOrder(ctx context.Context, orderId string, userId string) (*orderModel.Order, error)
}

type OrderService struct {
	repository orderRepository.OrderRepositoryInterface
}

func NewOrderService(
	repository orderRepository.OrderRepositoryInterface) *OrderService {
	return &OrderService{
		repository: repository,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userId int, req *orderRequest.Order) (*orderModel.Order, error) {
	ordId, err := utils.AlphaNumericId("ORD")
	if err != nil {
		log.Println("Failed to generate Order ID ", err)
		return nil, err
	}

	order := &orderModel.Order{
		ID:            ordId,
		UserID:        userId,
		TransactionID: req.TransactionID,
		AddressID:     req.AddressID,
		Status:        req.Status,
		Note:          req.Note,
		ServiceType:   req.ServiceType,
		OrderType:     req.OrderType,
		Price:         req.Price,
		CollectDate:   req.CollectDate,
		EstimateDate:  req.EstimateDate,
	}

	order, err = s.repository.CreateOrder(ctx, order)
	if err != nil {
		log.Println("Failed to create Order ", err)
		return nil, fmt.Errorf("failed to create order")
	}

	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context, userId string) ([]*orderModel.Order, error) {
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

func (s *OrderService) CancelOrder(ctx context.Context, orderId string, userId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(ctx, orderId, userId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", orderId)
	}

	utils.CopyTo(&order, &history)

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
