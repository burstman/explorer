package auth

import (
	"database/sql"
	"explorer/app/db"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Event name constants
const (
	UserSignupEvent         = "auth.signup"
	ResendVerificationEvent = "auth.resend.verification"
)

// UserWithVerificationToken is a struct that will be sent over the
// auth.signup event. It holds the User struct and the Verification token string.
type UserWithVerificationToken struct {
	User  User
	Token string
}

type Auth struct {
	UserID    uint
	Email     string
	LoggedIn  bool
	FirstName string
	Role string
}

func (auth Auth) Check() bool {
	return auth.LoggedIn
}

func (user Auth) HasRole(role string) bool {
	return user.Role == role
}

func (user Auth) IsAmin() bool {
	return user.Role == "admin"
}



type User struct {
	gorm.Model

	Email           string
	FirstName       string
	LastName        string
	PasswordHash    string
	EmailVerifiedAt sql.NullTime
	Role            string `gorm:"not null;default:user"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// createUserFromFormValues creates a new User from signup form values, hashing the password
// and persisting the user record to the database. It returns the created User and any error
// encountered during password hashing or database insertion.
func createUserFromFormValues(values SignupFormValues) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(values.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		Email:        values.Email,
		FirstName:    values.FirstName,
		LastName:     values.LastName,
		Role:         "user",
		PasswordHash: string(hash),
	}
	result := db.Get().Create(&user)
	return user, result.Error
}

type Session struct {
	gorm.Model

	UserID    uint
	Token     string
	IPAddress string
	UserAgent string
	ExpiresAt time.Time
	CreatedAt time.Time
	User      User
}

