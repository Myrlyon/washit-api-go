package historyService

import (
	"fmt"
	"log"
	"strconv"
	historyModel "washit-api/internal/history/dto/model"
	historyRequest "washit-api/internal/history/dto/request"
	historyRepository "washit-api/internal/history/repository"
	"washit-api/pkg/paging"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type IHistoryService interface {
	GetHistoryByID(c *gin.Context, historyID string, userID string) (*historyModel.History, error)
	GetHistoriesMe(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History, *paging.Pagination, error)
	GetHistoriesByUser(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History,*paging.Pagination, error)
	GetAllHistories(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History,*paging.Pagination, error)
}

type HistoryService struct {
	repository historyRepository.IHistoryRepository
	validator *validator.Validate
}

func NewHistoryService(repository historyRepository.IHistoryRepository, validator *validator.Validate) *HistoryService {
	return &HistoryService{
		repository: repository,
		validator: validator,
	}
}

func (s *HistoryService) GetHistoryByID(c *gin.Context, historyID string, userID string) (*historyModel.History, error) {
	history, err := s.repository.GetHistoryByID(c, historyID)
	if err != nil {
		log.Printf("Failed to get history by ID: %v", err)
		return nil, fmt.Errorf("failed to get history by ID: %v", err)
	}

	if strconv.FormatInt(history.UserID, 10) != userID && userID != "" {
		log.Printf("User ID mismatch: expected %v, got %v", userID, history.UserID)
		return nil, fmt.Errorf("user ID mismatch: %v", userID)
	}

	return history, err
}

func (s *HistoryService) GetHistoriesMe(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History,*paging.Pagination, error) {
	histories, pagination, err := s.repository.GetHistories(c, req)
	if err != nil {
		log.Printf("Failed to get histories by user id: %v", err)
		return nil, nil, fmt.Errorf("failed to get histories by user id: %v", err)
	}

	return histories, pagination, nil
}

func (s *HistoryService) GetHistoriesByUser(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History,*paging.Pagination, error) {
	histories, pagination, err := s.repository.GetHistories(c, req)
	if err != nil {
		log.Printf("Failed to get histories by user id: %v", err)
		return nil, nil, fmt.Errorf("failed to get histories by user id: %v", err)
	}

	return histories, pagination, nil
}

func (s *HistoryService) GetAllHistories(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History, *paging.Pagination, error) {
	histories, pagination, err := s.repository.GetHistories(c, req)
	if err != nil {
		log.Printf("Failed to get all histories: %v", err)
		return nil, nil, fmt.Errorf("failed to get all histories: %v", err)
	}

	return histories, pagination, nil
}
