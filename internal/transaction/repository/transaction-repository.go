package transactionRepository

import "washit-api/pkg/db/dbs"

type TransactionRepositoryInterface interface{}

type TransactionRepository struct {
	db dbs.IDatabase
}

func NewTransactionRepository(db dbs.IDatabase) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}
