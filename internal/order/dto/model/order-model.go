package orderModel

import (
	"time"

	userModel "washit-api/internal/user/dto/model"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID            string           `json:"id" gorm:"primaryKey unique"`
	UserID        int              `json:"userId" gorm:"not null;index"`
	TransactionID int              `json:"transactionId"`
	AddressID     int              `json:"addressId"`
	Status        string           `json:"status"`
	Note          string           `json:"note"`
	ServiceType   string           `json:"serviceType"`
	OrderType     string           `json:"orderType"`
	Price         *decimal.Decimal `json:"price" gorm:"type:numeric"`
	CollectDate   time.Time        `json:"collectDate"`
	EstimateDate  time.Time        `json:"estimateDate"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
	User          userModel.User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
