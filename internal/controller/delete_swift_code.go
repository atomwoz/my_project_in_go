package controller

import (
	"errors"
	"net/http"

	"atomwoz.com/remitly_task/internal/database"
	routerutils "atomwoz.com/remitly_task/internal/router/router_utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteSwiftRow(c *gin.Context) {
	code := c.Param("code")

	// Select record to delete
	recordToDelete, err := database.FetchSwiftRecord(code)

	// Handle wrong SWIFT code
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "SWIFT code not found",
		})
		return
	}

	// Handle any other database error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Deleting the record
	err = database.DeleteSwiftRecord(recordToDelete)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	routerutils.Ok(c, gin.H{"message": "Deleted successfully"})
}
