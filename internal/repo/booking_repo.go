package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/domain"
)

func BookEvent(userID int, eventID int, seats int) error {
	// Check if already booked
	var count int
	err := db.DBB.QueryRow(`
		SELECT COUNT(*) FROM bookings WHERE user_id = ? AND event_id = ?
	`, userID, eventID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing booking: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("you already booked this event")
	}

	// Check if event is full
	var totalSeats sql.NullInt64
	err = db.DBB.QueryRow(`
		SELECT SUM(seats) FROM bookings WHERE event_id = ?
	`, eventID).Scan(&totalSeats)
	if err != nil {
		return fmt.Errorf("failed to check event capacity")
	}

	currentSeats := 0
	if totalSeats.Valid {
		currentSeats = int(totalSeats.Int64)
	}

	var event domain.Event
	err = db.DBB.QueryRow(`
		SELECT id, capacity FROM events WHERE id = ?
	`, eventID).Scan(&event.ID, &event.Capacity)
	if err != nil {
		return fmt.Errorf("event not found or failed to fetch")
	}

	if currentSeats+seats > event.Capacity {
		return fmt.Errorf("event is fully booked")
	}

	// Insert booking
	stmt, err := db.DBB.Prepare(`
		INSERT INTO bookings (user_id, event_id, seats, booked_at)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare booking insert")
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, eventID, seats, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create booking")
	}

	return nil
}

func GetUserBookings(userID int) ([]gin.H, error) {
	query := `
		SELECT b.id, b.event_id, e.name, e.location_address
		FROM bookings b
		INNER JOIN events e ON b.event_id = e.id
		WHERE b.user_id = ?
	`

	rows, err := db.DBB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var bookings []gin.H
	for rows.Next() {
		var bookingID, eventID int
		var eventName, eventLocation string

		if err := rows.Scan(&bookingID, &eventID, &eventName, &eventLocation); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		bookings = append(bookings, gin.H{
			"booking_id":     bookingID,
			"event_id":       eventID,
			"event_name":     eventName,
			"event_location": eventLocation,
		})
	}

	return bookings, nil
}

func DeleteBookingByID(bookingID int64) error {
	query := "DELETE FROM bookings WHERE id = ?"
	stmt, err := db.DBB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(bookingID)
	if err != nil {
		return fmt.Errorf("failed to execute delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
