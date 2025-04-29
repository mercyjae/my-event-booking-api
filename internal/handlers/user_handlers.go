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
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.RegisterUser
	result := db.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user found with that email"})
		return
	}

	otp := utils.GenerateOTP()
	user.OTP = otp
	user.OTPExpiresAt = time.Now().Add(10 * time.Minute)
	db.DB.Save(&user)

	utils.SendEmail(user.Email, "Your Password Reset OTP", "Your OTP code is: "+otp)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset OTP sent to your email."})

}

func VerifyForgotPassword(c *gin.Context) {
	var req models.VerifyOTP

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.RegisterUser
	result := db.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
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

	user.OTP = ""
	user.OTPExpiresAt = time.Time{}

	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified. You can now reset your password."})

}

func ResetPassword(c *gin.Context) {

	var req models.ResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	var user models.RegisterUser
	result := db.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	user.Password = string(hashedPassword)

	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful!"})
}

func GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDInterface.(float64))

	var user models.RegisterUser
	result := db.DB.First(&user, userID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"full_name":     user.FullName,
		"email":         user.Email,
		"phone":         user.Phone,
		"date_of_birth": user.DoB,
		"verified":      user.Verified,
		"created_at":    user.CreatedAt,
	})
}

func EditProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDInterface.(float64))

	var user models.RegisterUser
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
		DoB      string `json:"date_of_birth"` // in ISO format e.g., "2000-01-01"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.DoB != "" {
		dob, err := time.Parse("2006-01-02", req.DoB)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		user.DoB = dob
	}

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "profile": gin.H{
		"full_name":     user.FullName,
		"phone":         user.Phone,
		"date_of_birth": user.DoB,
	}})

}

func ChangePassword(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDInterface.(float64))

	var user models.RegisterUser
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	if req.NewPassword == "" || req.ConfirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirmation are required"})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
