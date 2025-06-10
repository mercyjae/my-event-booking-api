package dto

import "time"

type UpdateEventRequest struct {
	Name            *string    `json:"name"`
	Description     *string    `json:"description"`
	LocationVenue   *string    `json:"location_venue"`
	LocationAddress *string    `json:"location_address"`
	EventDate       *time.Time `json:"event_date"`
	// StartTime       *string `json:"start_time"`
	// EndTime         *string `json:"end_time"`
	Capacity *int `json:"capacity"`
}
