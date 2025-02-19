package config

import (
	"log"

	"github.com/spf13/viper"
)

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
