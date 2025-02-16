package routerutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func FailDatabase(c *gin.Context, err error, code int) {
	c.JSON(http.StatusInternalServerError, gin.H{"error_msg": "Database error", "details": err.Error(), "error": code})
}

func Ok(c *gin.Context, data interface{}) {
	if viper.GetBool("debug") {
		c.IndentedJSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusOK, data)
	}

}
