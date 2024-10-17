package di

import (
	"os"

	"gopkg.in/gomail.v2"
)

// Send email function using GoMail
func SendNotificationToEmail(email, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "nuhmotp@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("APPEMAIL"), os.Getenv("APPPASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
