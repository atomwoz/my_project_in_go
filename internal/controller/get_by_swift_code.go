package controller

import (
	"net/http"

	"errors"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/logs"
	"atomwoz.com/remitly_task/internal/models"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SwiftCode retrieves bank details based on the SWIFT code.
func SwiftCode(c *gin.Context) {
	tableName := config.GetTable()
	code := c.Param("code")

	// Fetch the SWIFT record from the database
	bank, err := database.FetchSwiftRecord(code)

	// Handle wrong SWIFT code
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error_msg": "SWIFT code not found",
			"error":     ERRORS.ErrCodeSwiftRecordNotFound,
		})
		return
	}

	// Handle any other database error
	if routerutils.FailDatabaseIfError(c, err) {
		return
	}

	//Branching !!!
	// Query the database for branches
	if bank.IsHeadquarter {
		var branches []models.SwiftBranchRow
		if err := database.DB.Table(tableName).
			Where("bank_symbol = ? AND swift_code <> ?", bank.BankSymbol, bank.SwiftCode).
			Find(&branches).Error; err != nil {
			routerutils.FailDatabaseIfError(c, err)
			return
		}

		// Apply headquarter response structure
		headquarter := models.SwiftHeadquarterRow{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryCode:   bank.CountryCode,
			CountryName:   bank.CountryName,
			IsHeadquarter: true,
			SwiftCode:     bank.SwiftCode,
			Branches:      branches,
		}
		routerutils.Ok(c, headquarter)
		return
	}

	// Branch response structure
	branch := models.SwiftBranchRow{
		Address:       bank.Address,
		BankName:      bank.BankName,
		CountryCode:   bank.CountryCode,
		CountryName:   bank.CountryName,
		IsHeadquarter: false,
		SwiftCode:     bank.SwiftCode,
	}
	logs.Log("Fetched records for SWIFT code '%s'", code)
	routerutils.Ok(c, branch)
}
