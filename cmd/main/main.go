package main

import (
	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/router"
)

// Entry point
func main() {
	database.SetupDatabase()
	router.CreateRouter("/v1").Run(":8080")
}
