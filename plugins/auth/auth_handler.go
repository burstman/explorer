package auth

import (
	"database/sql"
	"explorer/app/db"
	"explorer/app/handlers"

	"explorer/app/types"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	userSessionName = "user-session"
)

var authSchema = v.Schema{
	"email":    v.Rules(v.Email),
	"password": v.Rules(v.Required),
}

func HandleLoginIndex(kit *kit.Kit) error {
	if kit.Auth().Check() {
		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")
		return kit.Redirect(http.StatusSeeOther, redirectURL)
	}
	return handlers.RenderWithLayout(kit, LoginIndex(LoginIndexPageData{}))
}

// HandleLoginCreate processes user login by validating credentials, checking email verification,
// creating a session, and redirecting the user to a specified URL upon successful authentication.
// It handles form validation, user lookup, password verification, and optional email verification.
// Returns an error if login fails or session creation encounters issues.
func HandleLoginCreate(kit *kit.Kit) error {
	var values LoginFormValues
	errors, ok := v.Request(kit.Request, &values, authSchema)
	if !ok {
		return kit.Render(LoginForm(values, errors))
	}

	var user types.User
	err := db.Get().Find(&user, "email = ?", values.Email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.Add("credentials", "invalid credentials")
			return kit.Render(LoginForm(values, errors))
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(values.Password))
	if err != nil {
		errors.Add("credentials", "invalid credentials")
		return kit.Render(LoginForm(values, errors))
	}

	skipVerify := kit.Getenv("SUPERKIT_AUTH_SKIP_VERIFY", "false")
	if skipVerify != "true" {
		if !user.EmailVerifiedAt.Valid {
			errors.Add("verified", "please verify your email")
			return kit.Render(LoginForm(values, errors))
		}
	}

	sessionExpiryStr := kit.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS", "48")
	sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
	if err != nil {
		sessionExpiry = 48
	}
	session := Session{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(sessionExpiry)),
	}
	if err = db.Get().Create(&session).Error; err != nil {
		return err
	}

	sess := kit.GetSession(userSessionName)
	sess.Values["sessionToken"] = session.Token
	sess.Save(kit.Request, kit.Response)
	redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/AreaAttraction")

	return kit.Redirect(http.StatusSeeOther, redirectURL)
}

func HandleLoginDelete(kit *kit.Kit) error {
	sess := kit.GetSession(userSessionName)
	defer func() {
		sess.Values = map[any]any{}
		sess.Save(kit.Request, kit.Response)
	}()
	err := db.Get().Delete(&Session{}, "token = ?", sess.Values["sessionToken"]).Error
	if err != nil {
		return err
	}
	return kit.Redirect(http.StatusSeeOther, "/")
}

// HandleEmailVerify processes an email verification request by validating a JWT token,
// verifying the user's email, and redirecting to the login page upon successful verification.
// It handles various error scenarios such as invalid tokens, expired tokens, or user lookup failures.
func HandleEmailVerify(kit *kit.Kit) error {
	tokenStr := kit.Request.URL.Query().Get("token")
	if tokenStr == "" {
		return kit.Render(EmailVerificationError("Invalid or missing token"))
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SUPERKIT_SECRET")), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil || !token.Valid {
		log.Println("Invalid token", "error", err)
		return kit.Render(EmailVerificationError("Invalid verification token"))
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return kit.Render(EmailVerificationError("Verification token expired"))
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return kit.Render(EmailVerificationError("Invalid user ID"))
	}

	var user types.User
	err = db.Get().First(&user, userID).Error
	if err != nil {
		return kit.Render(EmailVerificationError("User not found"))
	}

	if user.EmailVerifiedAt.Valid {
		return kit.Render(EmailVerificationInfo("Your email is already verified. You can log in now."))
	}

	user.EmailVerifiedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := db.Get().Save(&user).Error; err != nil {
		return kit.Render(EmailVerificationError("Failed to verify email. Please try again."))
	}

	return kit.Redirect(http.StatusSeeOther, "/login")
}

func AuthenticateUser(kit *kit.Kit) (kit.Auth, error) {
	auth := Auth{}
	sess := kit.GetSession(userSessionName)
	token, ok := sess.Values["sessionToken"]
	//fmt.Println("token", token)
	if !ok {
		return auth, nil
	}

	var session Session
	err := db.Get().
		Preload("User").
		Find(&session, "token = ? AND expires_at > ?", token, time.Now()).Error
	if err != nil || session.ID == 0 {
		return auth, nil
	}
	//fmt.Println(session.User.Role)

	return types.AuthUser{
		LoggedIn: true,
		UserID:   session.User.ID,
		Email:    session.User.Email,
		Role:     session.User.Role,
	}, nil
}
