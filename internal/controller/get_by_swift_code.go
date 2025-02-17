package controller

import (
	"fmt"
	"net/http"

	"errors"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/database/models"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// SwiftCode retrieves bank details based on the SWIFT code.
func SwiftCode(c *gin.Context) {
	tableName := viper.GetString("db.table")
	code := c.Param("code")
	var bank models.SwiftModel

	//SQL query
	if err := database.DB.Table(tableName).
		Where("swift_code = ?", code).
		Take(&bank).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error_msg": fmt.Sprintf("Swift code '%s' not found", code),
				"error":     ERRORS.ErrCodeSwiftRecordNotFound,
			})
			return
		}
		routerutils.FailDatabase(c, err, ERRORS.ErrCodeInternalDatabase)
		return
	}

	//Branching !!!
	// Query the database for branches
	if bank.IsHeadquarters {
		var branches []models.SwiftBranchRow
		if err := database.DB.Table(tableName).
			Where("bank_symbol = ? AND swift_code <> ?", bank.BankSymbol, bank.SwiftCode).
			Find(&branches).Error; err != nil {
			routerutils.FailDatabase(c, err, ERRORS.ErrCodeInternalDatabase)
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
	routerutils.Ok(c, branch)
}
