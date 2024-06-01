package main

import (
	"Base-Project/config"
	"Base-Project/database"
	"Base-Project/routes"
	"database/sql"
	"github.com/thedevsaddam/renderer"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	sqlDb, err := database.ConnectToDatabase(database.Connection{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		return
	}
	defer sqlDb.Close()

	render := renderer.New()
	routes := setupRoutes(render, sqlDb)
	routes.Run(cfg.AppPort)
}

func setupRoutes(render *renderer.Render, myDb *sql.DB) *routes.Routes {

	return &routes.Routes{}
}
