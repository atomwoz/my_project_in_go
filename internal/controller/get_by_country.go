package controller

import (
	"errors"
	"net/http"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/models"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	tableName := viper.GetString("db.table")
	var rows []models.SwiftBranchRow

	// SQL query
	err := database.DB.Table(tableName).
		Where("country_code = ?", countryCode).
		Find(&rows).Error

	// Handle wrong SWIFT code
	if errors.Is(err, gorm.ErrRecordNotFound) || len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error_msg": "Country ISO2 code not found",
			"error":     ERRORS.ErrCodeCountryCodeNotFound,
		})
		return
	}

	// Handle any other database error
	if err != nil {
		routerutils.FailDatabase(c, err, ERRORS.ErrCodeInternalDatabase)
		return
	}

	// Creating final response
	response := CountryBranchesResponse{
		CountryISO2: countryCode,
		CountryName: rows[0].CountryName,
		Branches:    rows,
	}

	routerutils.Ok(c, response)
}
