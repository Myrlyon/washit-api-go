package transactionService

import transactionRepository "washit-api/internal/transaction/repository"

type TransactionServiceInterface interface {
}

type transactionService struct {
	repository transactionRepository.TransactionRepositoryInterface
}

func NewTransactionService(repository transactionRepository.TransactionRepositoryInterface) *transactionService {
	return &transactionService{
		repository: repository,
	}
}
