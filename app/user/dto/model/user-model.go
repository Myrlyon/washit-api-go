package userModel

import "time"

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey unique"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email" gorm:"unique"`
	Role      string    `json:"role" gorm:"default:customer"`
	Password  string    `json:"-"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
