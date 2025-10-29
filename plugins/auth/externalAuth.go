package auth

import (
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
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

func CombinedAuthHandler(kit *kit.Kit) error {
	path := kit.Request.URL.Path
	log.Printf("OAuth path: %s", path)

	if strings.HasSuffix(path, "/callback") {
		user, err := gothic.CompleteUserAuth(kit.Response, kit.Request)
		if err != nil {
			log.Printf("Auth error: %v", err)

			return fmt.Errorf("authentication failed: %v", err)

		}

		// ✅ Store in Superkit session (optional)
		sess := kit.GetSession("user-session")
		sess.Values["name"] = user.Name
		sess.Values["email"] = user.Email
		sess.Values["provider"] = user.Provider
		sess.Values["social_id"] = user.UserID
		sess.Save(kit.Request, kit.Response)

		dbUser, err := FindOrCreateUserByEmail(user.Email, &user.FirstName, &user.LastName, &user.UserID, &user.Provider)
		if err != nil {
			log.Printf("Error finding or creating user: %v", err)
			return fmt.Errorf("failed to find or create user: %v", err)
		}
		dbUser.FirstName = user.FirstName // if available
		dbUser.LastName = user.LastName   // if available
		dbUser.SocialID = &user.UserID
		dbUser.Provider = &user.Provider

		log.Printf("User info from provider: %+v", user)

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

// FindOrCreateUserByEmail looks up a user by email or creates/updates one
func FindOrCreateUserByEmail(email string, firstName, lastName, socialID, provider *string) (*types.User, error) {
	var user types.User

	err := db.Get().Where("email = ?", email).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		// Create new user
		user = types.User{
			Email:     email,
			FirstName: *firstName,
			LastName:  *lastName,
			SocialID:  socialID,
			Provider:  provider,
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Get().Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	// Update existing user with new social info if missing
	updated := false
	if firstName != nil && user.FirstName == "" {
		user.FirstName = *firstName
		updated = true
	}
	if lastName != nil && (user.LastName == "") {
		user.LastName = *lastName
		updated = true
	}
	if socialID != nil && (user.SocialID == nil || *user.SocialID == "") {
		user.SocialID = socialID
		updated = true
	}
	if provider != nil && (user.Provider == nil || *user.Provider == "") {
		user.Provider = provider
		updated = true
	}

	if updated {
		user.UpdatedAt = time.Now()
		if err := db.Get().Save(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}
