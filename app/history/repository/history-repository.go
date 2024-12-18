package historyRepository

import dbs "washit-api/db"

type HistoryRepositoryInterface interface{}

type HistoryRepository struct {
	db dbs.DatabaseInterface
}

func NewHistoryRepository(db dbs.DatabaseInterface) *HistoryRepository {
	return &HistoryRepository{db: db}
}
