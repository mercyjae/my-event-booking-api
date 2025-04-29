package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/models"
)

func CreateEvent(c *gin.Context) {
	var req models.Event

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	event := models.Event{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Capacity:    req.Capacity,
	}

	db.DB.Create(&event)

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event": event})

}


func ListEvents(c *gin.Context) {
    var events []models.Event
    db.DB.Find(&events)

    c.JSON(http.StatusOK, gin.H{"events": events})
}


func GetEvent(c *gin.Context) {
    id := c.Param("id")

    var event models.Event
    result := db.DB.First(&event, id)

    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"event": event})
}

func DeleteEvent(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
        return
    }

    var event models.Event
    result := db.DB.First(&event, uint(id))
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
        return
    }

    if err := db.DB.Delete(&event).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
