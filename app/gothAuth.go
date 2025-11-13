package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
)

func InitializeGoth() {
	sessionSecret := os.Getenv("SUPERKIT_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET is required")
	}

	gothic.Store = sessions.NewCookieStore([]byte(sessionSecret))

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"), "email", "profile", "openid",
		),
		facebook.New(
			os.Getenv("FACEBOOK_KEY"),
			os.Getenv("FACEBOOK_SECRET"),
			os.Getenv("FACEBOOK_CALLBACK_URL"),
			"email", "public_profile",
		),
	)

	gothic.GetProviderName = func(r *http.Request) (string, error) {
		segments := strings.Split(r.URL.Path, "/")
		if len(segments) >= 3 {
			return segments[2], nil
		}
		return "", fmt.Errorf("no provider specified")
	}

}
