package order

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	orderRequest "washit-api/internal/order/dto/request"
	orderResource "washit-api/internal/order/dto/resource"
	orderService "washit-api/internal/order/service"
	"washit-api/pkg/configs"
	"washit-api/pkg/redis"
	"washit-api/pkg/response"
	"washit-api/pkg/utils"
)

type OrderHandler struct {
	service orderService.IOrderService
	cache   redis.IRedis
}

func NewOrderHandler(service orderService.IOrderService, cache redis.IRedis) *OrderHandler {
	return &OrderHandler{
		service: service,
		cache:   cache,
	}
}

var ordersCacheKey string = "/api/v1/orders"

// CreateOrder handles the creation of a new order.
//
//	@Summary	Create a new order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		orderRequest.Order	true	"Order details"
//	@Success	201	{object}	orderResource.Order
//	@Router		/order [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req orderRequest.Order
	var res orderResource.Order

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request body ", err)
		response.Error(c, http.StatusBadRequest, "failed to parse request body", err)
		return
	}

	userId, err := strconv.Atoi(c.GetString("userId"))
	if err != nil {
		log.Println("Failed to get user id ", err)
		response.Error(c, http.StatusBadRequest, "failed to get user id", err)
		return
	}

	order, err := h.service.CreateOrder(c, userId, &req)
	if err != nil {
		log.Println("Failed to create order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to create order", err)
		return
	}

	utils.CopyTo(order, &res)
	response.Success(c, http.StatusCreated, "order is created successfully", &res, links(res.ID))

	_ = h.cache.Remove(ordersCacheKey)
}

// CancelOrder handles the cancellation of an existing order.
//
//	@Summary	Cancel an existing order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	orderResource.Order
//	@Router		/order/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	var res orderResource.Order

	order, err := h.service.CancelOrder(c, c.Param("id"), c.GetString("userId"))
	if err != nil {
		log.Println("Failed to get order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to get order", err)
		return
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "order is cancelled successfully", &res, links(res.ID))

	_ = h.cache.Remove(ordersCacheKey)
}

// GetOrderById retrieves an order by its ID.
//
//	@Summary	Get order details by ID
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	orderResource.Order
//	@Router		/order/{id} [get]
func (h *OrderHandler) GetOrderById(c *gin.Context) {
	var res orderResource.Order
	var userId string

	if c.GetString("userRole") == "admin" {
		userId = "0"
	} else {
		userId = c.GetString("userId")
	}

	order, err := h.service.GetOrderById(c, c.Param("id"), userId)
	if err != nil {
		log.Println("Failed to get order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to get order", err)
		return
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "order is collected successfully", &res, links(res.ID))
}

// GetOrdersMe retrieves all orders for the authenticated user.
//
//	@Summary	Get all orders for the authenticated user
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	orderResource.Order
//	@Router		/orders [get]
func (h *OrderHandler) GetOrdersMe(c *gin.Context) {
	var res []orderResource.Order

	if err := h.cache.Get(ordersCacheKey, &res); err == nil {
		log.Println("Failed to get orders ", err)
		response.Success(c, http.StatusOK, "orders are collected successfully", &res, nil)
		return
	}

	userId := c.GetString("userId")
	orders, err := h.service.GetOrdersMe(c, userId)
	if err != nil {
		log.Println("Failed to get orders ", err)
		response.Error(c, http.StatusInternalServerError, "failed to get orders", err)
		return
	}

	utils.CopyTo(&orders, &res)
	response.Success(c, http.StatusOK, "orders are collected successfully", &res, nil)

	_ = h.cache.SetWithExpiration(ordersCacheKey, &res, configs.ProductCachingTime)
}

// GetOrdersAll retrieves all orders.
//
//	@Summary	Get all orders
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	orderResource.Order
//	@Router		/orders/all [get]
func (h *OrderHandler) GetOrdersAll(c *gin.Context) {
	var res []orderResource.Order

	err := h.cache.Get(ordersCacheKey, &res)
	if err == nil {
		log.Println("Failed to get orders ", err)
		response.Success(c, http.StatusOK, "orders are collected successfully", &res, nil)
		return
	}

	orders, err := h.service.GetOrdersAll(c)
	if err != nil {
		log.Println("Failed to get orders ", err)
		response.Error(c, http.StatusInternalServerError, "failed to get orders", err)
		return
	}

	utils.CopyTo(&orders, &res)
	response.Success(c, http.StatusOK, "orders are collected successfully", &res, nil)

	_ = h.cache.SetWithExpiration(ordersCacheKey, &res, configs.ProductCachingTime)
}

