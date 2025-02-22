package main

import (
	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Entry point
func main() {
	config.LoadConfig()
	database.SetupDatabase()
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router.CreateRouter("/v1").Run(":8080")
}
