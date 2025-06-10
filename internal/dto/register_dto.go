package dto

import "time"

type RegisterUserRequest struct {
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	FullName string    `json:"full_name"`
	Password string    `json:"password"`
	DoB      time.Time `json:"date_of_birth"`
}
