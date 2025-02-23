package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Test environment variables
	t.Run("Load from environment variables", func(t *testing.T) {
		// Setup test environment
		os.Setenv("DB_HOST", "testhost")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		defer func() {
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
		}()

		// Load configuration
		LoadConfig()

		// Verify values
		assert.Equal(t, "testhost", GetHost())
		assert.Equal(t, "testuser", GetUser())
		assert.Equal(t, "testpass", GetPassword())
	})
}

func TestConfigGetters(t *testing.T) {
	// Setup test values
	testCases := []struct {
		name     string
		setup    func()
		test     func(t *testing.T)
		teardown func()
	}{
		{
			name: "Test GetPassword",
			setup: func() {
				viper.Set("DB_PASSWORD", "testpass")
			},
			test: func(t *testing.T) {
				assert.Equal(t, "testpass", GetPassword())
			},
			teardown: func() {
				viper.Set("DB_PASSWORD", "")
			},
		},
		{
			name: "Test GetPort",
			setup: func() {
				viper.Set("DB_PORT", 5433)
			},
			test: func(t *testing.T) {
				assert.Equal(t, 5433, GetPort())
			},
			teardown: func() {
				viper.Set("DB_PORT", 5432)
			},
		},
		{
			name: "Test GetDebugMode",
			setup: func() {
				viper.Set("DEBUG_MODE", true)
			},
			test: func(t *testing.T) {
				assert.True(t, GetDebugMode())
			},
			teardown: func() {
				viper.Set("DEBUG_MODE", false)
			},
		},
		{
			name: "Test SetDebugMode",
			setup: func() {
				SetDebugMode(true)
			},
			test: func(t *testing.T) {
				assert.True(t, GetDebugMode())
			},
			teardown: func() {
				SetDebugMode(false)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			tc.test(t)
			tc.teardown()
		})
	}
}

func TestDefaultValues(t *testing.T) {
	// Clear any existing configuration
	viper.Reset()

	// Load config with defaults
	LoadConfig()

	// Test default values
	assert.Equal(t, 5432, GetPort())
	assert.Equal(t, "banks", GetTable())
	assert.Equal(t, "swiftdb", GetDBName())
	assert.Equal(t, "../../tests/test.db.back", GetTestDB())
	assert.Equal(t, "disable", GetSSLMode())
	assert.True(t, GetDebugMode())
	assert.Equal(t, "Europe/Warsaw", GetTimezone())
}

func TestMissingRequiredValues(t *testing.T) {
	// Clear any existing configuration
	viper.Reset()
	err := os.Rename("../../.env", "../../X")
	if err != nil {
		//t.Errorf("Failed to rename .env file: %v", err)
	}
	defer os.Rename("../../X", "../../.env")

	//Wait to ensure the file is renamed
	time.Sleep(1 * time.Second)

	assert.Panics(t, func() { LoadConfig() })
	// This should panic due to missing required values

}
