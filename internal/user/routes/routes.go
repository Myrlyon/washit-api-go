package userRoutes

import (
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	user "washit-api/internal/user/handler"
	userRepository "washit-api/internal/user/repository"
	userService "washit-api/internal/user/service"
	"washit-api/pkg/db/dbs"
	"washit-api/pkg/middleware"
	"washit-api/pkg/redis"
)

func Main(r *gin.RouterGroup, db dbs.IDatabase, cache redis.IRedis, app *firebase.App, validator *validator.Validate) {
	repository := userRepository.NewUserRepository(db)
	service := userService.NewUserService(repository, validator)
	handler := user.NewUserHandler(service, cache, app)

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JWTAuthAdmin()
	authRefreshMiddleware := middleware.JWTRefresh()

	r.POST("/auth/refresh", authRefreshMiddleware, handler.RefreshToken)

	// Auth
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)
	r.POST("/auth/logout", authMiddleware, handler.Logout)
	r.POST("/auth/google", handler.LoginWithGoogle)
	// r.POST("/auth/google/callback", handler.Login)

	// Profile Get
	r.GET("/profile/me", authMiddleware, handler.GetMe)

	// Profile Put
	r.PUT("/profile/update", authMiddleware, handler.UpdateMe)
	r.PUT("/profile/update/password", authMiddleware, handler.UpdatePassword)
	r.PUT("/profile/update/picture", authMiddleware, handler.UpdatePicture)

	// Admin Authority

	// Users
	r.GET("/users", adminAuthMiddleware, handler.GetUsers)
	r.GET("/users/banned", adminAuthMiddleware, handler.GetBannedUsers)
	r.GET("/user/:id", adminAuthMiddleware, handler.GetUserByID)
	r.PUT("/user/:id/ban", adminAuthMiddleware, handler.BanUser)
	r.PUT("/user/:id/unban", adminAuthMiddleware, handler.UnbanUser)
}
