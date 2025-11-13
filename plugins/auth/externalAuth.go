package auth

import (
	"errors"
	"explorer/app/db"
	"explorer/app/types"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anthdm/superkit/kit"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

func CombinedAuthHandler(kit *kit.Kit) error {
	path := kit.Request.URL.Path

	if strings.HasSuffix(path, "/callback") {
		user, err := gothic.CompleteUserAuth(kit.Response, kit.Request)
		if err != nil {
			log.Printf("Auth error: %v", err)
			return fmt.Errorf("authentication failed: %v", err)
		}

		// ✅ Check if a user already exists
		var dbUser *types.User

		// Try by email first (if available)
		if user.Email != "" {
			dbUser, err = FindUserByEmail(user.Email)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("database error in CombinedAuthHandler (email lookup): %v", err)
			}
		}

		// If not found by email, or email is empty — try by provider ID
		if dbUser == nil {
			dbUser, err = FindUserBySocialID(user.UserID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("database error in CombinedAuthHandler (ID lookup): %v", err)
			}
		}

		// ✅ If found, check provider consistency
		if dbUser != nil {
			if dbUser.Provider != nil && *dbUser.Provider != user.Provider {
				sess := kit.GetSession("user-session")
				sess.Values["flash"] = fmt.Sprintf(
					"You already registered using %s. Please log in with %s.",
					*dbUser.Provider, *dbUser.Provider,
				)
				sess.Save(kit.Request, kit.Response)
				return kit.Redirect(http.StatusSeeOther, "/login")
			}
		} else {
			// ✅ New user → create
			dbUser, err = CreateUserFromOAuth(user)
			if err != nil {
				return fmt.Errorf("failed to create user: %v", err)
			}
		}

		// ✅ Update latest info (safe even if exists)
		dbUser.FirstName = user.FirstName
		dbUser.LastName = user.LastName
		dbUser.SocialID = &user.UserID
		dbUser.Provider = &user.Provider
		dbUser.Email = user.Email
		if err := db.Get().Save(dbUser).Error; err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}

		// ✅ Create session
		sessionExpiryStr := kit.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS", "48")
		sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
		if err != nil {
			sessionExpiry = 48
		}

		session := Session{
			UserID:    dbUser.ID,
			Token:     uuid.New().String(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(sessionExpiry)),
		}
		if err = db.Get().Create(&session).Error; err != nil {
			return err
		}

		// ✅ Store user info in cookie session
		sess := kit.GetSession("user-session")
		sess.Values["sessionToken"] = session.Token
		sess.Values["email"] = dbUser.Email
		sess.Values["firstName"] = dbUser.FirstName
		sess.Values["lastName"] = dbUser.LastName
		sess.Values["phone"] = dbUser.PhoneNumber
		sess.Values["provider"] = dbUser.Provider
		sess.Values["social_id"] = dbUser.SocialID
		sess.Save(kit.Request, kit.Response)

		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/AreaAttraction")
		return kit.Redirect(http.StatusSeeOther, redirectURL)
	}

	// If not callback, start auth
	gothic.BeginAuthHandler(kit.Response, kit.Request)
	return nil
}

// ✅ Better naming for clarity
func FindUserByEmail(email string) (*types.User, error) {
	var user types.User
	if err := db.Get().Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ✅ Renamed for accuracy (you were looking by social ID)
func FindUserBySocialID(id string) (*types.User, error) {
	var user types.User
	if err := db.Get().Where("social_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUserFromOAuth(user goth.User) (*types.User, error) {
	if user.Email == "" {
		// generate a pseudo-unique placeholder, or set to NULL if allowed
		log.Printf("Warning: OAuth user has no email, generating placeholder for UserID %s", user.UserID)
		user.Email = fmt.Sprintf("noemail_%s@%s.local", user.UserID, user.Provider)
	}

	dbUser := types.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		SocialID:  &user.UserID,
		Provider:  &user.Provider,
	}
	if err := db.Get().Create(&dbUser).Error; err != nil {
		return nil, fmt.Errorf("error in CreateUserFromOAuth: %v", err)
	}
	return &dbUser, nil
}
