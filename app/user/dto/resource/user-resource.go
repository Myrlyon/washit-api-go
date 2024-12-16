package resource

import (
	"time"
)

type User struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email" `
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
