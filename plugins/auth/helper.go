package auth

import (
	"camping/app/mailer"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateEmailVerificationToken creates a signed JWT with the user ID as subject
func GenerateEmailVerificationToken(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("SUPERKIT_SECRET")
	return token.SignedString([]byte(secret))
}

func SendVerificationEmail(email, token string) error {
	url := fmt.Sprintf("%s/verify?token=%s", os.Getenv("APP_BASE_URL"), token)
	subject := "Verify your email address"
	body := fmt.Sprintf("Welcome!\n\nPlease verify your email by clicking this link:\n\n%s", url)
	return mailer.Send(email, subject, body)
}
