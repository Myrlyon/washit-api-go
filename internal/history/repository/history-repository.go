package historyRepository

import (
	historyModel "washit-api/internal/history/dto/model"
	"washit-api/pkg/db/dbs"

	"github.com/gin-gonic/gin"
)

type HistoryRepositoryInterface interface {
	GetHistoriesByUserId(c *gin.Context, userId string) (*[]historyModel.History, error)
}

type HistoryRepository struct {
	db dbs.IDatabase
}

func NewHistoryRepository(db dbs.IDatabase) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) GetHistoriesByUserId(c *gin.Context, userId string) (*[]historyModel.History, error) {
	var histories *[]historyModel.History
	query := []dbs.FindOption{
		dbs.WithLimit(10),
		dbs.WithOrder("created_at DESC"),
		dbs.WithPreload([]string{"User"}),
	}

	// if userId != "" {
	// 	query = append(query, dbs.WithQuery(dbs.NewQuery("user_id = ?", userId)))
	// }

	if err := r.db.Find(c, &histories, query...); err != nil {
		return nil, err
	}

	return histories, nil
}
