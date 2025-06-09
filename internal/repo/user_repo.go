package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

//import "github.com/mercyjae/event-booking-api/internal/db"

func SaveUser(u *domain.User) error {
	query := "INSERT INTO users(full_name, email, password, phone) VALUES (?, ?, ?, ?)"
	stmt, err := db.DBB.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare insert user failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.FullName, u.Email, u.Password.Hash, u.Phone)
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

// func GetUserByEmail(u *domain.User) error {
// 	fmt.Println("Searching for email:", u.Email)
// 	query := "SELECT id, password FROM users WHERE email = ?"
// 	row := db.DBB.QueryRow(query, u.Email)
// 	var retrievedPassword string
// 	err := row.Scan(&u.ID, &retrievedPassword)

// 	// if err != nil {
// 	// 	return errors.New("Credentials Invalid")
// 	// }
// 	if err == sql.ErrNoRows {
// 		// Return an error if the email is not found in the database
// 		return errors.New("user not found")
// 	} else if err != nil {
// 		// Return any other error encountered during the query
// 		return err
// 	}
// 	// match, err := doai..Matches(retrievedPassword)
// 	// passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

// 	// if !passwordIsValid {
// 	// 	return errors.New("Credentials Invalid")
// 	// }

//		return nil
//	}
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
	query := `SELECT id FROM users WHERE email = ?`
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
	_, err = db.DBB.Exec(updateQuery, string(hashedPassword), userID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(userID int) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, full_name, email, phone FROM users WHERE id = ?`

	err := db.DBB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		//&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // user not found
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	fmt.Println("Running query for userID:", userID)

	return &user, nil
}

func UpdateUserProfile(userID int, fullName string, phone string) error {
	// Optional: check user exists
	var id int
	checkQuery := `SELECT id FROM users WHERE id = ?`
	err := db.DBB.QueryRow(checkQuery, userID).Scan(&id)
	if err != nil {
		return err // could be sql.ErrNoRows
	}

	// Update query
	updateQuery := `UPDATE users SET full_name = ?, phone = ? WHERE id = ?`
	_, err = db.DBB.Exec(updateQuery, fullName, phone, userID)
	return err
}
