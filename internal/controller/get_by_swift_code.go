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
				"error_msg": "Swift code '" + c.Param("code") + "' not found", "error": ERRORS.ErrCodeSwiftRecordNotFound,
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
		err := database.DB.Table(TABLE_NAME).Where("bank_symbol = ? AND swift_code <> ?", bank.BankSymbol, bank.SwiftCode).Find(&branches).Error
		if err != nil {
			routerutils.FailDatabase(c, err, ERRORS.ErrCodeInternalDatabase)
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
