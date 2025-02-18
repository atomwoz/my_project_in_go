package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GLOBAL Database connection
var DB *gorm.DB

// LoadConfig initializes Viper and loads configuration from config.yaml
func LoadConfig(conf_location string) {
	viper.SetConfigName("config")      // Config file name without extension
	viper.SetConfigType("yaml")        // Config type
	viper.AddConfigPath(conf_location) // Path to look for config file

	// Read in environment variables as overrides
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	//log.Println("Config loaded successfully!")
}

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

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//log.Println("Connected to database successfully")
}

// SetupDatabase initializes the database connection
func SetupDatabase() {
	LoadConfig("config")
	setupDB(viper.GetString("db.dbname"))
}

// SetupTestDatabase initializes connection to the test database
func SetupTestDatabase() {
	LoadConfig("../../config")
	setupDB(viper.GetString("db.testdb"))
}

// SetupTestDatabaseWithConfig initializes connection to the test database with a custom configuration location
func SetupTestDatabaseWithConfig(location string) {
	LoadConfig(location)
	setupDB(viper.GetString("db.testdb"))
}

// SetupDatabaseForTesting initializes connection to the test database
func SetupDatabaseForTesting() {
	LoadConfig("../../config")
	setupDB(viper.GetString("db.dbname"))
}
