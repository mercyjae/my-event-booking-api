package repo

import (
	"fmt"

	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/domain"
)

func SaveEvent(e *domain.Event) error {
	query := `
	INSERT INTO events(name, description, location_venue, location_address, event_date, user_id, capacity)
	VALUES (?,?,?,?,?,?,?)`

	stmt, err := db.DBB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.LocationVenue, e.LocationAddress, e.EventDate, e.UserId, e.Capacity)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	e.ID = id
	return err
}

func GetAllEvents() ([]domain.Event, error) {
	query := "SELECT * FROM events"

	rows, err := db.DBB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event

		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.LocationVenue, &event.LocationAddress, &event.EventDate, &event.UserId, &event.Capacity)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil

}

func GetEventById(id int64) (*domain.Event, error) {
	query := "SELECT * FROM events WHERE id  = ?"
	row := db.DBB.QueryRow(query, id)

	var event domain.Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.LocationVenue, &event.LocationAddress, &event.EventDate, &event.UserId, &event.Capacity)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func Update(event *domain.Event) error {
	query := `
UPDATE events
SET name = ?, description = ?, location_venue = ?, location_address = ?, event_date = ?, user_id = ?, capacity = ?
WHERE id = ?
`
	stmt, err := db.DBB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(&event.Name, &event.Description, &event.LocationVenue, &event.LocationAddress, &event.EventDate, &event.UserId, &event.Capacity, &event.ID, )
	return err
}

func Delete(event *domain.Event) error {
	query := "DELETE FROM events WHERE id = ?"

	stmt, err := db.DBB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(event.ID)
	return err
}
