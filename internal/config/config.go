package config

import (

	//"strings"
	"log"

	"github.com/spf13/viper"
)

// LoadConfig initializes Viper and loads configuration from .env and environment variables.
func LoadConfig() {

	// Set default values
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_HOST", "")
	viper.SetDefault("DB_TABLE", "banks")
	viper.SetDefault("DB_DBNAME", "swiftdb")
	viper.SetDefault("DB_USER", "")
	viper.SetDefault("DB_TESTDB", "../tests/testdb")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DEBUG_MODE", true)
	viper.SetDefault("DB_TIMEZONE", "Europe/Warsaw")

	// Read in configuration from .env file
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath("../../")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	c := func(key string) bool {
		return viper.GetString(key) != ""
	}
	if !c("DB_HOST") || !c("DB_USER") || !c("DB_PASSWORD") {
		log.Fatalf("Missing required configuration values. \n Tried to get from .env file %s and enviroment \n Required are: DB_HOST, DB_USER, DB_PASSWORD", viper.ConfigFileUsed())
	}

}

// ============== Here is config API wrapper for Viper ==============

// GetPassword returns the database password.
func GetPassword() string {
	return viper.GetString("DB_PASSWORD")
}

// GetPort returns the database port.
func GetPort() int {
	return viper.GetInt("DB_PORT")
}

// GetHost returns the database host.
func GetHost() string {
	return viper.GetString("DB_HOST")
}

// GetTable returns the database table name.
func GetTable() string {
	return viper.GetString("DB_TABLE")
}

// GetDBName returns the database name.
func GetDBName() string {
	return viper.GetString("DB_DBNAME")
}

// GetUser returns the database user.
func GetUser() string {
	return viper.GetString("DB_USER")
}

// GetTestDB returns the test database path.
func GetTestDB() string {
	return viper.GetString("DB_TESTDB")
}

// GetSSLMode returns the database SSL mode.
func GetSSLMode() string {
	return viper.GetString("DB_SSLMODE")
}

// GetDebugMode returns the debug mode.
func GetDebugMode() bool {
	return viper.GetBool("DEBUG_MODE")
}

// GetTimezone returns the database timezone.
func GetTimezone() string {
	return viper.GetString("DB_TIMEZONE")
}

// SetDebugMode sets the debug mode.
func SetDebugMode(debug bool) {
	viper.Set("DEBUG_MODE", debug)
}
