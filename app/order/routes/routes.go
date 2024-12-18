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

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JTWAuthAdmin()

	// Order Get
	r.GET("/orders", authMiddleware, handler.GetOrdersMe)
	r.GET("/order/:id", authMiddleware, handler.GetOrderById)

	// Order Post
	r.POST("/order", authMiddleware, handler.CreateOrder)
	// r.POST("/order/:id/cancel", authMiddleware, handler.CancelOrder)

	// Order Update
	// r.PUT("/order/:id/update", authMiddleware, handler.UpdapteOrder)

	// Admin Authority

	// Order Get
	r.GET("/orders/all", adminAuthMiddleware, handler.GetOrdersAll)
	r.GET("/orders/user/:id", adminAuthMiddleware, handler.GetOrdersUser)

	// Order Update
	// r.PUT("/order/:id/accept", adminAuthMiddleware, handler.AcceptOrder)
	// r.PUT("/order/:id/reject", adminAuthMiddleware, handler.RejectOrder)
	// r.PUT("/order/:id/complete", adminAuthMiddleware, handler.CompleteOrder)
	// r.PUT("/order/:id/status", adminAuthMiddleware, handler.UpdateOrderStatus)
	// r.PUT("/order/:id/weight", adminAuthMiddleware, handler.UpdateOrderWeight)
}
