package mailer

import (
	"os"

	"gopkg.in/gomail.v2"
)

func Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_SENDER")) // e.g. "test@yourapp.com"
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),     // e.g. smtp.mailtrap.io
		587,                        // Port
		os.Getenv("SMTP_USER"),     // username from Mailtrap
		os.Getenv("SMTP_PASSWORD"), // password from Mailtrap
	)

	return d.DialAndSend(m)
}
