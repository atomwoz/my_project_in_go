package controller_tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/database"
	rtr "atomwoz.com/remitly_task/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

// Initialize database
func init() {
	database.SetupTestDatabase()
}

// Testing router
var post_router = rtr.CreateRouter("/v1")

// Helper function to test the response of the POST endpoint
func testResponse(t *testing.T, body gin.H, code int, expectedMessage gin.H) {
	w := httptest.NewRecorder()
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	post_router.ServeHTTP(w, req)

	assert.Equal(t, code, w.Code)

	// Comparing the expected message with the actual one in gin.H format, instead of a string
	// 	to avoid problems with the order of the keys, and blank spaces, new lines, etc.
	var actual gin.H
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.Equal(t, nil, err, "Error unmarshaling actual response")
	assert.Equal(t, actual, expectedMessage)
}

// Helper function to delete a record from the database, after successful insertion
func delRecord(swiftCode string) error {
	return database.DB.Exec("DELETE FROM "+config.GetTable()+" WHERE swift_code = ?", swiftCode).Error
}

// TestPostNewCorrectSwiftRow tests the case when good data is inserted
func TestPostNewCorrectSwiftRow(t *testing.T) {

	// Adding headquarter record
	const swiftHQCode = "REMIWWAAXXX"
	testResponse(t, gin.H{
		"swiftCode":     swiftHQCode,
		"countryISO2":   "WW",
		"countryName":   "WIETRZNE WIERCHY",
		"address":       "7th Floor, Penrose Two, Penrose Dock, Cork, Ireland. T23 YY09.",
		"bankName":      "REMITLY",
		"isHeadquarter": true,
	}, http.StatusCreated, gin.H{"message": "ok"})

	// Adding branch record
	const swiftBranchCode = "REMIWWAA111"
	testResponse(t, gin.H{
		"swiftCode":     swiftBranchCode,
		"countryISO2":   "WW",
		"countryName":   "WIETRZNE WIERCHY",
		"address":       "10th Floor, Penrose Three, Penrose Dock, Cork, Iceland. T23 YY11.",
		"bankName":      "REMITLY",
		"isHeadquarter": false,
	}, http.StatusCreated, gin.H{"message": "ok"})

	// Requesting the headquarter record
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/REMIWWAAXXX", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Checking the response, especially the branches in the same way as in helper function testResponse
	// 	by comparing the expected message with the actual one in gin.H format, instead of a string
	// 	to avoid problems with the order of the keys, and blank spaces, new lines, etc.
	var actual gin.H

	expectedResponse := map[string]interface{}{
		"bankName":      "REMITLY",
		"address":       "7th Floor, Penrose Two, Penrose Dock, Cork, Ireland. T23 YY09.",
		"countryISO2":   "WW",
		"countryName":   "WIETRZNE WIERCHY",
		"isHeadquarter": true,
		"swiftCode":     "REMIWWAAXXX",
		"branches": []interface{}{
			map[string]interface{}{
				"bankName":      "REMITLY",
				"address":       "10th Floor, Penrose Three, Penrose Dock, Cork, Iceland. T23 YY11.",
				"countryISO2":   "WW",
				"countryName":   "WIETRZNE WIERCHY",
				"isHeadquarter": false,
				"swiftCode":     "REMIWWAA111",
			},
		},
	}

	// Convert expectedResponse and actual to JSON strings and compare
	expectedJSON, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Error marshaling expected response: %v", err)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Error unmarshaling actual response: %v", err)
	}

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("Error marshaling actual response: %v", err)
	}

	expectedStr := string(expectedJSON)
	actualStr := string(actualJSON)

	assert.Equal(t, expectedStr, actualStr, "Expected and actual responses do not match")

	//If record was inserted, won't return an error
	err = delRecord(swiftHQCode)
	assert.Equal(t, nil, err, "Error deleting headquarter record")

	err = delRecord(swiftBranchCode)
	assert.Equal(t, nil, err, "Error deleting branch record")
}

// TestPostWrongJSON tests the case when the JSON is malformed
func TestPostWrongJSON(t *testing.T) {
	database.SetupTestDatabase()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/swift-codes", strings.NewReader("!@#$%^&*()ść😁<"))
	req.Header.Set("Content-Type", "application/json")
	post_router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "{\"message\":\"malformed JSON payload\"}")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/v1/swift-codes", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	post_router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "{\"message\":\"malformed JSON payload\"}", "Empty JSON should return malformed JSON error")
}

