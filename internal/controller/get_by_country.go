package controller

import (
	"errors"
	"net/http"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/logs"
	"atomwoz.com/remitly_task/internal/models"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountryBranchesResponse is a response structure for the GetByCountry controller.
type CountryBranchesResponse struct {
	CountryISO2 string                  `json:"countryISO2"`
	CountryName string                  `json:"countryName"`
	Branches    []models.SwiftBranchRow `json:"branches"`
}

// GetByCountry retrieves all branches in a country based on the country code.
func GetByCountry(c *gin.Context) {
	countryCode := c.Param("country_code")

	rows, err := database.GetSiftRecordsByCountryCode(countryCode)

	// Handle wrong SWIFT code
	if errors.Is(err, gorm.ErrRecordNotFound) || len(rows) == 0 {
		logs.Warn("Country ISO2 code '%s' not found", countryCode)
		c.JSON(http.StatusNotFound, gin.H{
			"error_msg": "Country ISO2 code not found",
			"error":     ERRORS.ErrCodeCountryCodeNotFound,
		})
		return
	}

	if routerutils.FailDatabaseIfError(c, err) {
		return
	}

	// Creating final response
	response := CountryBranchesResponse{
		CountryISO2: countryCode,
		CountryName: rows[0].CountryName,
		Branches:    rows,
	}

	logs.Log("Fetched records for country '%s'", countryCode)
	routerutils.Ok(c, response)
}
