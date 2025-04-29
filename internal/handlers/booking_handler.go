package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/models"
)

func BookEvent(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(float64))

	eventIDParam := c.Param("id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	result := db.DB.First(&event, eventID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var totalBookings int64
	db.DB.Model(&models.Booking{}).Where("event_id = ?", eventID).Count(&totalBookings)

	if int(totalBookings) >= event.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event is fully booked"})
		return
	}

	var existingBooking models.Booking
	bookingResult := db.DB.Where("user_id = ? AND event_id = ?", userID, eventID).First(&existingBooking)
	if bookingResult.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already booked this event"})
		return
	}

	booking := models.Booking{
		UserID:  userID,
		EventID: uint(eventID),
	}

	
	if err := db.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Booking successful", "booking": booking})
}

func CancelBooking(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := uint(userIDInterface.(float64)) 

	bookingIDParam := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var booking models.Booking
	result := db.DB.First(&booking, bookingID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only cancel your own bookings"})
		return
	}

	db.DB.Delete(&booking)

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

func GetBookings(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDInterface.(float64))

	var bookings []models.Booking
	result := db.DB.Where("user_id = ?", userID).Find(&bookings)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve bookings"})
	}
	var response []gin.H
	for _, booking := range bookings {
		var event models.Event
		db.DB.First(&event, booking.EventID)

		response = append(response, gin.H{
			"booking_id":       booking.ID,
			"event_id":         event.ID,
			"event_name":       event.Name,
			"event_location":   event.Location,
			"event_start_time": event.StartTime,
			"event_end_time":   event.EndTime,
		})
	}
	c.JSON(http.StatusOK, gin.H{"my_bookings": response})

}
