package repo

import (
	"database/sql"
	"fmt"

	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/domain"
)

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

func GetUserPasswordHash(userID int) (string, error) {
	var password string
	query := `SELECT password FROM users WHERE id = ?`
	err := db.DBB.QueryRow(query, userID).Scan(&password)
	return password, err
}

// UpdateUserPassword updates the user password in the DB
func UpdateUserPassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password = ? WHERE id = ?`
	_, err := db.DBB.Exec(query, hashedPassword, userID)
	return err
}

func GetAllUsers() ([]domain.User, error) {
	query := `SELECT id, full_name, email, phone FROM users ORDER BY id DESC`

	rows, err := db.DBB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Phone); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
