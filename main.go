package main

import (
	"patient-monitor-backend/internal/app"
	"patient-monitor-backend/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.SetupConfig()
	config.SetupDatabase()
	config.SetupLogger()
	app.StartApplication()
}
