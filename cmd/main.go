package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/routes"
)

func main() {
	r := gin.Default()
	db.ConnectDatabase()
	routes.UserRoutes(r)

	r.Run(":8080")

}
