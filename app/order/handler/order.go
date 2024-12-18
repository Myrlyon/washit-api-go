package order

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	orderRequest "washit-api/app/order/dto/request"
	orderResource "washit-api/app/order/dto/resource"
	orderService "washit-api/app/order/service"
	"washit-api/redis"
	"washit-api/utils"
)

type OrderHandler struct {
	service orderService.OrderServiceInterface
	cache   redis.RedisInterface
}

func NewOrderHandler(service orderService.OrderServiceInterface, cache redis.RedisInterface) *OrderHandler {
	return &OrderHandler{
		service: service,
		cache:   cache,
	}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req orderRequest.Order
	var res orderResource.Order

	if err := utils.ParseJson(ctx, &req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	userId, err := strconv.Atoi(ctx.GetString("userId"))
	if err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	order, err := h.service.CreateOrder(ctx, userId, &req)
	if err != nil {
		log.Println("Failed to create order ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&order, &res)
	utils.WriteJson(ctx, http.StatusCreated, map[string]interface{}{"order": res})
}

func (h *OrderHandler) GetOrderById(ctx *gin.Context) {
	var res orderResource.Order
	var userId string

	if ctx.GetString("userRole") == "admin" {
		userId = "0"
	} else {
		userId = ctx.GetString("userId")
	}

	order, err := h.service.GetOrderById(ctx, ctx.Param("id"), userId)
	if err != nil {
		log.Println("Failed to get order ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&order, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"order": res})
}

func (h *OrderHandler) GetOrdersMe(ctx *gin.Context) {
	var res []orderResource.Order

	userId := ctx.GetString("userId")
	orders, err := h.service.GetOrders(ctx, userId)
	if err != nil {
		log.Println("Failed to get orders ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&orders, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"orders": res})
}

func (h *OrderHandler) GetOrdersAll(ctx *gin.Context) {
	var res []orderResource.Order

	orders, err := h.service.GetOrders(ctx, "0")
	if err != nil {
		log.Println("Failed to get orders ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&orders, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"orders": res})
}

func (h *OrderHandler) GetOrdersUser(ctx *gin.Context) {
	var res []orderResource.Order

	orders, err := h.service.GetOrders(ctx, ctx.Param("id"))
	if err != nil {
		log.Println("Failed to get orders ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&orders, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"orders": res})
}
