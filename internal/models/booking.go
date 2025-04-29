package models

import "time"

type Booking struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `json:"user_id"`
	EventID   uint `json:"event_id"`
	CreatedAt time.Time
}
