package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/models"
	"github.com/mercyjae/event-booking-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var req models.RegisterUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	otp := utils.GenerateOTP()

	user := models.RegisterUser{
		FullName:     req.FullName,
		Email:        req.Email,
		Password:     string(hashedPassword),
		Verified:     false,
		OTP:          otp,
		OTPExpiresAt: time.Now().Add(10 * time.Minute),
	}
	db.DB.Create(&user)
	utils.SendEmail(user.Email, "Verify your email", "Your OTP code is: "+otp)
	c.JSON(http.StatusCreated, gin.H{"message": "Registration success"})
}

func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTP
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.RegisterUser
	result := db.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	if time.Now().After(user.OTPExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
		return
	}

	user.Verified = true
	user.OTP = ""
	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully!"})
}

func LoginUser(c *gin.Context) {

	var req models.LoginUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.RegisterUser
	result := db.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})


	
}
