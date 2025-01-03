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

func Main(r *gin.RouterGroup, db dbs.DatabaseInterface, cache redis.RedisInterface) {
	repository := historyRepository.NewHistoryRepository(db)
	service := historyService.NewHistoryService(repository)
	handler := history.NewHistoryHandler(service, cache)

	authMiddleware := middleware.JWTAuth()
	// adminAuthMiddleware := middleware.JTWAuthAdmin()

	r.GET("/histories", authMiddleware, handler.GetHistoriesMe)
}
