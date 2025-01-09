package userRequest

type Register struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=130"`
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
	OldPassword     string `json:"oldPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=6,max=130"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}
type UpdatePicture struct {
	Image []byte `json:"Image" validate:"required,file,ext=jpg|png|webp,max=2mb"`
}

type Google struct {
	TokenID  string `json:"tokenID"`
	FcmToken string `json:"fcmToken"`
}
