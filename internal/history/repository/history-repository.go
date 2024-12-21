package historyRepository

import "washit-api/pkg/db/dbs"

type HistoryRepositoryInterface interface{}

type HistoryRepository struct {
	db dbs.DatabaseInterface
}

func NewHistoryRepository(db dbs.DatabaseInterface) *HistoryRepository {
	return &HistoryRepository{db: db}
}
