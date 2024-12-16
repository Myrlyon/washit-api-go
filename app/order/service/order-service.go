package orderService

import (
	"context"
	"log"

	orderModel "washit-api/app/order/model"
	orderRepository "washit-api/app/order/repository"
)

type OrderServiceInterface interface {
	GetOrders(ctx context.Context) ([]*orderModel.Order, error)
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

func (s *OrderService) GetOrders(ctx context.Context) ([]*orderModel.Order, error) {
	Order, err := s.repository.GetOrders(ctx)
	if err != nil {
		log.Println("Failed to get Orders ", err)
		return nil, err
	}

	return Order, nil
}
