package mailer

import (
	"os"
	"testing"
)

func TestSend(t *testing.T) {
	// Set Mailtrap credentials for test (normally from env or secrets manager)
	os.Setenv("SMTP_SENDER", "test@hamedapp.com")
	os.Setenv("SMTP_HOST", "sandbox.smtp.mailtrap.io")
	os.Setenv("SMTP_USER", "e175a6d2d36620")     // Replace with real value
	os.Setenv("SMTP_PASSWORD", "791e9b169e5eb3") // Replace with real value

	to := "test@email.com" // Mailtrap test inbox email
	subject := "Test Email"
	body := "This is a test email from Go!"

	err := SendHTML(to, subject, body, body)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
}
