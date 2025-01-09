package history

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	historyRequest "washit-api/internal/history/dto/request"
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
	var userID string

	if c.GetString("userRole") == "admin" {
		userID = ""
	} else {
		userID = c.GetString("userID")
	}

	history, err := h.service.GetHistoryByID(c, c.Param("id"), userID)
	if err != nil {
		log.Println("Failed to get history by ID", err)
		response.Error(c, http.StatusInternalServerError, "failed to get history by ID", err)
	}

	utils.CopyTo(&history, &res)
	response.Success(c, http.StatusOK, "successfully retrieved history by ID", &res, nil)
}

func (h *HistoryHandler) GetHistoriesMe(c *gin.Context) {
	var res historyResource.ListHistory
	var req historyRequest.ListHistory

	userID, err := strconv.ParseInt(c.GetString("userID"), 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	req.UserID = userID

	histories, pagination, err := h.service.GetHistoriesMe(c, &req)
	if err != nil {
		log.Println("Failed to get histories me", err)
		response.Error(c, http.StatusInternalServerError, "failed to get histories me", err)
		return
	}

	utils.CopyTo(&histories, &res.Histories)
	res.Pagination = pagination
	response.Success(c, http.StatusOK, "successfully retrieved histories me", &res, nil)
}

func (h *HistoryHandler) GetHistoriesByUser(c *gin.Context) {
	var res historyResource.ListHistory
	var req historyRequest.ListHistory

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	req.UserID = userID

	histories, pagination, err := h.service.GetHistoriesByUser(c, &req)
	if err != nil {
		log.Printf("Failed to get histories by user: %v", err)
		response.Error(c, http.StatusInternalServerError, "failed to get histories by user", err)
		return
	}

	utils.CopyTo(&histories, &res.Histories)
	res.Pagination = pagination
	response.Success(c, http.StatusOK, "successfully retrieved histories by user", &res, nil)
}

func (h * HistoryHandler) GetAllHistories (c *gin.Context) {
	var res historyResource.ListHistory
	var req historyRequest.ListHistory

	histories, pagination, err := h.service.GetAllHistories(c, &req)
	if err != nil {
		log.Println("Failed to get all histories", err)
		response.Error(c, http.StatusInternalServerError, "failed to get all histories", err)
		return
	}

	utils.CopyTo(&histories, &res.Histories)
	res.Pagination = pagination
	response.Success(c, http.StatusOK, "successfully retrieved all histories", &res, nil)
}
