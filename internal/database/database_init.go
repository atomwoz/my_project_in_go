package database

import (
	"fmt"
	"log"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GLOBAL Database connection
var DB *gorm.DB
var DEFAULT_TABLE_NAME string

// Helper function to setup the PostgreSQL database connection
func setupDB(dbname string) {
	// Construct DSN from Viper settings
	var timezone string
	var sslmode string

	if viper.GetString("db.timezone") == "" {
		timezone = "Europe/Warsaw"
	} else {
		timezone = viper.GetString("db.timezone")
	}

	if viper.GetString("db.sslmode") == "" {
		sslmode = "disable"
	} else {
		sslmode = viper.GetString("db.sslmode")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		dbname,
		viper.GetInt("db.port"),
		sslmode,
		timezone,
	)

	DEFAULT_TABLE_NAME = viper.GetString("db.table")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// Helper function to setup the SQLite database connection for testing
func setupSQLite(dbPath string) {
	DEFAULT_TABLE_NAME = viper.GetString("db.table")

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}
}

// SetupDatabase initializes the PostgreSQL connection
func SetupDatabase() {
	setupDB(viper.GetString("db.dbname"))
}

// SetupTestDatabase initializes connection to a SQLite database in the test folder
func SetupTestDatabase() {
	config.LoadConfig()
	// Specify the SQLite file path. Adjust if a different path is desired.
	dbPath := "../../tests/test.db"
	setupSQLite(dbPath)
	DB.AutoMigrate(&models.SwiftModel{})
}
