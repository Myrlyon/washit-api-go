package userRoutes

import (
	"github.com/gin-gonic/gin"

	user "washit-api/app/user/handler"
	"washit-api/app/user/repository"
	"washit-api/app/user/service"
	dbs "washit-api/db"
	"washit-api/middleware"
)

func Main(r *gin.RouterGroup, db dbs.DatabaseInterface) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc)

	authMiddleware := middleware.JWTAuth()
	// refreshAuthMiddleware := middleware.JWTRefresh()
	authRoute := r.Group("/auth")
	profileRoute := r.Group("/profile")

	// Auth
	authRoute.POST("/register", userHandler.Register)
	authRoute.POST("/login", userHandler.Login)

	// Profile
	profileRoute.GET("/me", authMiddleware, userHandler.GetMe)
}
