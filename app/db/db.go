package db

import (
	"database/sql"
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
	port := "5432"
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Step 1: Connect to the default 'postgres' database to create the target database
	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		host, user, password, port,
	)

	fmt.Println("Connecting to default postgres database...", defaultDSN)
	dbSQL, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		log.Fatalf("failed to connect to default postgres database: %v", err)
	}
	defer dbSQL.Close()

	// Step 2: Create the database if it doesn't exist
	_, err = dbSQL.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil && !isDatabaseAlreadyExistsError(err) {
		log.Fatalf("failed to create database %s: %v", dbname, err)
	}

	// Step 3: Connect to the target database using GORM
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
	dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL database!")
}

// Helper function to check if the error is due to the database already existing
func isDatabaseAlreadyExistsError(err error) bool {
	return err != nil && err.Error() == fmt.Sprintf(`pq: database "%s" already exists`, os.Getenv("DB_NAME"))
}
