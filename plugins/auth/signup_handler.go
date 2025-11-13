package auth

import (
	"explorer/app/db"
	"explorer/app/handlers"
	"explorer/app/types"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anthdm/superkit/event"
	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var passwordRules = v.Rules(
	v.ContainsUpper,
	v.Min(7),
	v.Max(50),
)

var signupSchema = v.Schema{
	"email":              v.Rules(v.Email),
	"password":           passwordRules,
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

	// if values.SocialLink != "" && !isValidSocialLink(values.SocialLink) {

	// 	errors.Add("socialLink", "invalid social link")
	// 	return kit.Render(SignupForm(values, errors))
	// }

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

	//msg := fmt.Sprintf("A new verification token has been sent to %s", user.Email)

	return kit.Render(ConfirmEmail(user))
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

func HandelForgotPasswordPage(kit *kit.Kit) error {
	return kit.Render(ForgotPasswordPage(nil, nil))

}

func HandelResetPassEmailSend(kit *kit.Kit) error {
	email := kit.Request.FormValue("email")

	// check if user exists

	var user types.User
	if err := db.Get().Where("email = ?", email).First(&user).Error; err != nil {
		return kit.Text(http.StatusOK, `<div class="p-6 text-center">A reset link has been sent.</div>`)
	}

	// generate reset token
	token, err := createVerificationToken(user.ID)
	if err != nil {
		return kit.Text(http.StatusInternalServerError, "An unexpected error occurred while creating reset token")

	}

	// Emit event (async email sending)
	event.Emit(PasswordResetEvent, UserWithResetToken{
		User:  user,
		Token: token,
	})

	return kit.Text(http.StatusOK, `<div class="p-6 text-center">A reset link has been sent.</div>`)

}

type ResetPasswordForm struct {
	Password string `form:"password"`
}

func HandelResetPass(kit *kit.Kit) error {
	token := kit.Request.URL.Query().Get("token")
	if token == "" {
		return kit.Text(http.StatusBadRequest, "Invalid email link")
	}

	switch kit.Request.Method {
	case http.MethodGet:
		// Render form
		return handlers.RenderWithLayout(kit, ResetPasswordPage(token, nil))

	case http.MethodPost:
		// Parse and validate token
		claims := &jwt.RegisteredClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SUPERKIT_SECRET")), nil
		})
		if err != nil || !parsedToken.Valid {
			return kit.Text(http.StatusBadRequest, "Invalid or expired token")
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return kit.Text(http.StatusBadRequest, "Invalid token data")
		}
		password := kit.FormValue("password")
		confirm := kit.FormValue("confirm_password")

		if password == "" || confirm == "" || password != confirm {
			return kit.Render(ResetPasswordPage(token, map[string]string{
				"password": "Passwords do not match",
			}))
		}

		form := ResetPasswordForm{Password: password}

		errs, okPass := v.Validate(&form, v.Schema{
			"password": passwordRules,
		})

		if !okPass {
			return kit.Render(ResetPasswordPage(token, map[string]string{
				"password": strings.Join(errs["password"], ", "),
			}))
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err := db.Get().Model(&types.User{}).Where("id = ?", userID).Update("password_hash", string(hashed)).Error; err != nil {
			return kit.Text(http.StatusInternalServerError, "Failed to update password")
		}

		return kit.Text(http.StatusOK, `<div class="p-6 text-center">Your password has been reset successfully. You may now log in.</div>`)

	default:
		return kit.Text(http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func FacebookLogin(kit *kit.Kit) error {

	return kit.Text(http.StatusOK, "Facebook Login")
}

func FacebookCallback(kit *kit.Kit) error {

	return kit.Text(http.StatusOK, "Facebook Callback")
}
