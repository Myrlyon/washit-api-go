package userModel

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email" gorm:"unique"`
	Role      string    `json:"role" gorm:"default:customer"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
