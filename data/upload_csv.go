// main.go
package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds application configuration
type Config struct {
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	BatchSize int
}

// Bank represents a bank record
type Bank struct {
	ID          int64     `db:"id"`
	SwiftCode   string    `db:"swift_code"`
	CountryCode string    `db:"country_code"`
	CountryName string    `db:"country_name"`
	BankName    string    `db:"bank_name"`
	Address     string    `db:"address"`
	City        string    `db:"city"`
	TimeZone    string    `db:"time_zone"`
	CreatedAt   time.Time `db:"created_at"`
}

// BankService handles bank-related operations
type BankService struct {
	db     *sqlx.DB
	logger *zap.Logger
	config *Config
}

// NewBankService creates a new BankService
func NewBankService(db *sqlx.DB, logger *zap.Logger, config *Config) *BankService {
	return &BankService{
		db:     db,
		logger: logger,
		config: config,
	}
}

// ProcessCSVFile processes the SWIFT codes CSV file
func (s *BankService) ProcessCSVFile(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Begin transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare insert statement
	stmt, err := tx.PreparexContext(ctx, `
        INSERT INTO banks (
            swift_code, country_code, country_name, bank_name, 
            address, city, time_zone, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (swift_code) DO UPDATE SET
            country_code = EXCLUDED.country_code,
            country_name = EXCLUDED.country_name,
            bank_name = EXCLUDED.bank_name,
            address = EXCLUDED.address,
            city = EXCLUDED.city,
            time_zone = EXCLUDED.time_zone,
            created_at = EXCLUDED.created_at
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var records int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		// Process record
		_, err = stmt.ExecContext(ctx,
			strings.TrimSpace(record[1]),                  // swift_code
			strings.ToUpper(strings.TrimSpace(record[0])), // country_code
			strings.ToUpper(strings.TrimSpace(record[6])), // country_name
			strings.TrimSpace(record[3]),                  // bank_name
			strings.TrimSpace(record[4]),                  // address
			strings.TrimSpace(record[5]),                  // city
			strings.TrimSpace(record[7]),                  // time_zone
			time.Now().UTC(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}

		records++
		if records%1000 == 0 {
			s.logger.Info("Processing records", zap.Int("count", records))
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Completed processing", zap.Int("total_records", records))
	return nil
}

func loadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func initDB(config *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host, config.DB.Port, config.DB.User,
		config.DB.Password, config.DB.Name, config.DB.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize schema
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS banks (
            id SERIAL PRIMARY KEY,
            swift_code VARCHAR(11) UNIQUE NOT NULL,
            country_code VARCHAR(2) NOT NULL,
            country_name VARCHAR(100) NOT NULL,
            bank_name VARCHAR(200) NOT NULL,
            address VARCHAR(500) NOT NULL,
            city VARCHAR(100) NOT NULL,
            time_zone VARCHAR(50) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            CONSTRAINT swift_code_uppercase CHECK (swift_code = UPPER(swift_code)),
            CONSTRAINT country_code_uppercase CHECK (country_code = UPPER(country_code))
        );
        CREATE INDEX IF NOT EXISTS idx_banks_country_code ON banks(country_code);
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize database
	db, err := initDB(config)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Create service
	service := NewBankService(db, logger, config)

	// Process file
	ctx := context.Background()
	if err := service.ProcessCSVFile(ctx, "swift_codes.csv"); err != nil {
		logger.Fatal("Failed to process CSV file", zap.Error(err))
	}
}
