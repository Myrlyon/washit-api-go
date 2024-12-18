package orderModel

import (
	"time"

	userModel "washit-api/app/user/dto/model"
)

type Order struct {
	ID            string         `json:"id" gorm:"primaryKey unique"`
	UserID        int            `json:"userId" gorm:"not null;index"`
	TransactionID int            `json:"transactionId"`
	AddressID     int            `json:"addressId"`
	Status        string         `json:"status"`
	Note          string         `json:"note"`
	ServiceType   string         `json:"serviceType"`
	OrderType     string         `json:"orderType"`
	Price         float64        `json:"price"`
	CollectDate   time.Time      `json:"collectDate"`
	EstimateDate  time.Time      `json:"estimateDate"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	User          userModel.User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
