package controller

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/logs"
	"atomwoz.com/remitly_task/internal/models"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostNewSwiftRow adds a new SWIFT code entry to the database.
func PostNewSwiftRow(c *gin.Context) {
	tableName := config.GetTable()
	var candidate models.SwiftModel

	if err := c.ShouldBindJSON(&candidate); err != nil {
		logs.Warn("Malformed JSON payload: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "malformed JSON payload"})
		return
	}

	normalizeSwiftData(&candidate)

	if err := validateSwiftData(&candidate, tableName); err != nil {
		logs.Warn("Invalid SWIFT data: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := database.InsertSwiftRecord(&candidate); err != nil {
		if strings.Contains(err.Error(), "(SQLSTATE 23505)") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			logs.Warn("Duplicated SWIFT code: %s", candidate.SwiftCode)
			c.JSON(http.StatusConflict, gin.H{"message": "duplicated swift code"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	logs.Log("Added new SWIFT code: %s", candidate.SwiftCode)
	routerutils.Created(c)
}

func normalizeSwiftData(candidate *models.SwiftModel) {
	candidate.SwiftCode = strings.TrimSpace(strings.ToUpper(candidate.SwiftCode))
	candidate.CountryCode = strings.TrimSpace(strings.ToUpper(candidate.CountryCode))
	candidate.CountryName = strings.TrimSpace(strings.ToUpper(candidate.CountryName))
	candidate.Address = strings.TrimSpace(candidate.Address)
	candidate.BankName = strings.TrimSpace(candidate.BankName)
	candidate.BankSymbol = candidate.SwiftCode[:8]
}

func validateSwiftData(candidate *models.SwiftModel, tableName string) error {
	if candidate.SwiftCode == "" || candidate.CountryCode == "" || candidate.CountryName == "" || candidate.BankName == "" || candidate.Address == "" {
		return fmt.Errorf("required fields are empty")
	}

	if !regexp.MustCompile(`^[A-Z]{2}$`).MatchString(candidate.CountryCode) {
		return fmt.Errorf("invalid country code")
	}

	if !regexp.MustCompile(`^[A-Z ]+$`).MatchString(candidate.CountryName) {
		return fmt.Errorf("invalid country name")
	}

	if (len(candidate.SwiftCode) != 8 && len(candidate.SwiftCode) != 11) || !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(candidate.SwiftCode) {
		return fmt.Errorf("invalid swift code")
	}

	if strings.ToUpper(candidate.SwiftCode[4:6]) != candidate.CountryCode {
		return fmt.Errorf("country code does not match swift country code")
	}

	var swift models.SwiftModel
	if err := database.DB.Table(tableName).Where("country_code = ?", candidate.CountryCode).First(&swift).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			return err
		}
		if swift.CountryName != candidate.CountryName {
			return fmt.Errorf("wrong country name")
		}
	}

	if err := database.DB.Table(tableName).Where("country_name = ?", candidate.CountryName).First(&swift).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			return err
		}
		if swift.CountryCode != candidate.CountryCode {
			return fmt.Errorf("wrong country code")
		}
	}

	if (candidate.IsHeadquarter && !strings.HasSuffix(candidate.SwiftCode, "XXX")) || (!candidate.IsHeadquarter && strings.HasSuffix(candidate.SwiftCode, "XXX")) {
		return fmt.Errorf("headquarter flag does not match swift code suffix")
	}

	return nil
}
