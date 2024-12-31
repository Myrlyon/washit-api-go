package orderRequest

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID            string           `json:"id"`
	TransactionID string           `json:"transactionId"`
	AddressID     int              `json:"addressId" validate:"required"`
	Status        string           `json:"status"`
	Note          string           `json:"note"`
	ServiceType   string           `json:"serviceType" validate:"required"`
	OrderType     string           `json:"orderType" validate:"required"`
	Price         *decimal.Decimal `json:"price" validate:"required"`
	CollectDate   time.Time        `json:"collectDate" validate:"required"`
	EstimateDate  time.Time        `json:"estimateDate" validate:"required"`
}
