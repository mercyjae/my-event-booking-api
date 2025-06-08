package utils

import (
	"os"

	"github.com/go-gomail/gomail"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USERNAME"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_SENDER"), os.Getenv("SMTP_PASSWORD"))

	return d.DialAndSend(m)
}

// smtpHost := os.Getenv("SMTP_HOST")
// smtpPort := os.Getenv("SMTP_PORT")
// smtpUsername := os.Getenv("SMTP_USERNAME")
// smtpPassword := os.Getenv("SMTP_PASSWORD")
// smtpSender := os.Getenv("SMTP_SENDER")
