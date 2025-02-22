package database_tests

import (
	"log"
	"strings"
	"testing"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/models"
	"github.com/magiconair/properties/assert"
	"gorm.io/gorm"
)

func TestDBConnection(t *testing.T) {
	database.SetupTestDatabase()

	var test int64
	database.DB.Raw("SELECT 2+5").Find(&test)
	assert.Equal(t, test, int64(7))

}
func TestDBFetch(t *testing.T) {
	record, err := database.FetchSwiftRecord("ALBPPLP1BMW")
	log.Println(record)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, record == (*models.SwiftModel)(nil), false, "Record should not be nil")
	//assert.Equal(t, record.BankName, "AMAGIS CAPITAL FUNDS SICAV PLC")

	record, err = database.FetchSwiftRecord("test")
	assert.Equal(t, err, gorm.ErrRecordNotFound, "Error should be gorm.ErrRecordNotFound")
	assert.Equal(t, record, (*models.SwiftModel)(nil), "Record should be nil")

}

func TestDBInsertAndDelete(t *testing.T) {

	x, _ := database.FetchSwiftRecord("ACFCWWM1XXX")
	database.DeleteSwiftRecord(x)

	record := &models.SwiftModel{
		SwiftCode:     "ACFCWWM1XXX",
		CountryCode:   "WW",
		CountryName:   "WIETRZNE WIERCHY",
		Address:       "10th Floor, Penrose Three, Penrose Dock, Cork, Iceland. T23 YY11.",
		BankName:      "REMITLY",
		IsHeadquarter: false,
	}

	err := database.InsertSwiftRecord(record)
	assert.Equal(t, err, nil)

	record, err = database.FetchSwiftRecord("ACFCWWM1XXX")
	assert.Equal(t, err, nil)
	assert.Equal(t, record.BankName, "REMITLY")

	err = database.InsertSwiftRecord(record)
	assert.Equal(t, strings.Contains(err.Error(), "(SQLSTATE 23505)"), true, "Error should contain SQLSTATE 23505")

	x, _ = database.FetchSwiftRecord("ACFCWWM1XXX")
	err = database.DeleteSwiftRecord(x)
	assert.Equal(t, err, nil)

}
