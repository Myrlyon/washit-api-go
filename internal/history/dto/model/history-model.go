package historyModel

import (
	"time"
	userModel "washit-api/internal/user/dto/model"
)

type History struct {
	ID            string         `json:"id" gorm:"primaryKey unique"`
	UserID        int64          `json:"userID" gorm:"not null;index"`
	TransactionID int            `json:"transactionID"`
	AddressID     int            `json:"addressID"`
	Status        string         `json:"status"`
	Note          string         `json:"note"`
	ServiceType   string         `json:"serviceType"`
	OrderType     string         `json:"orderType"`
	Price         float64        `json:"price"`
	CollectDate   time.Time      `json:"collectDate"`
	EstimateDate  time.Time      `json:"estimateDate"`
	DeletedAt     time.Time      `json:"deletedAt"`
	Reason        string         `json:"reason"`
	User          userModel.User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
