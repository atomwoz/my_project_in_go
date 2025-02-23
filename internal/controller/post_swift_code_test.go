package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	database.SetupTestDatabase()
}

// Testing normalizeSwiftData function from post_swift_code.go.
// It tests if the function correctly normalizes SWIFT data.
func TestNormalizeSwiftData(t *testing.T) {
	tests := []struct {
		name     string
		input    models.SwiftModel
		expected models.SwiftModel
	}{
		{
			name: "Normalize spaces and case",
			input: models.SwiftModel{
				SwiftCode:   " plpwplw1xxx ",
				CountryCode: " pl ",
				CountryName: "poland ",
				Address:     " test address ",
				BankName:    " test bank ",
			},
			expected: models.SwiftModel{
				SwiftCode:   "PLPWPLW1XXX",
				CountryCode: "PL",
				CountryName: "POLAND",
				Address:     "test address",
				BankName:    "test bank",
				BankSymbol:  "PLPWPLW1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalizeSwiftData(&tt.input)
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

// Testing validateSwiftData function from post_swift_code.go.
// It tests edge cases for validating SWIFT data. Like country code, SWIFT code length and country code mismatch.
func TestValidateSwiftData(t *testing.T) {
	tests := []struct {
		name      string
		input     models.SwiftModel
		wantError bool
		errorMsg  string
	}{
		{
			name: "Valid headquarters SWIFT",
			input: models.SwiftModel{
				SwiftCode:     "PLPWPLW1XXX",
				CountryCode:   "PL",
				CountryName:   "POLAND",
				Address:       "Test Address",
				BankName:      "Test Bank",
				IsHeadquarter: true,
			},
			wantError: false,
		},
		{
			name: "Invalid country code",
			input: models.SwiftModel{
				SwiftCode:   "PLPWPLW1XXX",
				CountryCode: "POL",
				CountryName: "POLAND",
				Address:     "Test Address",
				BankName:    "Test Bank",
			},
			wantError: true,
			errorMsg:  "invalid country code",
		},
		{
			name: "Invalid SWIFT code length",
			input: models.SwiftModel{
				SwiftCode:   "PLPW",
				CountryCode: "PL",
				CountryName: "POLAND",
				Address:     "Test Address",
				BankName:    "Test Bank",
			},
			wantError: true,
			errorMsg:  "invalid swift code",
		},
		{
			name: "Mismatched country code in SWIFT",
			input: models.SwiftModel{
				SwiftCode:   "PLPWDEW1XXX",
				CountryCode: "PL",
				CountryName: "POLAND",
				Address:     "Test Address",
				BankName:    "Test Bank",
			},
			wantError: true,
			errorMsg:  "country code does not match swift country code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSwiftData(&tt.input)
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Testing PostNewSwiftRow function from post_swift_code.go.
// It tests if the function correctly adds new SWIFT code to the database.
func TestPostNewSwiftRow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		input      models.SwiftModel
		statusCode int
	}{
		{
			name: "Valid SWIFT code",
			input: models.SwiftModel{
				SwiftCode:     "PLPWPLW1XXX",
				CountryCode:   "PL",
				CountryName:   "POLAND",
				Address:       "Test Address",
				BankName:      "Test Bank",
				IsHeadquarter: true,
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Invalid SWIFT code",
			input: models.SwiftModel{
				SwiftCode:   "INVALID",
				CountryCode: "PL",
				CountryName: "POLAND",
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest("POST", "/v1/swift-code/", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			PostNewSwiftRow(c)
			defer database.DeleteSwiftRecord(&tt.input)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
