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
	GetOrdersMe(c context.Context, userId string) ([]*orderModel.Order, error)
	GetOrdersAll(c context.Context) ([]*orderModel.Order, error)
	GetOrderById(c context.Context, orderId string, userId string) (*orderModel.Order, error)
	GetOrdersByUser(c context.Context, userId string) ([]*orderModel.Order, error)
	CreateOrder(c context.Context, userId int, req *orderRequest.Order) (*orderModel.Order, error)
	CancelOrder(c context.Context, orderId string, userId string) (*orderModel.Order, error)
	UpdateWeight(c context.Context, orderId string, weight string) (*orderModel.Order, error)
	AcceptOrder(c context.Context, orderId string) (*orderModel.Order, error)
	CompleteOrder(c context.Context, orderId string, userId string) (*orderModel.Order, error)
	PayOrder(c context.Context, orderId string, req *orderRequest.Payment) (*orderModel.Order, error)
	RejectOrder(c context.Context, orderId string) (*orderModel.Order, error)
	EditOrder(c context.Context, orderId string, userId string, req *orderRequest.Order) (*orderModel.Order, error)
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

func (s *OrderService) CreateOrder(c context.Context, userId int, req *orderRequest.Order) (*orderModel.Order, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Printf("Failed to validate Order request: %v", err)
		return nil, err
	}

	order := &orderModel.Order{}
	ordId, err := generate.AlphaNumericId("ORD")
	if err != nil {
		log.Printf("Failed to generate Order ID: %v", err)
		return nil, err
	}

	utils.CopyTo(req, order)
	order.ID = ordId
	order.UserID = userId
	order.Status = "created"

	createdOrder, err := s.repository.CreateOrder(c, order)
	if err != nil {
		log.Printf("Failed to create Order: %v", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return createdOrder, nil
}

func (s *OrderService) GetOrdersMe(c context.Context, userId string) ([]*orderModel.Order, error) {
	orders, err := s.repository.GetOrdersByUser(c, userId)
	if err != nil {
		log.Printf("Failed to get orders for user %s: %v", userId, err)
		return nil, fmt.Errorf("failed to get orders for user %s: %w", userId, err)
	}

	return orders, nil
}

func (s *OrderService) GetOrdersAll(c context.Context) ([]*orderModel.Order, error) {
	orders, err := s.repository.GetAllOrders(c)
	if err != nil {
		log.Printf("Failed to get all Orders: %v", err)
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}

	return orders, nil
}

func (s *OrderService) GetOrdersByUser(c context.Context, userId string) ([]*orderModel.Order, error) {
	orders, err := s.repository.GetOrdersByUser(c, userId)
	if err != nil {
		log.Printf("Failed to get Orders from userID: %v", err)
		return nil, fmt.Errorf("failed to get orders from userID: %v", userId)
	}

	return orders, nil
}

func (s *OrderService) GetOrderById(c context.Context, orderId string, userId string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Printf("User ID mismatch: expected %v, got %v", userId, order.UserID)
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	return order, nil
}

func (s *OrderService) UpdateWeight(c context.Context, orderId string, weight string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	weightFloat, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		log.Printf("Failed to parse weight: %v", err)
		return nil, fmt.Errorf("failed to parse weight: %v", weight)
	}

	order.Weight = &weightFloat

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Printf("Failed to update order weight by ID: %v", err)
		return nil, fmt.Errorf("failed to update order weight by ID: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) AcceptOrder(c context.Context, orderId string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	if order.Status != "created" {
		log.Printf("Order is already in a non-acceptable status: %v", order.Status)
		return nil, fmt.Errorf("order is already in a non-acceptable status: %v", order.Status)
	}

	order.Status = "accepted"

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Printf("Failed to update order status to 'accepted' for order ID %s: %v", orderId, err)
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return order, nil
}

