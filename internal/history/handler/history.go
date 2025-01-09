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

func (h *HistoryHandler) GetHistoryByID(c *gin.Context) {
	var res historyResource.History

	history, err := h.service.GetHistoryByID(c, c.Param("id"), c.GetInt64("userID"))
	if err != nil {
		log.Println("Failed to get history by ID", err)
		response.Error(c, http.StatusInternalServerError, "failed to get history by ID", err)
	}

	utils.CopyTo(&history, &res)
	response.Success(c, http.StatusOK, "successfully retrieved history by ID", &res, nil)
}

func (h *HistoryHandler) GetHistoriesMe(c *gin.Context) {
	var res []historyResource.History

	histories, err := h.service.GetHistoriesMe(c, c.GetInt64("userID"))
	if err != nil {
		log.Println("Failed to get histories me", err)
		response.Error(c, http.StatusInternalServerError, "failed to get histories me", err)
	}

	utils.CopyTo(&histories, &res)
	response.Success(c, http.StatusOK, "succesfully retrieved histories me", &res, nil)
}

func (h *HistoryHandler) GetHistoriesByUser(c *gin.Context) {
	var res []historyResource.History

	userID, err := utils.StringToInt64(c.Param("id"))
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	histories, err := h.service.GetHistoriesByUser(c, userID)
	if err != nil {
		log.Printf("Failed to get histories by user: %v", err)
		response.Error(c, http.StatusInternalServerError, "failed to get histories by user", err)
		return
	}

	utils.CopyTo(&histories, &res)
	response.Success(c, http.StatusOK, "successfully retrieved histories by user", &res, nil)
}
