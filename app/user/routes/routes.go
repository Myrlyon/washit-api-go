package userRoutes

import (
	"github.com/gin-gonic/gin"

	user "washit-api/app/user/handler"
	userRepository "washit-api/app/user/repository"
	userService "washit-api/app/user/service"
	dbs "washit-api/db"
	"washit-api/middleware"
	"washit-api/redis"
)

func Main(r *gin.RouterGroup, db dbs.DatabaseInterface, cache redis.RedisInterface) {
	repository := userRepository.NewUserRepository(db)
	service := userService.NewUserService(repository)
	handler := user.NewUserHandler(service, cache)

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JTWAuthAdmin()

	// Auth
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)

	// Profile
	r.GET("/profile/me", authMiddleware, handler.GetMe)

	// User
	r.GET("/users", adminAuthMiddleware, handler.GetUsers)
}
