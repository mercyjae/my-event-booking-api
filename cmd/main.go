package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/routes"
	"github.com/mercyjae/event-booking-api/pkg/mailer"
)

var MailerInstance mailer.Mailer

func main() {
	smtp := mailer.LoadSmtpDetails() // ✅ Note the capital "L" if you want it exported

	// Convert SMTP port from string to int
	port, err := strconv.Atoi(smtp["smtp_port"])
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

	// Initialize your global mailer instance
	MailerInstance = mailer.New(
		smtp["smtp_host"],
		port,
		smtp["smtp_username"],
		smtp["smtp_password"],
		smtp["smtp_sender"],
	)
	r := gin.Default()
	db.InitDB()
	//db.ConnectDatabase()
	routes.UserRoutes(r)

	r.Run(":8070")

}
