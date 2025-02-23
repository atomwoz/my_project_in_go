package routerutils

import (
	"net/http"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/logs"
	"github.com/gin-gonic/gin"
)

// FailDatabase is a helper function to handle database errors.
// Arguments:
// c - gin.Context, got from the controller
// err - error, the database error to handle
// It returns true if the error is not nil.
func FailDatabaseIfError(c *gin.Context, err error) bool {
	if err != nil {
		logs.Error("Unknown database error: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error_msg": err.Error(), "error": 1})
		return true
	}
	return false
}

// Ok is a helper function to return a successful response.
// Arguments:
// c - gin.Context, got from the controller
// data - interface{}, the data to be JSONified and returned
func Ok(c *gin.Context, data interface{}) {
	if config.GetDebugMode() {
		c.IndentedJSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func Created(c *gin.Context) {
	var data = gin.H{"message": "ok"}
	if config.GetDebugMode() {
		c.IndentedJSON(http.StatusCreated, data)
	} else {
		c.JSON(http.StatusCreated, data)
	}
}
