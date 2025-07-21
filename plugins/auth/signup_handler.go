package auth

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anthdm/superkit/event"
	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
	"github.com/golang-jwt/jwt/v5"
)

var signupSchema = v.Schema{
	"email": v.Rules(v.Email),
	"password": v.Rules(
		v.ContainsSpecial,
		v.ContainsUpper,
		v.Min(7),
		v.Max(50),
	),
	"firstName":          v.Rules(v.Min(2), v.Max(50)),
	"lastName":           v.Rules(v.Min(2), v.Max(50)),
	"phoneNumber":        v.Rules(v.Min(8), v.Max(15)),
	"socialLink":         v.Rules(v.Max(64)),
	"cardIdentityNumber": v.Rules(v.Max(8), v.Min(8), DigitOnly),
}

func HandleSignupIndex(kit *kit.Kit) error {
	return kit.Render(SignupIndex(SignupIndexPageData{}))
}

// HandleSignupCreate handles the user signup process by validating form input, creating a new user,
// generating a verification token, and emitting a signup event. It renders either the signup form
// with validation errors or a confirmation email page upon successful user creation.
func HandleSignupCreate(kit *kit.Kit) error {
	var values SignupFormValues
	errors, ok := v.Request(kit.Request, &values, signupSchema)
	if !ok {
		return kit.Render(SignupForm(values, errors))
	}
	if values.Password != values.PasswordConfirm {
		errors.Add("passwordConfirm", "passwords do not match")
		return kit.Render(SignupForm(values, errors))
	}

	if values.SocialLink != "" && !isValidSocialLink(values.SocialLink) {
		log.Println(values.SocialLink)
		errors.Add("socialLink", "invalid social link")
		return kit.Render(SignupForm(values, errors))
	}

	user, err := createUserFormValues(values)
	if err != nil {
		return err
	}
	token, err := createVerificationToken(user.ID)
	if err != nil {
		return err
	}
	event.Emit(UserSignupEvent, UserWithVerificationToken{
		Token: token,
		User:  user,
	})
	return kit.Render(ConfirmEmail(user))
}

// HandleResendVerificationCode handles the process of resending an email verification token for a user.
// It retrieves the user by ID, checks if the email is already verified, generates a new verification token,
// and emits an event to resend the verification email. Returns an appropriate status message.
func HandleResendVerificationCode(kit *kit.Kit) error {
	idstr := kit.FormValue("userID")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return err
	}

	var user types.User
	if err = db.Get().First(&user, id).Error; err != nil {
		return kit.Text(http.StatusOK, "An unexpected error occured")
	}

	if user.EmailVerifiedAt.Time.After(time.Time{}) {
		return kit.Text(http.StatusOK, "Email already verified!")
	}

	token, err := createVerificationToken(uint(id))
	if err != nil {
		return kit.Text(http.StatusOK, "An unexpected error occured")
	}

	event.Emit(ResendVerificationEvent, UserWithVerificationToken{
		User:  user,
		Token: token,
	})

	msg := fmt.Sprintf("A new verification token has been sent to %s", user.Email)

	return kit.Text(http.StatusOK, msg)
}

// createVerificationToken generates a JWT token for email verification with a configurable expiry time.
// The token contains the user ID as the subject and expires after a specified number of hours (default 1).
// It uses the SUPERKIT_SECRET environment variable for signing the token.
// Returns the signed JWT token string and any error encountered during token creation.
func createVerificationToken(userID uint) (string, error) {
	expiryStr := kit.Getenv("SUPERKIT_AUTH_EMAIL_VERIFICATION_EXPIRY_IN_HOURS", "1")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {
		expiry = 1
	}

	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiry))),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SUPERKIT_SECRET")))
}
