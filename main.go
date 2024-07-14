package main

import (
	"cart-order-service/config"
	"cart-order-service/database"
	cartDto "cart-order-service/dto/cart"
	cartHandler "cart-order-service/handler/cart"
	"cart-order-service/repository/cart"
	"cart-order-service/routes"
	"database/sql"
)

// main is the main function that starts the application.
// It loads the configuration, connects to the database, sets up the routes, and starts the server.
func main() {
	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	// Connect to the database
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

	// Set up the renderer and routes
	routes := setupRoutes(sqlDb)

	// Run the server
	routes.Run(cfg.AppPort)
}

// setupRoutes sets up the routes for the application.
func setupRoutes(myDb *sql.DB) *routes.Routes {
	cartRepository := cart.NewStore(myDb)
	cartSvc := cartDto.NewCart(cartRepository)
	cartHandler := cartHandler.NewHandler(cartSvc)

	// Return a new Routes instance
	return &routes.Routes{
		Cart: cartHandler,
	}
}
