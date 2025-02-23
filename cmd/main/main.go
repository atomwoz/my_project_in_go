package main

import (
	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/router"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

const PORT = 8080

// Entry point
func main() {
	config.LoadConfig()
	database.SetupDatabase()
	if config.GetDebugMode() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Infof("Starting server on port %d", PORT)
	err := router.CreateRouter("/v1").Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
