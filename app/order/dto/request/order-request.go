package orderRequest

import (
	"time"
)

type Order struct {
	ID            string    `json:"id"`
	TransactionID int       `json:"transactionId"`
	AddressID     int       `json:"addressId" validate:"required"`
	Status        string    `json:"status"`
	Note          string    `json:"note"`
	ServiceType   string    `json:"serviceType" validate:"required"`
	OrderType     string    `json:"orderType" validate:"required"`
	Price         float64   `json:"price" validate:"required"`
	CollectDate   time.Time `json:"collectDate" validate:"required"`
	EstimateDate  time.Time `json:"estimateDate" validate:"required"`
}
