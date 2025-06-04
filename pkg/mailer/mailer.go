package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"

	"time"

	"github.com/joho/godotenv"
	"gopkg.in/mail.v2"
)

// this directive is important for embed
//

var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {

	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{dialer: dialer, sender: sender}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	//logger := logger.GetLogger(logger.Options{})
	// defer logger.Sync()
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		if nil == err {
			fmt.Println("Email Sent", nil)
			//	logger.Info("Email Sent", nil)
			return nil
		}
		//logger.Error(err.Error(), nil)
		time.Sleep(500 * time.Millisecond)
	}
	return err

}

func LoadSmtpDetails() map[string]string {

	godotenv.Load()
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpSender := os.Getenv("SMTP_SENDER")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || smtpSender == "" {
		//log.Fatal("Couldn't load smtp details", nil)
	}
	smtpMap := map[string]string{
		"smtp_host":     smtpHost,
		"smtp_port":     smtpPort,
		"smtp_username": smtpUsername,
		"smtp_password": smtpPassword,
		"smtp_sender":   smtpSender,
	}
	return smtpMap
}
