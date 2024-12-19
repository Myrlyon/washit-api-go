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
	authRefreshMiddleware := middleware.JWTRefresh()

	r.POST("/auth/refresh", authRefreshMiddleware, handler.RefreshToken)

	// Auth
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)
	// r.POST("/auth/google", handler.GoogleAuth)
	r.POST("/auth/logout", authMiddleware, handler.Logout)

	// Profile Get
	r.GET("/profile/me", authMiddleware, handler.GetMe)

	// Profile Put
	r.PUT("/profile/update", authMiddleware, handler.UpdateMe)
	// r.PUT("/profile/picture", authMiddleware, handler.UpdatePicture)

	// Admin Authority

	// Users
	r.GET("/users", adminAuthMiddleware, handler.GetUsers)
	r.GET("/user/:id", adminAuthMiddleware, handler.GetUserById)
	// r.POST)"/user/:id/ban", adminAuthMiddleware, handler.BanUser)
	// r.POST)"/user/:id/unban", adminAuthMiddleware, handler.UnbanUser)
	// r.DELETE("/user/:id", adminAuthMiddleware, handler.DeleteUser)
}
