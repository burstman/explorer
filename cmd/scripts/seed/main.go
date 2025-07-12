package main

import (
	"database/sql"
	"explorer/app/db"
	"explorer/plugins/auth"
	"log"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err:=godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
	database := db.Get() // or your actual DB connection method

	password := "admin1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	now := time.Now()
	admin := auth.User{
		Email:        "admin@camping.tn",
		PasswordHash: string(hashedPassword),
		FirstName:    "Hamrouni",
		LastName:     "Foued",
		Role:         "admin",
		EmailVerifiedAt: sql.NullTime{
			Time: now,
			Valid: true,
		},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := database.Create(&admin)
	if result.Error != nil {
		log.Fatalf("failed to create admin: %v", result.Error)
	}

	log.Println("✅ Admin user seeded successfully.")
}
