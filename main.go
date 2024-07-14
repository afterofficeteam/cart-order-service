package main

import (
	"cart-order-service/config"
	"cart-order-service/database"
	cartDto "cart-order-service/dto/cart"
	cartHandler "cart-order-service/handler/cart"

	orderDto "cart-order-service/dto/order"
	orderHandler "cart-order-service/handler/order"
	"cart-order-service/repository/cart"
	orderRepository "cart-order-service/repository/order"
	"cart-order-service/routes"
	"database/sql"

	"github.com/go-playground/validator/v10"
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

	validator := validator.New()

	// Set up the renderer and routes
	routes := setupRoutes(sqlDb, validator)

	// Run the server
	routes.Run(cfg.AppPort)
}

// setupRoutes sets up the routes for the application.
func setupRoutes(myDb *sql.DB, validator *validator.Validate) *routes.Routes {
	cartRepository := cart.NewStore(myDb)
	cartSvc := cartDto.NewCart(cartRepository)
	cartHandler := cartHandler.NewHandler(cartSvc)

	orderRepository := orderRepository.NewStore(myDb)
	orderSvc := orderDto.NewOrder(orderRepository)
	orderHandler := orderHandler.NewHandler(orderSvc, validator)

	// Return a new Routes instance
	return &routes.Routes{
		Cart:  cartHandler,
		Order: orderHandler,
	}
}
