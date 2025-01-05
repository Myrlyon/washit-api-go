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

type UpdateProfile struct {
	FirstName string `json:"firstName" validate:"min=2"`
	LastName  string `json:"lastName" validate:"min=2"`
	Email     string `json:"email"`
}

type UpdatePassword struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=3,max=130"`
}

type Google struct {
	IDToken  string `json:"idToken"`
	FcmToken string `json:"fcmToken"`
}
