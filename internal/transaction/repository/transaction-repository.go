package transactionRepository

import "washit-api/pkg/db/dbs"

type TransactionRepositoryInterface interface{}

type TransactionRepository struct {
	db dbs.DatabaseInterface
}

func NewTransactionRepository(db dbs.DatabaseInterface) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}
