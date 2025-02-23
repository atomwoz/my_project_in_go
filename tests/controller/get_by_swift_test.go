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
var router = rtr.CreateRouter("/v1")

// Initialize database
func init() {
	database.SetupTestDatabase()
}

// Helper function to send a GET request and validate response, ignoring the order of the keys and formatting
func testResponseGet(t *testing.T, swiftCode string, expectedStatus int, expectedBody interface{}) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/"+swiftCode, nil)
	router.ServeHTTP(w, req)

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

// TestGetBySwiftCodeBranch tests a branch swift code.
func TestGetBySwiftCodeBranch(t *testing.T) {
	testResponseGet(t, "ALBPPLP1BMW", http.StatusOK, gin.H{
		"bankName":      "ALIOR BANK SPOLKA AKCYJNA",
		"address":       "WARSZAWA, MAZOWIECKIE",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"isHeadquarter": false,
		"swiftCode":     "ALBPPLP1BMW",
	})
}

// TestGetBySwiftHQ tests a headquarters swift code with branches.
func TestGetBySwiftHQ(t *testing.T) {
	testResponseGet(t, "BKSACLRMXXX", http.StatusOK, gin.H{
		"bankName":      "SCOTIABANK CHILE",
		"address":       "AVENIDA COSTANERA SUR 2710, FLOOR 10 EDIFICIO PARQUE TITANIUM SANTIAGO, PROVINCIA DE SANTIAGO",
		"countryISO2":   "CL",
		"countryName":   "CHILE",
		"isHeadquarter": true,
		"swiftCode":     "BKSACLRMXXX",
		"branches": []map[string]interface{}{
			{"bankName": "SCOTIABANK CHILE", "address": "", "countryISO2": "CL", "countryName": "CHILE", "isHeadquarter": false, "swiftCode": "BKSACLRM055"},
			{"bankName": "SCOTIABANK CHILE", "address": "", "countryISO2": "CL", "countryName": "CHILE", "isHeadquarter": false, "swiftCode": "BKSACLRM061"},
			{"bankName": "SCOTIABANK CHILE", "address": "21 DE MAYO 187  ARICA, PROVINCIA DE ARICA, 1000000", "countryISO2": "CL", "countryName": "CHILE", "isHeadquarter": false, "swiftCode": "BKSACLRM064"},
			{"bankName": "SCOTIABANK CHILE", "address": "", "countryISO2": "CL", "countryName": "CHILE", "isHeadquarter": false, "swiftCode": "BKSACLRM068"},
		},
	})
}

// TestGetNoSwift tests when no swift code is provided.
func TestGetNoSwift(t *testing.T) {
	testResponseGet(t, "", http.StatusNotFound, "404 page not found")
}

// TestWrongSwift tests an incorrect swift code.
func TestWrongSwift(t *testing.T) {
	testResponseGet(t, "ALA_MA_KOTA", http.StatusNotFound, gin.H{
		"error":     float64(2), // JSON unmarshaling interprets numbers as float64
		"error_msg": "SWIFT code not found",
	})
}
