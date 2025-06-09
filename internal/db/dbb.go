package db

import (
	"database/sql"
	"log"

	//	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DBB *sql.DB

func InitDB() {
	var err error
	DBB, err = sql.Open("sqlite3", "event.db")

	if err != nil {
		panic("could not connect to database")
	}

	DBB.SetMaxOpenConns(10)
	DBB.SetMaxIdleConns(5)

	createTable()
}

func createTable() {
	// _, err := DBB.Exec(`DROP TABLE IF EXISTS users`)
	// if err != nil {
	// 	log.Fatal("❌ Failed to drop users table:", err)
	// }

	// _, err = DBB.Exec(`DROP TABLE IF EXISTS events`)
	// if err != nil {
	// 	log.Fatal("❌ Failed to drop events table:", err)
	// }

	// _, err = DBB.Exec(`DROP TABLE IF EXISTS users`)
	// if err != nil {
	// 	log.Fatal("❌ Failed to drop users table:", err)
	// }

	createUsersTable := `

	CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	full_name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	phone TEXT NOT NULL,
	password TEXT NOT NULL,
	otp TEXT,
	otp_expires_at DATETIME
	

	)
	`
	_, err := DBB.Exec(createUsersTable)
	if err != nil {
		panic("Could not create users table")
	}

	// _, err = DBB.Exec(`DROP TABLE IF EXISTS events`)
	// if err != nil {
	// 	log.Fatal("Failed to drop events table:", err)
	// }
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	location_venue TEXT NOT NULL,
	location_address TEXT NOT NULL,
	event_date DATETIME NOT NULL,
	user_id  INTEGER,
	capacity INTEGER,
	FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DBB.Exec(createEventsTable)

	if err != nil {
		panic("Could not create events table")
	}

	createBookingsTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		event_id INTEGER NOT NULL,
		seats INTEGER NOT NULL,
		booked_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(event_id) REFERENCES events(id)
	);
	`
	_, err = DBB.Exec(createBookingsTable)
	if err != nil {
		log.Fatal("Failed to create bookings table:", err)
	}

}
