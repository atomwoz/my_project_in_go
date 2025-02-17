package routerutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// FailDatabase is a helper function to handle database errors.
func FailDatabase(c *gin.Context, err error, code int) {
	c.JSON(http.StatusInternalServerError, gin.H{"error_msg": err.Error(), "error": code})
}

// Ok is a helper function to return a successful response.
func Ok(c *gin.Context, data interface{}) {
	if viper.GetBool("debug") {
		c.IndentedJSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusOK, data)
	}

}
