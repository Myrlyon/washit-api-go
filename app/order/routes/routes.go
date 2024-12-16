package orderRoutes

import (
	"github.com/gin-gonic/gin"

	dbs "washit-api/db"
)

func Main(r *gin.RouterGroup, db dbs.DatabaseInterface) {
	// userRepo := repository.NewUserRepository(db)
	// userSvc := service.NewUserService(userRepo)
	// userHandler := user.NewUserHandler(userSvc)

	// authMiddleware := middleware.JWTAuth()
	// refreshAuthMiddleware := middleware.JWTRefresh()
	// profileRoute := r.Group("/profile")

	// Profile
	// r.GET("/orders", authMiddleware, userHandler.GetMe)
	// r.GET("/order", authMiddleware, userHandler.GetMe)
	// r.POST("/order/new", authMiddleware, userHandler.GetMe)
	// r.POST("/order/update", authMiddleware, userHandler.GetMe)
	// r.POST("/order/cancel", authMiddleware, userHandler.GetMe)
}
