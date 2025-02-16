package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// LoadConfig initializes Viper and loads configuration from config.yaml
func LoadConfig() {
	viper.SetConfigName("config") // Config file name without extension
	viper.SetConfigType("yaml")   // Config type
	viper.AddConfigPath("config") // Path to look for config file

	// Read in environment variables as overrides
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	fmt.Println("Config loaded successfully!")
}

// SetupDatabase initializes the database connection
func SetupDatabase() {
	LoadConfig() // Ensure config is loaded before using it

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
		viper.GetString("db.dbname"),
		viper.GetInt("db.port"),
		viper.GetString("db.sslmode"),
		timezone,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database successfully")
}
