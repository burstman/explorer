package main

import (
	"camping/app/db"
	"camping/plugins/auth"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	database := db.Get() // or your actual DB connection method

	password := "admin1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	admin := auth.User{
		Email:        "admin@camping.tn",
		PasswordHash: string(hashedPassword),
		FirstName:    "Hamrouni",
		LastName:     "Foued",
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := database.Create(&admin)
	if result.Error != nil {
		log.Fatalf("failed to create admin: %v", result.Error)
	}

	log.Println("✅ Admin user seeded successfully.")
}
