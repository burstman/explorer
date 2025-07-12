package mailer

import (
	"os"

	"gopkg.in/gomail.v2"
)

// SendHTML sends an email with both plain text and HTML content.
func SendHTML(to, subject, plainText, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_SENDER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set both plain and HTML versions
	m.SetBody("text/plain", plainText)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	return d.DialAndSend(m)
}
