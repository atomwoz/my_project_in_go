package main

import (
	"atomwoz.com/remitly_task/src/database"
	"atomwoz.com/remitly_task/src/router"
)

func main() {
	database.SetupDatabase()
	router.CreateRouter("/v1").Run(":8080")
}
