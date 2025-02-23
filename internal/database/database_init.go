package database

import (
	"fmt"
	"io"
	"os"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/logs"
	"atomwoz.com/remitly_task/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GLOBAL Database connection
var DB *gorm.DB
var DEFAULT_TABLE_NAME string

// Helper function to setup the PostgreSQL database connection
func setupDB() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.GetHost(),
		config.GetUser(),
		config.GetPassword(),
		config.GetDBName(),
		config.GetPort(),
		config.GetSSLMode(),
		config.GetTimezone(),
	)

	DEFAULT_TABLE_NAME = config.GetTable()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.Fatal("Failed to connect to database: %v", err)
	}
}

// Helper function to setup the SQLite database connection for testing
func setupSQLite(dbPath string) {
	DEFAULT_TABLE_NAME = config.GetTable()

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logs.Fatal("Failed to connect to SQLite database: %v", err)
	}
}

// SetupDatabase initializes the prod PostgreSQL connection
func SetupDatabase() {
	setupDB()
}

// SetupTestDatabase initializes connection to a SQLite database in the test folder
func SetupTestDatabase() {
	config.LoadConfig()
	// Specify the SQLite file path.
	sourceDbPath := config.GetTestDB()
	dbPath := "../../tests/test.db"
	source, err := os.Open(sourceDbPath)
	if err != nil {
		logs.Fatal("Failed to open source SQLite DB: %v", err)
	}
	defer source.Close()

	dest, err := os.Create(dbPath)
	if err != nil {
		logs.Fatal("Failed to create destination SQLite DB: %v", err)
	}
	defer dest.Close()

	if _, err = io.Copy(dest, source); err != nil {
		logs.Fatal("Failed to copy SQLite DB: %v", err)
	}
	setupSQLite(dbPath)
	DB.AutoMigrate(&models.SwiftModel{})
}
