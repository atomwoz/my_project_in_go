package routerutils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"atomwoz.com/remitly_task/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFailDatabaseIfError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("with error", func(t *testing.T) {
		err := errors.New("database error")
		result := FailDatabaseIfError(c, err)

		assert.True(t, result)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})

	t.Run("without error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		result := FailDatabaseIfError(c, nil)

		assert.False(t, result)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestOk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name      string
		debugMode bool
		data      interface{}
	}{
		{
			name:      "debug mode enabled",
			debugMode: true,
			data:      gin.H{"test": "data"},
		},
		{
			name:      "debug mode disabled",
			debugMode: false,
			data:      gin.H{"test": "data"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			config.SetDebugMode(tc.debugMode)

			Ok(c, tc.data)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, strings.ReplaceAll(w.Body.String(), " ", ""), `"test":"data"`)
		})
	}
}

func TestCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name      string
		debugMode bool
	}{
		{
			name:      "debug mode enabled",
			debugMode: true,
		},
		{
			name:      "debug mode disabled",
			debugMode: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			config.SetDebugMode(tc.debugMode)

			Created(c)

			assert.Equal(t, http.StatusCreated, w.Code)
			assert.Contains(t, strings.ReplaceAll(w.Body.String(), " ", ""), `"message":"ok"`)
		})
	}
}
