package historyRepository

import (
	historyModel "washit-api/internal/history/dto/model"
	historyRequest "washit-api/internal/history/dto/request"
	"washit-api/pkg/db/dbs"
	"washit-api/pkg/paging"

	"github.com/gin-gonic/gin"
)

type IHistoryRepository interface {
	GetHistories(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History, *paging.Pagination, error)
	GetHistoryByID(c *gin.Context, historyID string) (*historyModel.History, error)
}

type HistoryRepository struct {
	db dbs.IDatabase
}

func NewHistoryRepository(db dbs.IDatabase) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) GetHistoryByID(c *gin.Context, historyID string) (*historyModel.History, error) {
	var history historyModel.History
	if err := r.db.FindByID(c, historyID, &history); err != nil {
		return nil, err
	}

	return &history, nil
}

func (r *HistoryRepository) GetHistories(c *gin.Context, req *historyRequest.ListHistory) ([]*historyModel.History, *paging.Pagination, error) {
	query := []dbs.Query{
		dbs.NewQuery("user_id = ?", req.UserID),
	}

	if req.Code != "" {
		query = append(query, dbs.NewQuery("code = ?", req.Code))
	}
	if req.Status != "" {
		query = append(query, dbs.NewQuery("status = ?", req.Status))
	}

	order := "deleted_at"
	if req.OrderBy != "" {
		order = req.OrderBy
		if req.OrderDesc {
			order += " DESC"
		}
	}

	var total int64
	if err := r.db.Count(c, &historyModel.History{}, &total, dbs.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var histories []*historyModel.History
	if err := r.db.Find(
		c,
		&histories,
		dbs.WithPreload([]string{"User"}),
		dbs.WithQuery(query...),
		dbs.WithLimit(int(pagination.Limit)),
		dbs.WithOffset(int(pagination.Skip)),
		dbs.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}

	return histories, pagination, nil
}
