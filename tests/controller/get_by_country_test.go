package controller_tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"atomwoz.com/remitly_task/internal/database"
	rtr "atomwoz.com/remitly_task/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

// Testing router
var routerCountry = rtr.CreateRouter("/v1")

// Initialize database
func init() {
	database.SetupTestDatabase()
}

// Helper function to send a GET request and validate response, ignoring the order of the keys and formatting
func testResponseGetL(t *testing.T, path string, expectedStatus int, expectedBody gin.H) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	routerCountry.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)

	body := w.Body.Bytes()

	// Comparing the expected message with the actual one in gin.H format, instead of a string
	// 	to avoid problems with the order of the keys, and blank spaces, new lines, etc.
	var actual map[string]interface{}
	err := json.Unmarshal(body, &actual)

	if err != nil {
		// If it's not JSON, treat it as a string and compare raw response
		assert.Equal(t, string(body), expectedBody, "Unexpected plaintext response")
	} else {
		// Normalize expected JSON
		expectedJSON, _ := json.Marshal(expectedBody)
		var expected map[string]interface{}
		_ = json.Unmarshal(expectedJSON, &expected)

		assert.Equal(t, expected, actual)
	}
}

// TestGetByCountryCode tests the case when the country code is correct.
func TestGetByCountryCode(t *testing.T) {
	testResponseGetL(t, "/v1/swift-codes/country/AW", http.StatusOK, gin.H{
		"countryISO2": "AW",
		"countryName": "ARUBA",
		"branches": []gin.H{
			{"bankName": "AIB BANK NV", "address": "WILHELMINASTRAAT 36  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "ANIBAWA1XXX"},
			{"bankName": "ARUBA BANK, LTD", "address": "CAMACURI 12  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "ARUBAWAXXXX"},
			{"bankName": "CENTRALE BANK VAN ARUBA", "address": "J.E. IRAUSQUIN 8  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "CBARAWAWXXX"},
			{"bankName": "CARIBBEAN MERCANTILE BANK N.V.", "address": "KAYA GILBERTO FRANCOIS CROES 53 ORANJESTAD, ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "CMBAAWAXXXX"},
			{"bankName": "INTERBANK ARUBA NV", "address": "CAYA G.F. CROES 38  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "IARUAWA1XXX"},
			{"bankName": "IMTRADEX INTERNATIONAL N.V.", "address": "TANKI LENDEERT 143  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "IMIEAWA1XXX"},
			{"bankName": "RBC ROYAL BANK (ARUBA) N.V. (FORMERLY RBTT BANK ARUBA N.V.)", "address": "ITALIESTRAAT 36  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "RBTTAWAWXXX"},
			{"bankName": "BANCO DI CARIBE (ARUBA) N.V", "address": "VONDELLAAN 31  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST", "countryISO2": "AW", "countryName": "ARUBA", "isHeadquarter": true, "swiftCode": "BDCCAWAWXXX"},
		},
	})
}

// TestGetNoCountry tests the case when the country code is not found.
func TestGetNoCountry(t *testing.T) {
	testResponseGetL(t, "/v1/swift-codes/country", http.StatusNotFound, gin.H{"error": 2, "error_msg": "Swift code not found"})
}

// TestWrongCountry tests the case when the country code is wrong.
func TestWrongCountry(t *testing.T) {
	testResponseGetL(t, "/v1/swift-codes/country/XX", http.StatusNotFound, gin.H{"error": 3, "error_msg": "Country ISO2 code not found"})
}
