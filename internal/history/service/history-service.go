package historyService

import (
	"fmt"
	"log"
	historyModel "washit-api/internal/history/dto/model"
	historyRepository "washit-api/internal/history/repository"

	"github.com/gin-gonic/gin"
)

type IHistoryService interface {
	GetHistoryByID(c *gin.Context, historyID string, userID int64) (*historyModel.History, error)
	GetHistoriesMe(c *gin.Context, userID int64) ([]*historyModel.History, error)
	GetHistoriesByUser(c *gin.Context, userID int64) ([]*historyModel.History, error)
}

type HistoryService struct {
	repository historyRepository.IHistoryRepository
}

func NewHistoryService(repository historyRepository.IHistoryRepository) *HistoryService {
	return &HistoryService{
		repository: repository,
	}
}

func (s *HistoryService) GetHistoryByID(c *gin.Context, historyID string, userID int64) (*historyModel.History, error) {
	return nil, nil
}

func (s *HistoryService) GetHistoriesMe(c *gin.Context, userID int64) ([]*historyModel.History, error) {
	histories, err := s.repository.GetHistoriesByUser(c, userID)
	if err != nil {
		log.Println("Failed to get histories by user id ", err)
		return nil, fmt.Errorf("failed to get histories by user id: %v", err)
	}

	return histories, nil
}

func (s *HistoryService) GetHistoriesByUser(c *gin.Context, userID int64) ([]*historyModel.History, error) {
	histories, err := s.repository.GetHistoriesByUser(c, userID)
	if err != nil {
		log.Println("Failed to get histories by user id ", err)
		return nil, fmt.Errorf("failed to get histories by user id: %v", err)
	}

	return histories, nil
}
