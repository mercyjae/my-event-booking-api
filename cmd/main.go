package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/handlers"
	"github.com/mercyjae/event-booking-api/internal/middlewares"
)

func main() {
	r := gin.Default()
	db.ConnectDatabase()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/register", handlers.RegisterUser)
	r.POST("/verify-otp", handlers.VerifyOTP)
	r.POST("/login", handlers.LoginUser)
	r.POST("/verify-forgot-password", handlers.VerifyForgotPassword)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)

	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("/events", handlers.CreateEvent)
		
		auth.GET("/events", handlers.ListEvents)
		auth.GET("/events/:id", handlers.GetEvent)
		auth.POST("/events/:id/book", handlers.BookEvent) // Book event (next step)
	}
	//authenticated.Use(milddlewares.Authenticate)
	//authenticated.POST("/events", createEvents)

	r.Run(":8080")
	//fmt.Println("Project works!")
}
