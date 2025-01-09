package transactionModel

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID             string           `json:"id" gorm:"primaryKey unique"`
	OrderID        string           `json:"orderID"`
	UserID         int64            `json:"userID" gorm:"not null;index"`
	ExternalID     string           `json:"externalID"`
	PaymentMethod  string           `json:"paymentMethod"`
	Status         string           `json:"status"`
	Amount         *decimal.Decimal `json:"amount" gorm:"type:numeric"`
	PaymentChannel string           `json:"paymentChannel"`
	Description    string           `json:"description"`
	PaidAt         time.Time        `json:"paidAt"`
}
