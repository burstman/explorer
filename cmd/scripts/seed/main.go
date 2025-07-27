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
	database := db.Get()

	// Hash password
	passwordAdmin := "admin1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordAdmin), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	now := time.Now()

	// Create Admin
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
	if err := database.Create(&admin).Error; err != nil {
		log.Fatalf("failed to create admin: %v", err)
	}

	// Create Normal User
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
	if err := database.Create(&user).Error; err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	// Create Bus Types
	buses := []types.BuseType{
		{Name: "MAN", Capacity: 54},
		{Name: "MiniBus", Capacity: 23},
	}
	if err := database.Create(&buses).Error; err != nil {
		log.Fatalf("failed to create bus types: %v", err)
	}

	// Create Services
	services := []types.Service{
		{Name: "tente", Price: 22.5},
		{Name: "Matela", Price: 10.5},
	}
	if err := database.Create(&services).Error; err != nil {
		log.Fatalf("failed to create services: %v", err)
	}

	// Create Campsite (Camp needs to be created before we can use its ID)
	camp := types.CampSite{
		Name:        "best campsite",
		Description: "campsite in Tunisia",
		ImageURL:    "https://i.imgur.com/LmzEIwg.jpeg",
		Location:    "Bizert",
		Price:       80,
	}
	if err := database.Create(&camp).Error; err != nil {
		log.Fatalf("failed to create campsite: %v", err)
	}

	// Now insert buses to campsite_buses using real camp.ID and buses[i].ID
	campBus := []types.CampsiteBus{
		{
			CampsiteID: camp.ID,
			BusTypeID:  buses[0].ID, // MAN
			Quantity:   1,
		},
		{
			CampsiteID: camp.ID,
			BusTypeID:  buses[1].ID, // MiniBus
			Quantity:   1,
		},
	}
	if err := database.Create(&campBus).Error; err != nil {
		log.Fatalf("failed to create campsite buses: %v", err)
	}

	log.Println("✅ Seed successfully.")
}
