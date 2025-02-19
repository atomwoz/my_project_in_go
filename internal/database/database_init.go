package database

import (
	"fmt"
	"log"

	"atomwoz.com/remitly_task/internal/config"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GLOBAL Database connection
var DB *gorm.DB
var DEFAULT_TABLE_NAME string

// Helper function to setup the database connection
func setupDB(dbname string) {

	// Construct DSN from Viper settings
	var timezone string

	if viper.GetString("db.timezone") == "" {
		timezone = "Europe/Warsaw"
	} else {
		timezone = viper.GetString("db.timezone")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		dbname,
		viper.GetInt("db.port"),
		viper.GetString("db.sslmode"),
		timezone,
	)

	DEFAULT_TABLE_NAME = viper.GetString("db.table")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//log.Println("Connected to database successfully")
}

// SetupDatabase initializes the database connection
func SetupDatabase() {
	setupDB(viper.GetString("db.dbname"))
}

// SetupTestDatabase initializes connection to the test database
func SetupTestDatabase() {
	config.LoadConfig("../../config")
	setupDB(viper.GetString("db.testdb"))
}

// SetupTestDatabaseWithConfig initializes connection to the test database with a custom configuration location
func SetupTestDatabaseWithConfig(location string) {
	config.LoadConfig(location)
	setupDB(viper.GetString("db.testdb"))
}

// SetupDatabaseForTesting initializes connection to the test database
func SetupDatabaseForTesting() {
	config.LoadConfig("../../config")
	setupDB(viper.GetString("db.dbname"))
}
