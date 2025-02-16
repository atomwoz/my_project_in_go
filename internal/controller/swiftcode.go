package controller

import (
	"net/http"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/database/model"
	"github.com/gin-gonic/gin"
)

func SwiftCode(c *gin.Context) {
	var row model.SwiftRow

	err := database.DB.Table("banks").Model(model.SwiftRow).
		Select("swift_code", "country_code", "country_name", "bank_name", "address", "city", "time_zone").
		Where("swift_code = ?", c.Param("code")).
		First(model.SwiftRow).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swift code not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, row)
}
