package database

import (
	"testing"

	"atomwoz.com/remitly_task/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDatabaseAPI(t *testing.T) {
	// Setup test database
	SetupTestDatabase()

	t.Run("FetchSwiftRecord", func(t *testing.T) {
		// Test fetching existing record
		record, err := FetchSwiftRecord("ALBPPLP1BMW")
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, "ALBPPLP1BMW", record.SwiftCode)

		// Test fetching non-existent record
		record, err = FetchSwiftRecord("NONEXISTENT")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, record)
	})

	t.Run("InsertSwiftRecord", func(t *testing.T) {
		testRecord := &models.SwiftModel{
			SwiftCode:     "TESTPL00XXX",
			CountryCode:   "PL",
			CountryName:   "POLAND",
			Address:       "Test Address",
			BankName:      "Test Bank",
			IsHeadquarter: true,
			BankSymbol:    "TESTPL00",
		}

		// Test inserting new record
		err := InsertSwiftRecord(testRecord)
		assert.NoError(t, err)

		// Verify insertion
		inserted, err := FetchSwiftRecord("TESTPL00XXX")
		assert.NoError(t, err)
		assert.Equal(t, testRecord.SwiftCode, inserted.SwiftCode)
		assert.Equal(t, testRecord.BankName, inserted.BankName)
	})

	t.Run("DeleteSwiftRecord", func(t *testing.T) {
		// First fetch a record we know exists
		record, err := FetchSwiftRecord("TESTPL00XXX")
		assert.NoError(t, err)

		// Test deleting the record
		err = DeleteSwiftRecord(record)
		assert.NoError(t, err)

		// Verify deletion
		_, err = FetchSwiftRecord("TESTPL00XXX")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("GetSiftRecordsByCountryCode", func(t *testing.T) {
		// Test getting records by country code
		records, err := GetSiftRecordsByCountryCode("PL")
		assert.NoError(t, err)
		assert.NotEmpty(t, records)

		// Verify each record has correct country code
		for _, record := range records {
			assert.Equal(t, "PL", record.CountryCode)
		}

		// Test with non-existent country code
		records, err = GetSiftRecordsByCountryCode("XX")
		assert.NoError(t, err)
		assert.Empty(t, records)
	})
}
