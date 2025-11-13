package auth

import (
	"explorer/app/mailer"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	v "github.com/anthdm/superkit/validate"
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
	return mailer.SendHTML(email, subject, body, body)
}

func isValidSocialLink(link string) bool {
	return strings.HasPrefix(link, "https://www.facebook.com/") ||
		strings.HasPrefix(link, "https://facebook.com/") ||
		strings.HasPrefix(link, "https://www.instagram.com/") ||
		strings.HasPrefix(link, "https://instagram.com/")
}

var DigitOnly = v.RuleSet{
	Name: "digitOnly",
	ValidateFunc: func(rule v.RuleSet) bool {
		str, ok := rule.FieldValue.(string)
		if !ok {
			return false
		}
		for _, ch := range str {
			if !unicode.IsDigit(ch) {
				return false
			}
		}
		return true
	},
	MessageFunc: func(set v.RuleSet) string {
		return "must be digits"
	},
}
