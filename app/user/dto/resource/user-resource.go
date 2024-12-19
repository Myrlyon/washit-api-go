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

type ListUser struct {
	User []Base `json:"users"`
}

type HideToken struct {
	User Base `json:"user"`
}

type ShowToken struct {
	User         Base   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
