package controller

import (
	"net/http"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/database/models"
	"atomwoz.com/remitly_task/internal/router/routerutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// SwiftCode retrieves bank details based on the SWIFT code
func SwiftCode(c *gin.Context) {
	var bank models.SwiftModel
	TABLE_NAME := viper.GetString("db.table")

	// Query the database
	err := database.DB.Table(TABLE_NAME).
		Where("swift_code = ?", c.Param("code")).
		Take(&bank).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error_msg": "Swift code not found", "error": 1,
			})
			return
		}
		routerutils.FailDatabase(c, err, 2)
		return
	}

	//Branching !!!
	// Query the database for branches
	if bank.IsHeadquarters {
		var branches []models.SwiftBranchRow
		err := database.DB.Table(TABLE_NAME).Where("bank_symbol = ? AND swift_code <> ?", bank.BankSymbol, bank.SwiftCode).Find(&branches).Error
		if err != nil {
			routerutils.FailDatabase(c, err, 2)
		}
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
