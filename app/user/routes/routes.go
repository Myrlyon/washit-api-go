package userRoutes

import (
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	user "washit-api/app/user/handler"
	userRepository "washit-api/app/user/repository"
	userService "washit-api/app/user/service"
	dbs "washit-api/db"
	"washit-api/middleware"
	"washit-api/redis"
)

func Main(r *gin.RouterGroup, db dbs.DatabaseInterface, cache redis.RedisInterface, app *firebase.App, validator *validator.Validate) {
	repository := userRepository.NewUserRepository(db)
	service := userService.NewUserService(repository)
	handler := user.NewUserHandler(service, cache, app, validator)

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JTWAuthAdmin()
	authRefreshMiddleware := middleware.JWTRefresh()

	r.POST("/auth/refresh", authRefreshMiddleware, handler.RefreshToken)

	// Auth
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)
	r.POST("/auth/logout", authMiddleware, handler.Logout)
	r.GET("/auth/google", handler.Login)
	r.POST("/auth/google/callback", handler.Login)

	// Profile Get
	r.GET("/profile/me", authMiddleware, handler.GetMe)

	// Profile Put
	r.PUT("/profile/update", authMiddleware, handler.UpdateMe)
	// r.PUT("/profile/picture", authMiddleware, handler.UpdatePicture)

	// Admin Authority

	// Users
	r.GET("/users", adminAuthMiddleware, handler.GetUsers)
	r.GET("/users/banned", adminAuthMiddleware, handler.GetBannedUsers)
	r.GET("/user/:id", adminAuthMiddleware, handler.GetUserById)
	r.PUT("/user/:id/ban", adminAuthMiddleware, handler.BanUser)
	r.PUT("/user/:id/unban", adminAuthMiddleware, handler.UnbanUser)
	// r.DELETE("/user/:id", adminAuthMiddleware, handler.DeleteUser)
}
