package models

import "time"

type Booking struct {
	ID       uint      `json:"id"`
	UserID   int       `json:"user_id"`
	EventID  int       `json:"event_id"`
	Seats    int       `json:"seats"`
	BookedAt time.Time `json:"booked_at"`
}
