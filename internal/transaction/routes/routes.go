package transactionRoutes

import (
	transaction "washit-api/app/transaction/handler"
	transactionRepository "washit-api/app/transaction/repository"
	transactionService "washit-api/app/transaction/service"
	dbs "washit-api/db"
	"washit-api/middleware"
	"washit-api/redis"

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
