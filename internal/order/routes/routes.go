package orderRoutes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	order "washit-api/internal/order/handler"
	orderRepository "washit-api/internal/order/repository"
	orderService "washit-api/internal/order/service"
	"washit-api/pkg/db/dbs"
	"washit-api/pkg/middleware"
	"washit-api/pkg/redis"
)

func Main(r *gin.RouterGroup, db dbs.IDatabase, cache redis.IRedis, validator *validator.Validate) {
	repository := orderRepository.NewOrderRepository(db)
	service := orderService.NewOrderService(repository, validator)
	handler := order.NewOrderHandler(service, cache)

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JTWAuthAdmin()

	// Order Get
	r.GET("/orders", authMiddleware, handler.GetOrdersMe)
	r.GET("/order/:id", authMiddleware, handler.GetOrderById)

	// Order Post
	r.POST("/order", authMiddleware, handler.CreateOrder)

	// Order Update
	r.PUT("/order/:id/edit", authMiddleware, handler.EditOrder)
	r.PUT("/order/:id/cancel", authMiddleware, handler.CancelOrder)
	r.PUT("/order/:id/complete", authMiddleware, handler.CompleteOrder)
	r.PUT("/order/:id/pay", authMiddleware, handler.PayOrder)

	// Admin Authority

	// Order Get
	r.GET("/orders/all", adminAuthMiddleware, handler.GetOrdersAll)
	r.GET("/orders/user/:id", adminAuthMiddleware, handler.GetOrdersByUser)

	// Order Update
	// r.PUT("/order/:id/update", )
	r.PUT("/order/:id/accept", adminAuthMiddleware, handler.AcceptOrder)
	r.PUT("/order/:id/reject", adminAuthMiddleware, handler.RejectOrder)
	// r.PUT("/order/:id/status", adminAuthMiddleware, handler.UpdateOrderStatus)
	r.PUT("/order/:id/weight/:weight", adminAuthMiddleware, handler.UpdateWeight)
}
