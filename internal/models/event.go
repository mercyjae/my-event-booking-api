package models

import "time"

type Event struct {
	ID              int64  `json:"id"`
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	LocationVenue   string `json:"location_venue" binding:"required"`
	LocationAddress string `json:"location_address" binding:"required"`
	EventDate       string `json:"event_date" binding:"required"`
	UserId          int    `json:"user_id"`
	// StartTime       string `json:"start_time"`
	// EndTime         string `json:"end_time"`
	Capacity  int `json:"capacity" binding:"required"`
	CreatedAt time.Time
}

type UpdateEventRequest struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	LocationVenue   *string `json:"location_venue"`
	LocationAddress *string `json:"location_address"`
	EventDate       *string `json:"event_date"`
	// StartTime       *string `json:"start_time"`
	// EndTime         *string `json:"end_time"`
	Capacity *int `json:"capacity"`
}
