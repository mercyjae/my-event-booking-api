package models

import "time"

type RegisterUser struct {
	ID           uint    `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	FullName     string    `json:"full_name"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	DoB          time.Time `json:"date_of_birth"`
	Verified     bool      `gorm:"default:false"`
	OTP          string    `gorm:"size:6"`
	OTPExpiresAt time.Time
}

type VerifyOTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
