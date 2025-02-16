package controller

import (
	"net/http"

	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/database/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SwiftCode retrieves bank details based on the SWIFT code
func SwiftCode(c *gin.Context) {
	var bank models.SwiftRow

	// Query the database
	err := database.DB.Table("banks").
		Select("*").
		Where("swift_code = ?", c.Param("code")).
		Take(&bank).Error

	// Handle database errors
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Swift code not found",
			})
			return
		}
		// Handle other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	// Return the JSON response
	c.IndentedJSON(http.StatusOK, bank)
}
