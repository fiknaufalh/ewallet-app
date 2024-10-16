package main

import (
	"log"
	"net/http"

	"ewallet-app/internal/config"
	"ewallet-app/internal/controller"
	"ewallet-app/internal/db"
	"ewallet-app/internal/repository"
	"ewallet-app/internal/services"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to the database
	database, err := db.NewDatabase(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize repository
	userRepo := repository.NewUserRepository(database)

	// Initialize service
	walletService := services.NewWalletService(userRepo)

	// Initialize controller
	walletController := controller.NewWalletController(walletService)

	// Set up HTTP routes
	http.HandleFunc("/users", walletController.CreateUser)
	http.HandleFunc("/users/", walletController.GetUser)
	http.HandleFunc("/topup", walletController.TopUp)
	http.HandleFunc("/withdraw", walletController.Withdraw)

	// Start the server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}