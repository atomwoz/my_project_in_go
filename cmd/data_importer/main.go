package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strings"

	"atomwoz.com/remitly_task/internal/database"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// SwiftCode represents a bank's headquarters or branch.
type SwiftCode struct {
	SwiftCode   string `gorm:"primaryKey"`
	BankName    string
	CountryCode string
	CountryName string
	City        string
	Address     string
	TimeZone    string
	Headquarter string // Stores HQ SwiftCode if it's a branch
}

// ImportSwiftCodes reads a CSV file and imports Swift codes into the database.
func ImportSwiftCodes(db *gorm.DB, filePath string) {
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

	swiftMap := make(map[string]string) // Stores HQ Swift codes

	// Iterate over records and store valid Swift codes
	var swiftCodes []SwiftCode
	for _, row := range records {
		if len(row) < 8 {
			continue // Skip invalid rows
		}

		countryCode := strings.ToUpper(strings.TrimSpace(row[0]))
		swift := strings.TrimSpace(row[1])
		bank := strings.TrimSpace(row[3])
		address := strings.TrimSpace(row[4])
		city := strings.TrimSpace(row[5])
		countryName := strings.ToUpper(strings.TrimSpace(row[6]))
		timeZone := strings.TrimSpace(row[7])

		// Identify headquarters and branches
		hqSwift := swift
		if !strings.HasSuffix(swift, "XXX") {
			hqSwift = swift[:8] + "XXX" // Derive HQ Swift from first 8 characters
		}
		swiftMap[swift[:8]] = hqSwift

		swiftCodes = append(swiftCodes, SwiftCode{
			SwiftCode:   swift,
			BankName:    bank,
			CountryCode: countryCode,
			CountryName: countryName,
			City:        city,
			Address:     address,
			TimeZone:    timeZone,
			Headquarter: hqSwift, // Store HQ reference
		})
	}

	// Check if the database exists, if not create it
	dbName := viper.GetString("db.dbname")
	var exists int
	checkQuery := "SELECT 1 FROM pg_database WHERE datname = ?"
	if err := db.Raw(checkQuery, dbName).Scan(&exists).Error; err != nil {
		log.Fatalf("Failed to check database existence: %v", err)
	}

	if exists != 1 {
		// Build the CREATE DATABASE query string (placeholders can't be used here)
		createQuery := "CREATE DATABASE " + dbName
		if err := db.Exec(createQuery).Error; err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		log.Printf("Database %s created successfully!", dbName)
	} else {
		log.Printf("Database %s already exists.", dbName)
	}

	// Migrate the database
	if err := db.AutoMigrate(&SwiftCode{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create a table for Swift codes
	if err := db.Exec("TRUNCATE TABLE swift_codes").Error; err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}

	// Push data to the database
	if err := db.Create(&swiftCodes).Error; err != nil {
		log.Fatalf("Failed to insert Swift codes: %v", err)
	}

	log.Println("Swift codes imported successfully!")
}

func main() {

	filePath := flag.String("file", "data/swift_codes.csv", "Path to the CSV file")
	flag.Parse()

	database.SetupDatabase()

	ImportSwiftCodes(database.DB, *filePath)
}
