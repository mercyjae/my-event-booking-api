package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/repo"
)

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	user, err := repo.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user", "devError": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "devError": fmt.Sprintf("No user with ID %d found", uid)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"full_name": user.FullName,
		"email":     user.Email,
		"phone":     user.Phone,
		//"created_at": user.CreatedAt,
	})
}

func EditProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uid := userID.(int)

	var req struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
		
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repo.UpdateUserProfile(uid, req.FullName, req.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update profile", "devError": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": gin.H{
			"full_name": req.FullName,
			"phone":     req.Phone,
		},
	})
}
