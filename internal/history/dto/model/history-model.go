package historyModel

import (
	"time"
	userModel "washit-api/internal/user/dto/model"
)

type History struct {
	ID            string         `json:"id" gorm:"primaryKey unique"`
	UserID        int            `json:"userId" gorm:"not null;index"`
	TransactionID int            `json:"transactionId"`
	AddressID     int            `json:"addressId"`
	Status        string         `json:"status"`
	Note          string         `json:"note"`
	ServiceType   string         `json:"serviceType"`
	OrderType     string         `json:"orderType"`
	Price         float64        `json:"price"`
	Reason        string         `json:"reason"`
	CollectDate   time.Time      `json:"collectDate"`
	EstimateDate  time.Time      `json:"estimateDate"`
	DeletedAt     time.Time      `json:"deletedAt"`
	User          userModel.User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
