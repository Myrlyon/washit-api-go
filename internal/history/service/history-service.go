package historyService

import (
	historyRepository "washit-api/internal/history/repository"
)

type HistoryServiceInterface interface{}

type HistoryService struct {
	repository historyRepository.HistoryRepositoryInterface
}

func NewHistoryService(repository historyRepository.HistoryRepositoryInterface) *HistoryService {
	return &HistoryService{
		repository: repository,
	}
}
