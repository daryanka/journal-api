package auth

type CreateUser struct {
	Email                string `json:"email" validate:"required,email,max=255"`
	Password             string `json:"password" validate:"required,max=50,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,max=550,min=6"`
}

type TokenWithClaims struct {
	Token     string `json:"token"`
	Expires   int64  `json:"expires"`
	ExpiresIn int    `json:"expires_in"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=50,min=6"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
