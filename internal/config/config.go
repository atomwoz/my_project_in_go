package config

import (

	//"strings"

	"strings"

	"github.com/spf13/viper"
)

// LoadConfig initializes Viper and loads configuration from config.yaml
func LoadConfig() {
	//viper.SetConfigName("config")      // Config file name without extension
	//iper.SetConfigType("yaml")        // Config type
	//viper.AddConfigPath(conf_location) // Path to look for config file

	// Read in environment variables as overrides
	// viper.SetConfigType("env")
	// viper.SetConfigName(".env")
	// viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Replace dots with underscores to match environment variable style.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//Bind the environment variable DB_PASSWORD to the key "db.password".
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.table", "DB_TABLE")
	viper.BindEnv("db.dbname", "DB_DBNAME")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.testdb", "DB_TESTDB")
	viper.BindEnv("db.sslmode", "DB_SSLMODE")
	viper.BindEnv("batch_size", "DB_BATCH_SIZE")
	viper.BindEnv("debug", "DEBUG_MODE")
	viper.BindEnv("db.timezone", "DB_TIMEZONE")

	// Read config file
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("Error reading config file: %v", err)
	// }

	//log.Println("Config loaded successfully!")
}
