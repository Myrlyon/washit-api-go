package historyResource

import (
	"time"
	"washit-api/pkg/paging"
)

type ListHistory struct {
	Histories     []*History           `json:"histories,omitempty"`
	Pagination *paging.Pagination `json:"pagination,omitempty"`
}

type History struct {
	ID            string    `json:"id" gorm:"primaryKey unique"`
	User          User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	TransactionID int       `json:"transactionID"`
	AddressID     int       `json:"addressID"`
	Status        string    `json:"status"`
	Note          string    `json:"note"`
	ServiceType   string    `json:"serviceType"`
	OrderType     string    `json:"orderType"`
	Price         float64   `json:"price"`
	CollectDate   time.Time `json:"collectDate"`
	EstimateDate  time.Time `json:"estimateDate"`
	DeletedAt     time.Time `json:"deletedAt"`
	Reason        string    `json:"reason"`
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
