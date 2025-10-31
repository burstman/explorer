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
	//log.Printf("OAuth path: %s", path)

	if strings.HasSuffix(path, "/callback") {
		user, err := gothic.CompleteUserAuth(kit.Response, kit.Request)
		if err != nil {
			log.Printf("Auth error: %v", err)

			return fmt.Errorf("authentication failed: %v", err)

		}

		// ✅ Check if a user already exists
		dbUser, err := FindUserByEmail(user.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("database error: %v", err)
		}

		if dbUser != nil {
			// ✅ User exists → check provider consistency
			if dbUser.Provider != nil && *dbUser.Provider != user.Provider {
				// The user registered with another provider
				log.Printf("Provider mismatch: existing=%s, tried=%s", *dbUser.Provider, user.Provider)
				// ✅ Add flash message
				sess := kit.GetSession("user-session")
				sess.Values["flash"] = fmt.Sprintf("You already registered using %s. Please log in with %s.", *dbUser.Provider, *dbUser.Provider)
				sess.Save(kit.Request, kit.Response)

				return kit.Redirect(http.StatusSeeOther, "/login")

			}
		} else {
			// ✅ New user → allow creation
			dbUser, err = CreateUserFromOAuth(user)
			if err != nil {
				return fmt.Errorf("failed to create user: %v", err)
			}
		}

		// ✅ Store in Superkit session (optional)
		sess := kit.GetSession("user-session")
		sess.Values["name"] = user.Name
		sess.Values["email"] = user.Email
		sess.Values["provider"] = user.Provider
		sess.Values["social_id"] = user.UserID
		sess.Save(kit.Request, kit.Response)

		dbUser.FirstName = user.FirstName // if available
		dbUser.LastName = user.LastName   // if available
		dbUser.SocialID = &user.UserID
		dbUser.Provider = &user.Provider

		//log.Printf("User info from provider: %+v", user)

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

		sess.Values["sessionToken"] = session.Token
		sess.Save(kit.Request, kit.Response)
		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/AreaAttraction")

		return kit.Redirect(http.StatusSeeOther, redirectURL)

	}

	gothic.BeginAuthHandler(kit.Response, kit.Request)
	return nil
}

func FindUserByEmail(email string) (*types.User, error) {
	var user types.User
	if err := db.Get().Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUserFromOAuth(user goth.User) (*types.User, error) {
	dbUser := types.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		SocialID:  &user.UserID,
		Provider:  &user.Provider,
	}
	if err := db.Get().Create(&dbUser).Error; err != nil {
		return nil, err
	}
	return &dbUser, nil
}
