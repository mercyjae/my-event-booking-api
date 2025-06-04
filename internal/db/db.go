package db

import (
	"log"
	"github.com/mercyjae/event-booking-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := "host=localhost user=mac dbname=eventdb port=5432 sslmode=disable"

	//"host=localhost user=youruser password=yourpassword dbname=eventdb port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	DB = database

	database.AutoMigrate(&models.RegisterUser{}, &models.Event{}, &models.Booking{}, &models.VerifyOTP{}, &models.ResetPassword{})
 }
