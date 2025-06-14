package dto

import (
	"time"
)

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Message              string    `json:"message"`
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	// RefreshToken          string    `json:"refresh_token"`
	// RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type ResetPassword struct {
	Email           string `json:"email"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ChangePassword struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
