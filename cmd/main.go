package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/handlers"
)

func main() {
	r := gin.Default()
	db.ConnectDatabase()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/register", handlers.RegisterUser)
	r.POST("/verify-otp", handlers.VerifyOTP)
	r.Run(":8080")
	//fmt.Println("Project works!")
}
