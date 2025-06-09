package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mercyjae/event-booking-api/internal/db"
	"github.com/mercyjae/event-booking-api/internal/domain"
	"github.com/mercyjae/event-booking-api/internal/dto"
	"github.com/mercyjae/event-booking-api/internal/repo"
	"github.com/mercyjae/event-booking-api/pkg/mailer"
	"github.com/mercyjae/event-booking-api/pkg/utils"
)

func RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	//	otp := utils.GenerateOTP()

	user := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		//Password: string(hashedPassword),
		// Verified:     false,
		// OTP:          otp,
		// OTPExpiresAt: time.Now().Add(10 * time.Minute),
	}
	exists, err := repo.IsEmailTaken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while checking email"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}
	err = user.Password.Set(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err, "devError": err.Error()})
		return
	}
	err = repo.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user", "devError": err.Error()})
		return
	}
	//db.DB.Create(&user)
	// data := map[string]any{
	// 	"name":            req.FullName,
	// 	"expiryDate":      user.OTPExpiresAt.Format("Monday, 02 January 2006 at 15:04"),
	// 	"activationToken": user.OTP}

	// // err := app.Mailer.Send("dolagookun@icloud.com", "reset-token.html", data)
	// err := MailerInstance.Send(user.Email, "reset-token.html", data)
	// //	err := mailer.Mailer.Send(user.Email, "reset-token.html", data)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	//	app.Logger.Error(err.Error(), nil)
	// }
	//utils.SendEmail(user.Email, "Verify your email", "Your OTP code is: "+otp)
	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// func VerifyOTP(c *gin.Context) {
// 	var req models.VerifyOTP
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user models.RegisterUser
// 	result := db.DB.Where("email = ?", req.Email).First(&user)

// 	if result.Error != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
// 		return
// 	}

// 	if user.OTP != req.OTP {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
// 		return
// 	}

// 	if time.Now().After(user.OTPExpiresAt) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
// 		return
// 	}

// 	user.Verified = true
// 	user.OTP = ""
// 	db.DB.Save(&user)

// 	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully!"})
// }

func LoginUser(c *gin.Context) {

	var req dto.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//var user dto.RegisterUserRequest
	user, err := repo.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}
	match, err := user.Password.Matches(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password validation failed"})
		return
	}
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	//result := db.DBB.Where("email = ?", req.Email).First(&user)

	// if !user.Verified {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified"})
	// 	return
	// }
	//match, err := user.Password.Matches(req.Password)
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	// 	return
	// }

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

	var user struct {
		ID       int
		FullName string
		Email    string
	}
	// Fetch user by email
	query := "SELECT id, full_name, email FROM users WHERE email = ?"
	err := db.DBB.QueryRow(query, req.Email).Scan(&user.ID, &user.FullName, &user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user found with that email"})
		return
	}

	// Generate OTP and expiry
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(10 * time.Minute)

	//Update OTP in the database
	updateQuery := "UPDATE users SET otp = ?, otp_expires_at = ? WHERE id = ?"
	_, err = db.DBB.Exec(updateQuery, otp, otpExpiry, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP", "devError": err.Error()})
		return
	}

	// Send OTP email
	data := map[string]any{
		"name":            user.FullName,
		"expiryDate":      otpExpiry,
		"activationToken": otp,
	}
	mailerService := mailer.Newi()
	err = mailerService.Send(user.Email, "token.html", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset OTP sent to your email."})
}

func VerifyForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var passwordHash string
	// Query user from DB
	var user domain.User
	query := `SELECT id, full_name, email, password, otp, otp_expires_at FROM users WHERE email = ?`
	err := db.DBB.QueryRow(query, req.Email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&passwordHash,
		&user.OTP,
		&user.OTPExpiresAt,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email", "devError": err.Error()})
		return
	}
	user.Password = domain.Password{
		Hash: []byte(passwordHash),
	}
	// Validate OTP
	if user.OTP != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}
	if time.Now().After(user.OTPExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
		return
	}

	// Clear OTP fields
	updateQuery := `UPDATE users SET otp = NULL, otp_expires_at = NULL WHERE id = ?`
	_, err = db.DBB.Exec(updateQuery, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear OTP", "devError": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified. You can now reset your password."})
}

func ResetPassword(c *gin.Context) {
	var req dto.ResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	err := repo.ResetPasswordByEmail(req.Email, req.NewPassword)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password", "devError": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful!"})
}