// TestPostInvalidSwiftCode tests the case when the swift code is too short or too long
func TestPostInvalidSwiftCode(t *testing.T) {
	testResponse(t, gin.H{
		"swiftCode":     "XABPLP1BM",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "invalid swift code"})

	testResponse(t, gin.H{
		"swiftCode":     "XABPLP1BMW123",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "invalid swift code"})

}

// TestPostNewCountry tests the case when the country is not in the database
func TestPostNewCountry(t *testing.T) {
	const swiftCode = "XABZAAPBMW1"
	testResponse(t, gin.H{
		"swiftCode":     swiftCode,
		"countryISO2":   "AA",
		"countryName":   "FICTIONAL COUNTRY",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusCreated, gin.H{"message": "ok"})

	//If record was inserted, won't return an error
	err := delRecord(swiftCode)
	assert.Equal(t, err, nil, "Error deleting record")
}

// TestPostInvalidCountryCode tests the case when the country code is too short or too long
func TestPostInvalidCountryCode(t *testing.T) {
	testResponse(t, gin.H{
		"swiftCode":     "XABPLP1BMW1",
		"countryISO2":   "PLA",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "invalid country code"})

	testResponse(t, gin.H{
		"swiftCode":     "XABPLP1BMW1",
		"countryISO2":   "P",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "invalid country code"})
}

// TestPostInvalidCountryName tests the case when the country name is invalid
func TestPostInvalidCountryName(t *testing.T) {
	testResponse(t, gin.H{
		"swiftCode":     "XABXPLP1MW1",
		"countryISO2":   "PL",
		"countryName":   "VATI",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "wrong country name"})
}

// TestPostCountryHasDifferentName tests the case when the country name has a different code
func TestPostCountryHasDifferentCode(t *testing.T) {
	testResponse(t, gin.H{
		"swiftCode":     "XABXZFP1MW1",
		"countryISO2":   "ZF",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "wrong country code"})
}

// TestPostEmptyFields tests the case when the fields are empty
func TestPostEmptyFields(t *testing.T) {
	const swiftCode = "XABAAA1BMW1"
	testResponse(t, gin.H{
		"swiftCode":     swiftCode,
		"countryISO2":   "",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "required fields are empty"})

	testResponse(t, gin.H{
		"swiftCode":     swiftCode,
		"countryISO2":   "PL",
		"countryName":   "",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "required fields are empty"})

	testResponse(t, gin.H{
		"swiftCode":     swiftCode,
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "",
		"bankName":      "ALIOR BANK SPOLKA AKCYJNA",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "required fields are empty"})

	testResponse(t, gin.H{
		"swiftCode":     swiftCode,
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "WARSZAWA, MAZOWIECKIE",
		"bankName":      "",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "required fields are empty"})

}

// TestPostHQFlag tests the case when the isHeadquarter flag is wrong compared to the swift code suffix
func TestPostHearthquatersWrongFlag(t *testing.T) {

	testResponse(t, gin.H{
		"swiftCode":     "AAAAPLAA111",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": true,
	}, http.StatusBadRequest, gin.H{"message": "headquarter flag does not match swift code suffix"})

	testResponse(t, gin.H{
		"swiftCode":     "AAAAPLAAXXX",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "headquarter flag does not match swift code suffix"})

}

// TestPostDuplicateSwiftCode tests the case when the swift code is already in the database
func TestPostDuplicateSwiftCode(t *testing.T) {
	const swiftCode = "AAAAPLAA111"
	testResponse(t, gin.H{
		"swiftCode":     swiftCode,
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusCreated, gin.H{"message": "ok"})

	testResponse(t, gin.H{
		"swiftCode":     swiftCode,
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusConflict, gin.H{"message": "duplicated swift code"})

	// Removed that first record
	err := delRecord(swiftCode)
	assert.Equal(t, err, nil, "Error deleting record")
}

// TestPostWrongRunesInSwiftCode tests the case when the swift code contains wrong characters
func TestPostWrongRunesInSwiftCode(t *testing.T) {
	testResponse(t, gin.H{
		"swiftCode":     "AAAAPLAA1!1",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "invalid swift code"})

	testResponse(t, gin.H{
		"swiftCode":     "AAAAPLAA1 1",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "invalid swift code"})

	testResponse(t, gin.H{
		"swiftCode":     "AAAAPL💗A111",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"address":       "KRAKÓW, MAŁOPOLSKIE",
		"bankName":      "FICTIONAL BANK",
		"isHeadquarter": false,
	}, http.StatusBadRequest, gin.H{"message": "invalid swift code"})
}
