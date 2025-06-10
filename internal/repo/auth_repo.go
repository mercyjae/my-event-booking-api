package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func SaveUser(u *domain.User) error {
	query := "INSERT INTO users(full_name, email, password, phone, otp, otp_expires_at, verified) VALUES (?, ?, ?, ?,?,?,?)"
	stmt, err := db.DBB.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare insert user failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.FullName, u.Email, u.Password.Hash, u.Phone, u.OTP, u.OTPExpiresAt, u.Verified)
	if err != nil {
		return fmt.Errorf("exec insert user failed: %w", err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("fetching last insert ID failed: %w", err)
	}

	u.ID = uint(userId)
	return nil
}

func IsEmailTaken(email string) (bool, error) {
	var count int
	err := db.DBB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserByEmail(email string) (*domain.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = ?"
	row := db.DBB.QueryRow(query, email)

	var user domain.User
	var hashedPassword string

	err := row.Scan(&user.ID, &user.Email, &hashedPassword)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	user.Password = domain.Password{Hash: []byte(hashedPassword)}
	return &user, nil
}

func ResetPasswordByEmail(email, newPassword string) error {
	var userID int

	// Normalize email comparison
	query := `SELECT id FROM users WHERE TRIM(LOWER(email)) = TRIM(LOWER(?))`
	err := db.DBB.QueryRow(query, email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE users SET password = ? WHERE id = ?`
	result, err := db.DBB.Exec(updateQuery, hashedPassword, userID)
	if err != nil {
		return err
	}
	fmt.Println("🔐 New hashed password:", hashedPassword)
	fmt.Println("🔐 New hashed password:", string(hashedPassword))

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

// func ResetPasswordByEmail(email, newPassword string) error {
// 	var userID int
// 	query := `SELECT id FROM users WHERE email = ?`
// 	err := db.DBB.QueryRow(query, email).Scan(&userID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return errors.New("user not found")
// 		}
// 		return err
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	updateQuery := `UPDATE users SET password = ? WHERE id = ?`
// 	_, err = db.DBB.Exec(updateQuery, string(hashedPassword), userID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
