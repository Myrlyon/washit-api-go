package transactionRoutes

import (
	transaction "washit-api/internal/transaction/handler"
	transactionRepository "washit-api/internal/transaction/repository"
	transactionService "washit-api/internal/transaction/service"
	"washit-api/pkg/db/dbs"
	"washit-api/pkg/middleware"
	"washit-api/pkg/redis"

	"github.com/gin-gonic/gin"
)

func Main(r *gin.RouterGroup, db dbs.DatabaseInterface, cache redis.RedisInterface) {
	repository := transactionRepository.NewTransactionRepository(db)
	service := transactionService.NewTransactionService(repository)
	handler := transaction.NewTransactionHandler(service, cache)

	authMiddleware := middleware.JWTAuth()
	// adminAuthMiddleware := middleware.JTWAuthAdmin()

	r.GET("/transactions", authMiddleware, handler.GetTransactions)
}
