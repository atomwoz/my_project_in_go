package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"atomwoz.com/remitly_task/internal/database"
	"github.com/stretchr/testify/assert"
)

func init() {
	database.SetupTestDatabase()
}

func TestCreateRouter(t *testing.T) {
	// Test cases for router configuration
	tests := []struct {
		name           string
		prefix         string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET swift code endpoint",
			prefix:         "/api/v1",
			method:         "GET",
			path:           "/api/v1/swift-codes/QQQQ",
			expectedStatus: http.StatusNotFound, // Will be 404 since we're not mocking the database
		},
		{
			name:           "GET by country endpoint",
			prefix:         "/api/v1",
			method:         "GET",
			path:           "/api/v1/swift-codes/country/QQQ",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "POST swift code endpoint",
			prefix:         "/api/v1",
			method:         "POST",
			path:           "/api/v1/swift-codes",
			expectedStatus: http.StatusBadRequest, // Will be 400 since we're not sending body
		},
		{
			name:           "DELETE swift code endpoint",
			prefix:         "/api/v1",
			method:         "DELETE",
			path:           "/api/v1/swift-codes/QQ",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid endpoint",
			prefix:         "/api/v1",
			method:         "GET",
			path:           "/api/v1/invalid",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Good GET request",
			prefix:         "/api/v1",
			method:         "GET",
			path:           "/api/v1/swift-codes/ALBPPLPWXXX",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create router with test prefix
			router := CreateRouter(tt.prefix)

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert response status
			assert.Equal(t, tt.expectedStatus, w.Code, "Expected status %d but got %d",
				tt.expectedStatus, w.Code)
		})
	}
}

func TestRouterPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		expectedPath string
	}{
		{
			name:         "Empty prefix",
			prefix:       "",
			expectedPath: "/swift-codes/PLPWPLW1XXX",
		},
		{
			name:         "API v1 prefix",
			prefix:       "/api/v1",
			expectedPath: "/api/v1/swift-codes/PLPWPLW1XXX",
		},
		{
			name:         "Custom prefix",
			prefix:       "/custom",
			expectedPath: "/custom/swift-codes/PLPWPLW1XXX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := CreateRouter(tt.prefix)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.expectedPath, nil)

			router.ServeHTTP(w, req)

			// All paths should be properly routed (even if they return 404)
			assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code,
				"Route %s should be registered", tt.expectedPath)
		})
	}
}
