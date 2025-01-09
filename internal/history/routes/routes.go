package historyRoutes

import (
	history "washit-api/internal/history/handler"
	historyRepository "washit-api/internal/history/repository"
	historyService "washit-api/internal/history/service"
	"washit-api/pkg/db/dbs"
	"washit-api/pkg/middleware"
	"washit-api/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func Main(r *gin.RouterGroup, db dbs.IDatabase, cache redis.IRedis, validator *validator.Validate) {
	repository := historyRepository.NewHistoryRepository(db)
	service := historyService.NewHistoryService(repository, validator)
	handler := history.NewHistoryHandler(service, cache)

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JWTAuthAdmin()

	r.GET("/histories/me", authMiddleware, handler.GetHistoriesMe)
	r.GET("/history/:id", authMiddleware, handler.GetHistoryByID)

	//ADMIN
	r.GET("/histories/user/:id", adminAuthMiddleware, handler.GetHistoriesByUser)
	r.GET("/histories/all", adminAuthMiddleware, handler.GetAllHistories)
}
