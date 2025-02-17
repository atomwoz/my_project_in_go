package controller

import (
	"net/http"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/database/models"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func GetByCountry(c *gin.Context) {
	var rows []models.SwiftBranchRow
	TABLE_NAME := viper.GetString("db.table")

	// Query the database
	err := database.DB.Table(TABLE_NAME).
		Where("country_code = ?", c.Param("country_code")).
		Find(&rows).Error

	if err == gorm.ErrRecordNotFound || len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error_msg": "Country ISO2 code '" + c.Param("country_code") + "' not found", "error": ERRORS.ErrCodeCountryCodeNotFound,
		})
		return
	}

	if err != nil {
		routerutils.FailDatabase(c, err, 2)
		return
	}

	var rowsWithCountryName struct {
		CountryISO2 string `json:"countryISO2"`
		CountryName string `json:"countryName"`
		Branches    []models.SwiftBranchRow
	}
	{
		rowsWithCountryName.CountryISO2 = c.Param("country_code")
		rowsWithCountryName.CountryName = rows[0].CountryName
		rowsWithCountryName.Branches = rows
	}

	routerutils.Ok(c, rowsWithCountryName)
}
