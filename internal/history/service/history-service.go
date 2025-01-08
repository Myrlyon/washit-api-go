package historyService

import (
	"fmt"
	"log"
	historyModel "washit-api/internal/history/dto/model"
	historyRepository "washit-api/internal/history/repository"

	"github.com/gin-gonic/gin"
)

type IHistoryService interface {
	GetHistoriesMe(c *gin.Context, userId string) ([]*historyModel.History, error)
}

type HistoryService struct {
	repository historyRepository.HistoryRepositoryInterface
}

func NewHistoryService(repository historyRepository.HistoryRepositoryInterface) *HistoryService {
	return &HistoryService{
		repository: repository,
	}
}

func (s *HistoryService) GetHistoriesMe(c *gin.Context, userId string) ([]*historyModel.History, error) {
	histories, err := s.repository.GetHistoriesByUserId(c, userId)
	if err != nil {
		log.Println("Failed to get histories by user id ", err)
		return nil, fmt.Errorf("failed to get histories by user id: %v", err)
	}

	return histories, nil
}
