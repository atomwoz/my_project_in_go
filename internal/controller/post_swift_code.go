package controller

import (
	"net/http"
	"regexp"
	"strings"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/database/models"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// PostNewSwiftRow adds a new SWIFT code entry to the database.
func PostNewSwiftRow(c *gin.Context) {
	tableName := viper.GetString("db.table")
	var candidate models.SwiftBranchRow

	// Bind incoming JSON request to struct
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "malformed JSON in request"})
		return
	}
	candidate.SwiftCode = strings.TrimSpace(strings.ToUpper(candidate.SwiftCode))
	candidate.CountryCode = strings.TrimSpace(strings.ToUpper(candidate.CountryCode))
	candidate.CountryName = strings.TrimSpace(strings.ToUpper(candidate.CountryName))
	candidate.Address = strings.TrimSpace(strings.ToUpper(candidate.Address))
	candidate.BankName = strings.TrimSpace(strings.ToUpper(candidate.BankName))

	{
		// Helper function to trim whitespace
		t := func(s string) string { return strings.TrimSpace(s) }
		// Check for empty fields
		if t(candidate.SwiftCode) == "" || t(candidate.CountryCode) == "" || t(candidate.CountryName) == "" || t(candidate.BankName) == "" || t(candidate.Address) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "required fields are empty"})
			return
		}

	}

	// Check for invalid country code (2 uppercase letters)
	if !regexp.MustCompile(`^[A-Z]{2}$`).MatchString(candidate.CountryCode) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid country code"})
		return
	}

	// Check for invalid swift code (8 or 11 characters)
	if len(candidate.SwiftCode) != 8 && len(candidate.SwiftCode) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid swift code"})
		return
	}

	// Check for non matchin county from swift code and country code
	if strings.ToUpper(candidate.SwiftCode[3:5]) != candidate.CountryCode {
		c.JSON(http.StatusBadRequest, gin.H{"message": "country code does not match swift country code"})
		return
	}

	// Check for wrong country name
	{
		var swift models.SwiftModel
		// Query the database for the country names from other swift codes
		err := database.DB.Table(tableName).Where("country_code = ?", candidate.CountryCode).First(&swift).Error
		if err != gorm.ErrRecordNotFound {
			if err != nil {
				routerutils.FailDatabase(c, err, ERRORS.ErrCodeInternalDatabase)
				return
			}
			if swift.CountryName != candidate.CountryName {
				c.JSON(http.StatusBadRequest, gin.H{"message": "wrong country name"})
				return
			}
		}

	}

	//Check for existing country name
	{
		var swift models.SwiftModel
		// Query the database for the country names from other swift codes
		err := database.DB.Table(tableName).Where("country_name = ?", candidate.CountryName).First(&swift).Error
		if err != gorm.ErrRecordNotFound {
			if err != nil {
				routerutils.FailDatabase(c, err, ERRORS.ErrCodeInternalDatabase)
				return
			}
			if swift.CountryCode != candidate.CountryCode {
				c.JSON(http.StatusBadRequest, gin.H{"message": "wrong country code"})
				return
			}
		}

	}

	//Check for headquarter flag and swift code
	if (candidate.IsHeadquarter && !strings.HasSuffix(candidate.SwiftCode, "XXX")) || (!candidate.IsHeadquarter && strings.HasSuffix(candidate.SwiftCode, "XXX")) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "headquarter flag does not match swift code suffix"})
		return
	}

	//SQL INSERT query
	err := database.DB.Table(tableName).Create(&candidate).Error

	if err != nil && strings.Contains(err.Error(), "(SQLSTATE 23505)") {
		c.JSON(http.StatusConflict, gin.H{"message": "duplicated swift code entry"})
		return
	}

	// Check for any other database error
	if err != nil {
		routerutils.FailDatabase(c, err, ERRORS.ErrCodeInternalDatabase)
		return
	}

	// Success response
	routerutils.Ok(c, gin.H{"message": "ok"})
}
