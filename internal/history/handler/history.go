package history

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	historyResource "washit-api/internal/history/dto/resource"
	historyService "washit-api/internal/history/service"
	"washit-api/pkg/redis"
	"washit-api/pkg/response"
	"washit-api/pkg/utils"
)

type HistoryHandler struct {
	service historyService.IHistoryService
	cache   redis.IRedis
}

func NewHistoryHandler(service historyService.IHistoryService, cache redis.IRedis) *HistoryHandler {
	return &HistoryHandler{
		service: service,
		cache:   cache,
	}
}

func (h *HistoryHandler) GetHistoriesMe(c *gin.Context) {
	var res []historyResource.History

	history, err := h.service.GetHistoriesMe(c, c.GetString("userId"))
	if err != nil {
		log.Println("Failed to get histories me", err)
		response.Error(c, http.StatusInternalServerError, "failed to get histories me", err)
	}

	utils.CopyTo(&history, &res)
	response.Success(c, http.StatusOK, "succesfully retrieved histories me", &res, nil)
}