// GetOrdersByUser retrieves all orders for a specific user.
//
//	@Summary	Get all orders for a specific user
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	orderResource.Order
//	@Router		/orders/user/{id} [get]
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	var res []orderResource.Order

	orders, err := h.service.GetOrdersByUser(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to get orders ", err)
		response.Error(c, http.StatusInternalServerError, "failed to get orders", err)
		return
	}

	utils.CopyTo(&orders, &res)
	response.Success(c, http.StatusOK, "orders are collected successfully", &res, nil)
}

// EditOrder handles the editing of an existing order.
//
//	@Summary	Edit an existing order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string				true	"Order ID"
//	@Param		_	body		orderRequest.Order	true	"Order details"
//	@Success	200	{object}	orderResource.Order
//	@Router		/order/{id} [put]
func (h *OrderHandler) EditOrder(c *gin.Context) {
	var req orderRequest.Order
	var res orderResource.Order

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		response.Error(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	order, err := h.service.EditOrder(c, c.Param("id"), c.GetString("userId"), &req)
	if err != nil {
		log.Println("Failed to update order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to update order", err)
		return
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "order is updated successfully", &res, links(res.ID))

	_ = h.cache.SetWithExpiration(ordersCacheKey, &res, configs.ProductCachingTime)
}

// AcceptOrder handles the acceptance of an order.
//
//	@Summary	Accept an order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	orderResource.Order
//	@Router		/order/{id}/accept [put]
func (h *OrderHandler) AcceptOrder(c *gin.Context) {
	var res orderResource.Order

	order, err := h.service.AcceptOrder(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to accept order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to accept order", err)
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "order is accepted successfully", &res, links(res.ID))
}

// CompleteOrder handles the completion of an order.
//
//	@Summary	Complete an order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	orderResource.Order
//	@Router		/order/{id}/complete [put]
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	var res orderResource.Order

	order, err := h.service.CompleteOrder(c, c.Param("id"), c.GetString("userId"))
	if err != nil {
		log.Println("Failed to complete order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to complete order", err)
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "order is completed successfully", &res, links(res.ID))
}

// RejectOrder handles the rejection of an order.
//
//	@Summary	Reject an order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	orderResource.Order
//	@Router		/order/{id}/reject [put]
func (h *OrderHandler) RejectOrder(c *gin.Context) {
	var res orderResource.Order

	order, err := h.service.RejectOrder(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to reject order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to reject order", err)
		return
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "order is rejected successfully", &res, links(res.ID))
}

// UpdateWeight handles the updating of an order's weight.
//
//	@Summary	Update the weight of an order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id		path		string	true	"Order ID"
//	@Param		weight	path		string	true	"Weight"
//	@Success	200		{object}	orderResource.Order
//	@Router		/order/{id}/weight/{weight} [put]
func (h *OrderHandler) UpdateWeight(c *gin.Context) {
	var res orderResource.Order

	if _, err := strconv.ParseFloat(c.Param("weight"), 64); err != nil {
		response.Error(c, http.StatusBadRequest, "weight must be a number", err)
		return
	}

	order, err := h.service.UpdateWeight(c, c.Param("id"), c.Param("weight"))
	if err != nil {
		log.Println("Failed to update weight ", err)
		response.Error(c, http.StatusInternalServerError, "failed to update weight", err)
		return
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "weight is updated successfully", &res, links(res.ID))
}

// PayOrder handles the payment of an order.
//
//	@Summary	Pay for an order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string					true	"Order ID"
//	@Param		_	body		orderRequest.Payment	true	"Payment details"
//	@Success	200	{object}	orderResource.Order
//	@Router		/order/{id}/pay [put]
func (h *OrderHandler) PayOrder(c *gin.Context) {
	var res orderResource.Order
	var req orderRequest.Payment

	order, err := h.service.PayOrder(c, c.Param("id"), &req)
	if err != nil {
		log.Println("Failed to pay order ", err)
		response.Error(c, http.StatusInternalServerError, "failed to pay order", err)
	}

	utils.CopyTo(&order, &res)
	response.Success(c, http.StatusOK, "order is paid successfully", &res, links(res.ID))
}

var links = func(orderId string) map[string]response.HypermediaLink {
	return map[string]response.HypermediaLink{
		"self": {
			Href:   "/order/" + orderId,
			Method: "GET",
		},
	}
}
