package orderService

import (
	"context"
	"log"

	orderModel "washit-api/app/order/dto/model"
	orderRequest "washit-api/app/order/dto/request"
	orderRepository "washit-api/app/order/repository"
	"washit-api/utils"
)

type OrderServiceInterface interface {
	GetOrders(ctx context.Context, userId string) ([]*orderModel.Order, error)
	GetOrderById(ctx context.Context, orderId string, userId string) (*orderModel.Order, error)
	CreateOrder(ctx context.Context, userId int, req *orderRequest.Order) (*orderModel.Order, error)
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
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context, userId string) ([]*orderModel.Order, error) {
	order, err := s.repository.GetOrders(ctx, userId)
	if err != nil {
		log.Println("Failed to get Orders ", err)
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, orderId string, userId string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(ctx, orderId, userId)
	if err != nil {
		log.Println("Failed to get Orders ", err)
		return nil, err
	}

	return order, nil
}
