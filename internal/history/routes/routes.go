package historyRoutes

import (
	history "washit-api/internal/history/handler"
	historyRepository "washit-api/internal/history/repository"
	historyService "washit-api/internal/history/service"
	"washit-api/pkg/db/dbs"
	"washit-api/pkg/middleware"
	"washit-api/pkg/redis"

	"github.com/gin-gonic/gin"
)

func Main(r *gin.RouterGroup, db dbs.IDatabase, cache redis.IRedis) {
	repository := historyRepository.NewHistoryRepository(db)
	service := historyService.NewHistoryService(repository)
	handler := history.NewHistoryHandler(service, cache)

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JWTAuthAdmin()

	r.GET("/histories", authMiddleware, handler.GetHistoriesMe)
	r.GET("/history/:id", authMiddleware, handler.GetHistoryByID)

	//ADMIN
	r.GET("/histories/user/:id", adminAuthMiddleware, handler.GetHistoriesByUser)
}
