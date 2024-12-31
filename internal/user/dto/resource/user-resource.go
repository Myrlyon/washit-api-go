package userResource

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

type WithToken struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
