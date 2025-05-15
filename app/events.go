package app

import (
	"camping/app/events"
	"camping/plugins/auth"
	"log"

	"github.com/anthdm/superkit/event"
)

// Events are functions that are handled in separate goroutines.
// They are the perfect fit for offloading work in your handlers
// that otherwise would take up response time.
// - sending email
// - sending notifications (Slack, Telegram, Discord)
// - analytics..

// Register your events here.
func RegisterEvents() {
	event.Subscribe(auth.UserSignupEvent, events.OnUserSignup)
	event.Subscribe(auth.ResendVerificationEvent, events.OnResendVerificationToken)
}

func OnUserSignup(payload any) {
	data, ok := payload.(auth.UserWithVerificationToken)
	if !ok {
		log.Println("invalid payload for UserSignupEvent")
		return
	}

	err := auth.SendVerificationEmail(data.User.Email, data.Token)
	if err != nil {
		log.Printf("failed to send verification email to %s: %v", data.User.Email, err)
	}
}
