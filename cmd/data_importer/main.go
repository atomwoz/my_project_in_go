package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strings"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// ImportSwiftCodes reads a CSV file and imports Swift codes into the database.
func ImportSwiftCodes(db *gorm.DB, filePath string) {

	countryCodes := make(map[string]string)

	// CSV file handling
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	// Iterate over records and create model from each csv row
	var swiftCodes []models.SwiftModel

	for _, row := range records {
		if len(row) < 8 {
			continue // Skip invalid rows
		}

		//Creating a data model
		countryCode := strings.ToUpper(strings.TrimSpace(row[0]))
		swift := strings.TrimSpace(row[1])
		bank := strings.TrimSpace(row[3])
		address := strings.TrimSpace(row[4])
		city := strings.TrimSpace(row[5])
		countryName := strings.ToUpper(strings.TrimSpace(row[6]))
		timeZone := strings.TrimSpace(row[7])

		symbol := swift[:8]
		isHQ := strings.HasSuffix(swift, "XXX")

		if _, ok := countryCodes[countryCode]; !ok {
			countryCodes[countryCode] = countryName
		} else if countryCodes[countryCode] != countryName {
			log.Fatalf("Country code %s has multiple names: %s and %s", countryCode, countryCodes[countryCode], countryName)
			return
		}

		swiftCodes = append(swiftCodes, models.SwiftModel{
			SwiftCode:     swift,
			BankName:      bank,
			CountryCode:   countryCode,
			CountryName:   countryName,
			City:          city,
			Address:       address,
			TimeZone:      timeZone,
			IsHeadquarter: isHQ,
			BankSymbol:    symbol,
		})
	}

	for i, dbName := range []string{viper.GetString("db.dbname"), viper.GetString("db.testdb")} {

		var exists int

		// Check if the database exists
		checkQuery := "SELECT 1 FROM pg_database WHERE datname = ?"
		if err := db.Raw(checkQuery, dbName).Scan(&exists).Error; err != nil {
			log.Fatalf("Failed to check database existence: %v", err)
		}

		// Create the database if it doesn't exist
		if exists != 1 {
			createQuery := "CREATE DATABASE " + dbName
			if err := db.Exec(createQuery).Error; err != nil {
				log.Fatalf("Failed to create database: %v", err)
			}
			log.Printf("Database %s created successfully!", dbName)
		} else {
			log.Printf("Database %s already exists.", dbName)
		}

		// Handling test database
		if i == 1 {
			database.SetupTestDatabase()
			db = database.DB
		}

		// Migrate the database
		if err := db.AutoMigrate(&models.SwiftModel{}); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}

		// Clearing a table for data
		if err := db.Exec("TRUNCATE TABLE " + viper.GetString("db.table")).Error; err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}

		//Insert data
		if err := db.Create(swiftCodes).Error; err != nil {
			log.Fatalf("Failed to insert Swift codes: %v", err)
		}

	}
	database.SetupDatabase()
	log.Println("Swift codes imported successfully!")
}

func main() {

	filePath := flag.String("file", "swift_codes.csv", "Path to the CSV file")
	flag.Parse()

	config.LoadConfig()
	database.SetupDatabase()

	ImportSwiftCodes(database.DB, *filePath)
}
