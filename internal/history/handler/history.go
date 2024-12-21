package history

import (
	historyService "washit-api/internal/history/service"
	"washit-api/pkg/redis"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	service historyService.HistoryServiceInterface
	cache   redis.RedisInterface
}

func NewHistoryHandler(service historyService.HistoryServiceInterface, cache redis.RedisInterface) *HistoryHandler {
	return &HistoryHandler{
		service: service,
		cache:   cache,
	}
}

func (h *HistoryHandler) GetHistoriesMe(ctx *gin.Context) {}
