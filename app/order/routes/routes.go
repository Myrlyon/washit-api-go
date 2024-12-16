package orderRoutes

import (
	"github.com/gin-gonic/gin"

	order "washit-api/app/order/handler"
	orderRepository "washit-api/app/order/repository"
	orderService "washit-api/app/order/service"
	dbs "washit-api/db"
	"washit-api/middleware"
	"washit-api/redis"
)

func Main(r *gin.RouterGroup, db dbs.DatabaseInterface, redis redis.RedisInterface) {
	repository := orderRepository.NewOrderRepository(db)
	service := orderService.NewOrderService(repository)
	handler := order.NewOrderHandler(service, redis)

	// authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JTWAuthAdmin()

	r.GET("/orders", adminAuthMiddleware, handler.GetOrders)
}
