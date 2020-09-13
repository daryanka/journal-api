package auth

type CreateUser struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required"`
}

type TokenWithClaims struct {
	Token string `json:"token"`
	Expires int64 `json:"expires"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID    int64 `json:"id"`
	Email string `json:"email"`
}