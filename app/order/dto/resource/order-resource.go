package orderResource

import (
	"time"
	"washit-api/utils"

	"github.com/shopspring/decimal"
)

type Base struct {
	ID            string           `json:"id" gorm:"primaryKey"`
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
}

type Order struct {
	Message string `json:"message"`
	Order   Base   `json:"order"`
}

type OrderWithLinks struct {
	Order      Base                  `json:"order"`
	Hypermedia map[string]utils.Link `json:"_links"`
}
