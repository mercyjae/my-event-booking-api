package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/domain"
	"github.com/mercyjae/event-booking-api/internal/dto"
	"github.com/mercyjae/event-booking-api/internal/models"
	"github.com/mercyjae/event-booking-api/internal/repo"
)

func CreateEvent(c *gin.Context) {
	var req dto.EventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := c.GetInt("user_id")

	event := domain.Event{
		Name:            req.Name,
		Description:     req.Description,
		LocationAddress: req.LocationAddress,
		LocationVenue:   req.LocationVenue,
		//StartTime:       req.StartTime,
		EventDate: req.EventDate,
		// EndTime:         req.EndTime,
		Capacity: req.Capacity,
		UserId:   userId,
	}

	//event.UserId = int(userId)
	err := repo.SaveEvent(&event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create events, Try again later", "devError": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event": event})

}

func ListEvents(c *gin.Context) {

	events, err := repo.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events, Try again later", "devError": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

//c.JSON(http.StatusOK, gin.H{"events": events})
//}

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

func UpdateEvent(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id"})
		return
	}

	var event models.Event
	result := db.DB.First(&event, uint(id))

	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event"})
		return
	}
	// userId := context.GetInt("user_id")
	// if event.UserId != userId {
	// 	context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event"})
	// 	return
	// }
	var req models.UpdateEventRequest
	err = context.ShouldBindJSON(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	if req.Name != nil {
		event.Name = *req.Name
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.LocationVenue != nil {
		event.LocationVenue = *req.LocationVenue
	}
	if req.LocationAddress != nil {
		event.LocationAddress = *req.LocationAddress
	}
	if req.EventDate != nil {
		event.EventDate = *req.EventDate
	}
	// if req.StartTime != nil {
	// 	event.StartTime = *req.StartTime
	// }
	// if req.EndTime != nil {
	// 	event.EndTime = *req.EndTime
	// }
	if req.Capacity != nil {
		event.Capacity = *req.Capacity
	}

	// Save updated event
	if err := db.DB.Save(&event).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}
