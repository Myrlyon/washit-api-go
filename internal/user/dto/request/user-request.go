package userRequest

type Register struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	FcmToken string `json:"fcmToken"`
}

type Update struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password" validate:"min=3,max=130"`
}

type Google struct {
	IDToken  string `json:"idToken"`
	FcmToken string `json:"fcmToken"`
}
