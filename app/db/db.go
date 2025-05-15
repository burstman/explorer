package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// By default this is a pre-configured Gorm DB instance.
// Change this type based on the database package of your likings.
var dbInstance *gorm.DB

// Get returns the instantiated DB instance.
func Get() *gorm.DB {
	return dbInstance
}

func init() {
	// Read PostgreSQL config from environment
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var err error
	dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}
