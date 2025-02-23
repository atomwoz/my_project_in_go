package main

import (
	"fmt"

	"atomwoz.com/remitly_task/internal/config"
	"atomwoz.com/remitly_task/internal/database"
	"atomwoz.com/remitly_task/internal/logs"
	"github.com/charmbracelet/log"
)

func main() {
	config.LoadConfig()
	database.SetupDatabase()
	var records int64
	err := database.DB.Table(database.DEFAULT_TABLE_NAME).Select("*").Count(&records).Error
	if err != nil {
		logs.Fatal("Failed to connect to prod PostgreSQL database: %v", err)
	}
	if records == 0 {
		logs.Fatal("No data in the database")
	}
	log.Infof("Connected to prod PostgreSQL database.")
	log.Infof("There are %d records in the database", records)
	fmt.Println("===== OK! ====== ")
}
