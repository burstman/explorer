package main

import (
	"database/sql"
	"explorer/app/db"
	"explorer/app/types"
	"log"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
	database := db.Get() // or your actual DB connection method

	passwordAdmin := "admin1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordAdmin), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	now := time.Now()
	admin := types.User{
		Email:        "admin@camping.tn",
		PasswordHash: string(hashedPassword),
		FirstName:    "Hamrouni",
		LastName:     "Foued",
		Role:         "admin",
		EmailVerifiedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create the user

	user := types.User{
		Email:        "user@camping.tn",
		PasswordHash: string(hashedPassword),
		FirstName:    "Flissi",
		LastName:     "Hamed",
		Role:         "user",
		EmailVerifiedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := database.Create(&admin)
	if result.Error != nil {
		log.Fatalf("failed to create admin: %v", result.Error)
	}
	result = database.Create(&user)
	if result.Error != nil {
		log.Fatalf("failed to create admin: %v", result.Error)
	}

	log.Println("✅ Admin user seeded successfully.")
}
