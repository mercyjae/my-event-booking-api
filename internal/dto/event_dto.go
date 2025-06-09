package dto

import "time"

type EventRequest struct {

	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	LocationVenue   string `json:"location_venue" binding:"required"`
	LocationAddress string `json:"location_address" binding:"required"`
	EventDate       time.Time `json:"event_date" binding:"required"`
	//UserId          int    `json:"user_id"`
	// StartTime       string `json:"start_time"`
	// EndTime         string `json:"end_time"`
	Capacity  int `json:"capacity" binding:"required"`
	CreatedAt time.Time
}
