package transactionRepository

import dbs "washit-api/db"

type TransactionRepositoryInterface interface{}

type TransactionRepository struct {
	db dbs.DatabaseInterface
}

func NewTransactionRepository(db dbs.DatabaseInterface) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}
