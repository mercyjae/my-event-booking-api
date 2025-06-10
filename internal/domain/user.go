package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	FullName     string    `json:"full_name"`
	Password     Password  `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	DoB          time.Time `json:"date_of_birth"`
	Verified     bool      `json:"verified"`
	OTP          string    `json:"otp"`
	OTPExpiresAt time.Time `json:"otp_expires_at"`
}

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
