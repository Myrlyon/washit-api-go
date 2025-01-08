package orderResource

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID string `json:"id" gorm:"primaryKey"`
	// UserID        int              `json:"userId" gorm:"not null;index"`
	User          User             `json:"user" gorm:"foreignKey:UserID;references:ID"`
	TransactionID string           `json:"transactionId"`
	AddressID     int              `json:"addressId"`
	Status        string           `json:"status"`
	Note          string           `json:"note"`
	ServiceType   string           `json:"serviceType"`
	OrderType     string           `json:"orderType"`
	Weight        *float64         `json:"weight"`
	Price         *decimal.Decimal `json:"price" gorm:"type:numeric"`
	CollectDate   time.Time        `json:"collectDate"`
	EstimateDate  time.Time        `json:"estimateDate"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Image     string `json:"image"`
	// CreatedAt time.Time `json:"createdAt"`
}
