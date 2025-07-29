package types

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// AuthUser represents an user that might be authenticated.
type AuthUser struct {
	UserID    uint
	Email     string
	LoggedIn  bool
	Role      string
	FirstName string
}

// Check should return true if the user is authenticated.
// See handlers/auth.go.
func (user AuthUser) Check() bool {
	return user.UserID > 0 && user.LoggedIn
}

func (user AuthUser) GetRole() string {
	return user.Role
}

func (user AuthUser) GetUserID() uint {
	return user.UserID
}

func (user AuthUser) GetFirstName() string {
	return user.FirstName
}

func (user AuthUser) HasRole(role string) bool {
	return user.Role == role
}

func (user AuthUser) IsAmin() bool {
	return user.Role == "admin"
}

type User struct {
	gorm.Model

	Email           string
	FirstName       string
	LastName        string
	PasswordHash    string
	PhoneNumber     string
	SocialLink      string
	Cin             string
	EmailVerifiedAt sql.NullTime
	Role            string `gorm:"not null;default:user"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
