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
	cache   redis.RedisInterface
}

func NewOrderHandler(service orderService.IOrderService, cache redis.RedisInterface) *OrderHandler {
	return &OrderHandler{
		service: service,
		cache:   cache,
	}
}

var ordersCacheKey string = "/api/v1/orders"

// @Summary	Create Order
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Param		_	body		orderRequest.Order	true	"Body"
// @Success	201	{object}	orderResource.Order
// @Router		/order [post]
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

// @Summary	Cancel Order
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Param		id	path		string	true	"Order ID"
// @Success	200	{object}	orderResource.Order
// @Router		/order/{id}/cancel [put]
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

// @Summary	Get Order By ID
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Param		id	path		string	true	"Order ID"
// @Success	200	{object}	orderResource.Order
// @Router		/order/{id} [get]
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

// @Summary	Get Orders Me
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	orderResource.OrderList
// @Router		/orders [get]
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

// @Summary	Get Orders All
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	orderResource.OrderList
// @Router		/orders/all [get]
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

// @Summary	Get Orders By User
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	orderResource.OrderList
// @Router		/orders/user/{id} [get]
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

var links = func(orderId string) map[string]response.HypermediaLink {
	return map[string]response.HypermediaLink{
		"self": {
			Href:   "/order/" + orderId,
			Method: "GET",
		},
	}
}
