package userResource

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	Image     string    `json:"image"`
}

type ListUser struct {
	User []User `json:"users"`
}

type HideToken struct {
	User User `json:"user"`
}

type ShowToken struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