func (s *OrderService) CompleteOrder(c context.Context, orderId string, userId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Printf("User ID mismatch: expected %v, got %v", userId, order.UserID)
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	if order.Status != "delivered" {
		log.Printf("Order cannot be completed due to its current status: %v", order.Status)
		return nil, fmt.Errorf("order cannot be completed due to its current status: %v", order.Status)
	}

	if order.TransactionID == "" {
		log.Printf("Order cannot be completed due to missing transaction ID")
		return nil, fmt.Errorf("order cannot be completed due to missing transaction ID")
	}

	utils.CopyTo(&order, &history)
	history.Status = "completed"
	history.DeletedAt = time.Now()

	if err := s.repository.CreateHistory(c, &history); err != nil {
		log.Printf("Failed to move order to history: %v", err)
		return nil, fmt.Errorf("failed to move order to history: %v", err)
	}

	if err := s.repository.DeleteOrder(c, order); err != nil {
		log.Printf("Failed to delete order by ID %s: %v", orderId, err)
		return nil, fmt.Errorf("failed to delete order by ID: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) PayOrder(c context.Context, orderId string, req *orderRequest.Payment) (*orderModel.Order, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Printf("Failed to validate Order request: %v", err)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	if order.Price == nil {
		log.Printf("Payment is not allowed, invalid price: %v", order.Price)
		return nil, fmt.Errorf("payment is not allowed, invalid price: %v", order.Price)
	}

	order.TransactionID = req.TransactionID

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Printf("Failed to update transaction ID: %v", err)
		return nil, fmt.Errorf("failed to update transaction by id: %v", err)
	}

	return order, nil
}

func (s *OrderService) RejectOrder(c context.Context, orderId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	if order.Status != "created" {
		log.Printf("Order cannot be rejected in its current status: %v", order.Status)
		return nil, fmt.Errorf("order cannot be rejected in its current status: %v", order.Status)
	}

	utils.CopyTo(&order, &history)
	history.Reason = "rejected"
	history.DeletedAt = time.Now()

	if err := s.repository.CreateHistory(c, &history); err != nil {
		log.Printf("Failed to move order to history: %v", err)
		return nil, fmt.Errorf("failed to move order to history: %w", err)
	}

	if err := s.repository.DeleteOrder(c, order); err != nil {
		log.Printf("Failed to delete order by ID %s: %v", orderId, err)
		return nil, fmt.Errorf("failed to delete order by ID %s: %w", orderId, err)
	}

	return order, nil
}

func (s *OrderService) CancelOrder(c context.Context, orderId string, userId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Printf("User ID mismatch")
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	if order.Status != "created" {
		log.Printf("Order cannot be cancelled in its current status: %v", order.Status)
		return nil, fmt.Errorf("order cannot be cancelled in its current status: %v", order.Status)
	}

	utils.CopyTo(&order, &history)
	history.Reason = "cancelled"
	history.DeletedAt = time.Now()

	if err := s.repository.CreateHistory(c, &history); err != nil {
		log.Printf("Failed to move order to history: %v", err)
		return nil, fmt.Errorf("failed to move order to history: %w", err)
	}

	if err := s.repository.DeleteOrder(c, order); err != nil {
		log.Printf("Failed to delete order by ID %s: %v", orderId, err)
		return nil, fmt.Errorf("failed to delete order by ID %s: %w", orderId, err)
	}

	return order, nil
}

func (s *OrderService) EditOrder(c context.Context, orderId string, userId string, req *orderRequest.Order) (*orderModel.Order, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Printf("Validation failed for update profile request: %v", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Printf("Failed to get Order by id: %v", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Printf("User ID mismatch")
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	if order.Status != "created" {
		log.Printf("Editing is not allowed for orders with status: %v", order.Status)
		return nil, fmt.Errorf("editing is not allowed for orders with status: %v", order.Status)
	}

	utils.CopyTo(&req, order)

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Printf("Failed to update order with ID %s: %v", orderId, err)
		return nil, fmt.Errorf("failed to update order with ID %s: %w", orderId, err)
	}

	return order, nil
}
