package transaction

import (
	"net/http"
	transactionService "washit-api/internal/transaction/service"
	"washit-api/pkg/redis"
	"washit-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service transactionService.TransactionServiceInterface
	cache   redis.RedisInterface
}

func NewTransactionHandler(service transactionService.TransactionServiceInterface, cache redis.RedisInterface) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		cache:   cache,
	}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	utils.WriteJson(c, http.StatusOK, map[string]interface{}{"transactions": "transactions"})
}
