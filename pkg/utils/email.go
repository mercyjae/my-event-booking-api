package utils

import "github.com/go-gomail/gomail"

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "youremail@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "youremail@gmail.com", "yourapppassword")

	return d.DialAndSend(m)
}
