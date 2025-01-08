package transaction

import (
	"net/http"
	transactionService "washit-api/internal/transaction/service"
	"washit-api/pkg/redis"
	"washit-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service transactionService.TransactionServiceInterface
	cache   redis.IRedis
}

func NewTransactionHandler(service transactionService.TransactionServiceInterface, cache redis.IRedis) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		cache:   cache,
	}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	response.Success(c, http.StatusOK, "Test", nil, nil)
}
