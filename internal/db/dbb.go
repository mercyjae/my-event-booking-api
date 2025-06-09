package db

import (
	"database/sql"

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

	createUsersTable := `

	CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	full_name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	phone TEXT NOT NULL,
	password TEXT NOT NULL

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

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id INTEGER,
	user_id INTEGER,
	FOREIGN KEY(event_id) REFERENCES events(id),
	FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DBB.Exec(createRegistrationsTable)

	if err != nil {
		panic("Could not create registrations table")
	}
}
