package events

import (
	"context"
	"explorer/app/mailer"
	"explorer/plugins/auth"
	"os"

	"fmt"
	"log"
)

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// }

var ngrokURL = os.Getenv("NGROK_ADDRESS")

// Event handlers
func OnUserSignup(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithVerificationToken)
	if !ok {
		return
	}
	//userWithToken, _ := json.MarshalIndent(userWithToken, "   ", "    ")
	//fmt.Println(string(b))
	subject := "Verify your email"
	link := fmt.Sprintf(ngrokURL+"/email-verify?token=%s", userWithToken.Token)
	htmlLink := fmt.Sprintf(`<a href="%s">Verify your email</a>`, link)
	plainText := fmt.Sprintf("Hi %s,\n\nThanks for signing up!\nPlease verify your email by clicking this link: %s\n\nIf you didn’t sign up, please ignore this email.\n\nThanks,\nExplorer Team", userWithToken.User.FirstName, link)
	htmlText := fmt.Sprintf(`
    <p>Hi %s,</p>
    <p>Thanks for signing up!</p>
    <p>Please verify your email by clicking the link below:</p>
    <p>%s</p>
    <p>If you didn’t sign up, please ignore this email.</p>
    <p>Thanks,<br/>The Explorer Team</p>`, userWithToken.User.FirstName, htmlLink)
	error := mailer.SendHTML(userWithToken.User.Email, subject, plainText, htmlText)
	if error != nil {
		log.Printf("Error sending email: %v", error)
	}

}

// OnResendVerificationToken handles the resend event by re-sending the verification email.
func OnResendVerificationToken(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithVerificationToken)
	if !ok {
		return
	}

	// 	b, _ := json.MarshalIndent(userWithToken, "   ", "    ")

	subject := "Verify your email"
	link := fmt.Sprintf(ngrokURL+"/resend-email-verification?token=%s", userWithToken.Token)
	htmlLink := fmt.Sprintf(`<a href="%s">Verify your email</a>`, link)
	plainText := fmt.Sprintf("Hi %s,\n\nThanks for signing up!\nPlease verify your email by clicking this link: %s\n\nIf you didn’t sign up, please ignore this email.\n\nThanks,\nExplorer Team", userWithToken.User.FirstName, link)
	htmlText := fmt.Sprintf(`
    <p>Hi %s,</p>
    <p>Thanks for signing up!</p>
    <p>Please verify your email by clicking the link below:</p>
    <p>%s</p>
    <p>If you didn’t sign up, please ignore this email.</p>
    <p>Thanks,<br/>The Explorer Team</p>`, userWithToken.User.FirstName, htmlLink)
	error := mailer.SendHTML(userWithToken.User.Email, subject, plainText, htmlText)
	if error != nil {
		log.Printf("Error sending email: %v", error)
	}

}

func OnPasswordReset(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithResetToken)
	if !ok {
		log.Println("invalid payload for PasswordResetEvent")
		return
	}

	subject := "Reset your password"
	link := fmt.Sprintf(ngrokURL+"/resetPass?token=%s", userWithToken.Token)
	htmlLink := fmt.Sprintf(`<a href="%s">Reset your password</a>`, link)

	plainText := fmt.Sprintf(
		"Hi %s,\n\nWe received a request to reset your password.\nPlease click the link below to reset it:\n%s\n\nIf you did not request a password reset, please ignore this email.\n\nThanks,\nExplorer Team",
		userWithToken.User.FirstName,
		link,
	)

	htmlText := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>We received a request to reset your password.</p>
		<p>Please reset it by clicking the link below:</p>
		<p>%s</p>
		<p>If you did not request this, please ignore this email.</p>
		<p>Thanks,<br/>The Explorer Team</p>`,
		userWithToken.User.FirstName,
		htmlLink,
	)

	err := mailer.SendHTML(userWithToken.User.Email, subject, plainText, htmlText)
	if err != nil {
		log.Printf("Error sending password reset email to %s: %v", userWithToken.User.Email, err)
	}
}
