package order

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	orderRequest "washit-api/internal/order/dto/request"
	orderResource "washit-api/internal/order/dto/resource"
	orderService "washit-api/internal/order/service"
	"washit-api/pkg/configs"
	"washit-api/pkg/redis"
	"washit-api/pkg/utils"
)

type OrderHandler struct {
	service   orderService.OrderServiceInterface
	cache     redis.RedisInterface
	validator *validator.Validate
}

func NewOrderHandler(service orderService.OrderServiceInterface, cache redis.RedisInterface, validator *validator.Validate) *OrderHandler {
	return &OrderHandler{
		service:   service,
		cache:     cache,
		validator: validator,
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
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	userId, err := strconv.Atoi(c.GetString("userId"))
	if err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	order, err := h.service.CreateOrder(c, userId, &req)
	if err != nil {
		log.Println("Failed to create order ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(order, &res.Order)
	res.Hypermedia = links(order.ID)
	res.Message = "order is created succesfully"
	utils.WriteJson(c, http.StatusCreated, &res)

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
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&order, &res.Order)
	res.Hypermedia = links(order.ID)
	res.Message = "order is cancelled succesfully"
	utils.WriteJson(c, http.StatusOK, &res)

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
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&order, &res.Order)
	res.Hypermedia = links(order.ID)
	res.Message = "order is collected successfully"
	utils.WriteJson(c, http.StatusOK, &res)
}

// @Summary	Get Orders Me
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	orderResource.OrderList
// @Router		/orders [get]
func (h *OrderHandler) GetOrdersMe(c *gin.Context) {
	var res orderResource.OrderList

	cacheKey := c.Request.URL.RequestURI()
	log.Println("Cache key: ", cacheKey)
	if err := h.cache.Get(cacheKey, &res); err == nil {
		utils.WriteJson(c, http.StatusOK, &res)
		return
	}

	userId := c.GetString("userId")
	orders, err := h.service.GetOrders(c, userId)
	if err != nil {
		log.Println("Failed to get orders ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&orders, &res.Orders)
	res.Message = "orders are collected successfully"
	utils.WriteJson(c, http.StatusOK, &res)

	_ = h.cache.SetWithExpiration(cacheKey, &res, configs.ProductCachingTime)
}

// @Summary	Get Orders All
// @Tags		Order
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	orderResource.OrderList
// @Router		/orders/all [get]
func (h *OrderHandler) GetOrdersAll(c *gin.Context) {
	var res orderResource.OrderList

	cacheKey := c.Request.URL.RequestURI()
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		utils.WriteJson(c, http.StatusOK, &res)
		return
	}

	orders, err := h.service.GetOrdersAll(c)
	if err != nil {
		log.Println("Failed to get orders ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&orders, &res.Orders)
	res.Message = "orders are collected successfully"
	utils.WriteJson(c, http.StatusOK, &res)

	_ = h.cache.SetWithExpiration(cacheKey, &res, configs.ProductCachingTime)
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
	var res orderResource.OrderList

	orders, err := h.service.GetOrdersByUser(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to get orders ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&orders, &res.Orders)
	res.Message = "orders are collected successfully"
	utils.WriteJson(c, http.StatusOK, &res)
}

var links = func(orderID string) orderResource.Hypermedia {
	return orderResource.Hypermedia{
		Self:   map[string]string{"href": "/api/v1/order/" + orderID, "method": "GET"},
		Create: map[string]string{"href": "/api/v1/order", "method": "POST"},
		Cancel: map[string]string{"href": "/api/v1/order/" + orderID + "/cancel", "method": "PUT"},
	}
}
