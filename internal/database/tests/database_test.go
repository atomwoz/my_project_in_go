package tests

import (
	"testing"

	"atomwoz.com/remitly_task/internal/database"
	"github.com/magiconair/properties/assert"
)

func TestDBConnection(t *testing.T) {
	database.SetupDatabaseForTesting("../../../config")

	var test int64
	database.DB.Raw("SELECT 2+5").Find(&test)
	assert.Equal(t, test, int64(7))

}
