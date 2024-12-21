package orderResource

import (
	"time"

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

type OrderList struct {
	Message string `json:"message,omitempty"`
	Orders  []Base `json:"orders"`
}

type Order struct {
	Message    string     `json:"message,omitempty"`
	Order      Base       `json:"order"`
	Hypermedia Hypermedia `json:"_links,omitempty"`
}

type Hypermedia struct {
	Self   map[string]string `json:"self,omitempty"`
	Create map[string]string `json:"create,omitempty"`
	Cancel map[string]string `json:"cancel,omitempty"`
}
