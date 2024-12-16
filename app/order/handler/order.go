package order

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

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

func (h *OrderHandler) GetOrders(ctx *gin.Context) {
	var res []orderResource.Order

	orders, err := h.service.GetOrders(ctx)
	if err != nil {
		log.Println("Failed to get orders ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&orders, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"data": res})
}
