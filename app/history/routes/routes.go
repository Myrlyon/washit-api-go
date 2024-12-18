package historyRoutes

import (
	history "washit-api/app/history/handler"
	historyRepository "washit-api/app/history/repository"
	historyService "washit-api/app/history/service"
	dbs "washit-api/db"
	"washit-api/middleware"
	"washit-api/redis"

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
