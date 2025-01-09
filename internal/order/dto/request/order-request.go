package orderRequest

import (
	"time"
)

type Order struct {
	AddressID   int       `json:"addressID" validate:"required"`
	Note        string    `json:"note"`
	ServiceType string    `json:"serviceType" validate:"required"`
	OrderType   string    `json:"orderType" validate:"required"`
	CollectDate time.Time `json:"collectDate" validate:"required"`
}

type Payment struct {
	TransactionID string `json:"transactionID" validate:"required"`
}
