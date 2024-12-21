package userResource

import (
	"time"
)

type Base struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserList struct {
	Message string `json:"message,omitempty"`
	Users   []Base `json:"users"`
}

type User struct {
	Message      string `json:"message,omitempty"`
	User         Base   `json:"user"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
