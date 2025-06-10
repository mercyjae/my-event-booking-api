package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/handlers"
	"github.com/mercyjae/event-booking-api/internal/middlewares"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/alive", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to my event booking app!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/register", handlers.RegisterUser)
	r.POST("/verify-otp", handlers.VerifyOTP)
	r.POST("/login", handlers.LoginUser)
	r.POST("/verify-forgot-password", handlers.VerifyForgotPassword)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)
	r.GET("/users", handlers.ListUsers)

	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("/events", handlers.CreateEvent)
		auth.GET("/events", handlers.ListEvents)
		auth.GET("/events/:id", handlers.GetEvent)
		auth.DELETE("/events/:id", handlers.DeleteEvent)
		auth.PUT("/events/:id", handlers.UpdateEvent)
		auth.POST("/events/:id/book", handlers.BookEvent)
		auth.DELETE("/booking/:id/cancel", handlers.CancelBooking)
		auth.GET("/bookings", handlers.GetBookings)
		auth.GET("/profile", handlers.GetProfile)
		auth.PUT("/profile/edit", handlers.EditProfile)
		auth.PUT("/change-password", handlers.ChangePassword)
	}

}
