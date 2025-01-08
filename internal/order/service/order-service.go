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
	order.Status = "created"

	order, err = s.repository.CreateOrder(c, order)
	if err != nil {
		log.Println("Failed to create Order ", err)
		return nil, fmt.Errorf("failed to create order")
	}

	return order, nil
}

func (s *OrderService) GetOrdersMe(c context.Context, userId string) ([]*orderModel.Order, error) {
	order, err := s.repository.GetOrders(c, userId)
	if err != nil {
		log.Println("Failed to get Orders me")
		return nil, fmt.Errorf("failed to get orders me")
	}

	return order, nil
}

func (s *OrderService) GetOrdersAll(c context.Context) ([]*orderModel.Order, error) {
	order, err := s.repository.GetOrders(c, "")
	if err != nil {
		log.Println("Failed to get all Orders ", err)
		return nil, fmt.Errorf("failed to get all orders")
	}

	return order, nil
}

func (s *OrderService) GetOrdersByUser(c context.Context, userId string) ([]*orderModel.Order, error) {
	order, err := s.repository.GetOrders(c, userId)
	if err != nil {
		log.Println("Failed to get Orders from userID ", err)
		return nil, fmt.Errorf("failed to get orders from userid: %v", userId)
	}

	return order, nil
}

func (s *OrderService) GetOrderById(c context.Context, orderId string, userId string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Println("User ID mismatch")
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	return order, nil
}

func (s *OrderService) UpdateWeight(c context.Context, orderId string, weight string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	weightFloat, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		log.Println("Failed to parse weight", err)
		return nil, fmt.Errorf("failed to parse weight: %v", weight)
	}

	order.Weight = &weightFloat

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Println("Failed to update order weight by ID:", err)
		return nil, fmt.Errorf("failed to update order weight by ID: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) AcceptOrder(c context.Context, orderId string) (*orderModel.Order, error) {
	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if order.Status != "created" {
		log.Printf("Order is already in a non-acceptable status: %v", order.Status)
		return nil, fmt.Errorf("order is already in a non-acceptable status: %v", order.Status)
	}

	order.Status = "accepted"

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Println("Failed to update status")
		return nil, fmt.Errorf("failed to update status: %v", err)
	}

	return order, nil
}

func (s *OrderService) CompleteOrder(c context.Context, orderId string, userId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Panicln("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Println("User ID mismatch")
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	if order.Status != "created" {
		log.Printf("Order is already in a non-acceptable status: %v", order.Status)
		return nil, fmt.Errorf("order is already in a non-acceptable status: %v", order.Status)
	}

	utils.CopyTo(&order, &history)
	history.Status = "completed"

	if err := s.repository.CreateHistory(c, &history); err != nil {
		log.Println("Failed to move order to history", err)
		return nil, fmt.Errorf("failed to move order to history: %v", err)
	}

	if err := s.repository.DeleteOrder(c, order); err != nil {
		log.Println("Failed to delete order by ID:", err)
		return nil, fmt.Errorf("failed to delete order by ID: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) PayOrder(c context.Context, orderId string, req *orderRequest.Payment) (*orderModel.Order, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Println("Failed to validate Order request ", err)
		return nil, err
	}

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	order.TransactionID = req.TransactionID

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Println("Failed to update transaction Id ", err)
		return nil, fmt.Errorf("failed to update transaction by id: %v", err)
	}

	return order, nil
}

func (s *OrderService) RejectOrder(c context.Context, orderId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if order.Status != "created" {
		log.Println("Rejecting is forbidden")
		return nil, fmt.Errorf("Rejecting is forbidden in status: %v", order.Status)
	}

	utils.CopyTo(&order, &history)
	history.Reason = "rejected"
	history.DeletedAt = time.Now()

	if err := s.repository.CreateHistory(c, &history); err != nil {
		log.Println("Failed to move order to history", err)
		return nil, fmt.Errorf("failed to move order to history: %v", err)
	}

	if err := s.repository.DeleteOrder(c, order); err != nil {
		log.Println("Failed to delete order by ID:", err)
		return nil, fmt.Errorf("failed to delete order by ID: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) CancelOrder(c context.Context, orderId string, userId string) (*orderModel.Order, error) {
	var history historyModel.History

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Println("User ID mismatch")
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	if order.Status != "created" {
		log.Println("Cancelling is forbidden")
		return nil, fmt.Errorf("Cancelling is forbidden in status: %v", order.Status)
	}

	utils.CopyTo(&order, &history)
	history.Reason = "cancelled"
	history.DeletedAt = time.Now()

	if err := s.repository.CreateHistory(c, &history); err != nil {
		log.Println("Failed to move order to history", err)
		return nil, fmt.Errorf("failed to move order to history: %v", err)
	}

	if err := s.repository.DeleteOrder(c, order); err != nil {
		log.Println("Failed to delete order by ID:", err)
		return nil, fmt.Errorf("failed to delete order by ID: %v", orderId)
	}

	return order, nil
}

func (s *OrderService) EditOrder(c context.Context, orderId string, userId string, req *orderRequest.Order) (*orderModel.Order, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Println("Failed to validate update profile request ", err)
		return nil, err
	}

	order, err := s.repository.GetOrderById(c, orderId)
	if err != nil {
		log.Println("Failed to get Order by id", err)
		return nil, fmt.Errorf("failed to get order by id: %v", err)
	}

	if strconv.Itoa(order.UserID) != userId {
		log.Println("User ID mismatch")
		return nil, fmt.Errorf("user ID mismatch: %v", userId)
	}

	if order.Status != "created" {
		log.Println("Editing is forbidden")
		return nil, fmt.Errorf("Editing is forbidden in status: %v", order.Status)
	}

	utils.CopyTo(&req, order)

	if err := s.repository.UpdateOrder(c, order); err != nil {
		log.Println("Failed to update order by id ", err)
		return nil, fmt.Errorf("failed to update order by id: %v", orderId)
	}

	return order, nil
}
