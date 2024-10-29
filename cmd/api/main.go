package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ewallet-app/internal/config"
	"ewallet-app/internal/domain/repository"
	"ewallet-app/internal/domain/usecase"
	"ewallet-app/internal/handler"
	"ewallet-app/internal/middleware"
	"ewallet-app/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize router with logging middleware
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	idempotencyRepo := repository.NewIdempotencyRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(db, userRepo, walletRepo)
	topUpUseCase := usecase.NewTopUpUseCase(db, walletRepo, transactionRepo, idempotencyRepo, cfg)
	withdrawalUseCase := usecase.NewWithdrawalUseCase(db, walletRepo, transactionRepo, idempotencyRepo, cfg)
	balanceUseCase := usecase.NewBalanceUseCase(walletRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase)
	walletHandler := handler.NewWalletHandler(topUpUseCase, withdrawalUseCase, balanceUseCase)

	// API routes group with CORS middleware
	api := router.Group("/api/v1")
	api.Use(corsMiddleware())
	{
		// User routes
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)

		// Wallet routes
		api.POST("/topup", middleware.RequireIdempotencyKey(), walletHandler.TopUp)
		api.POST("/withdraw", middleware.RequireIdempotencyKey(), walletHandler.Withdraw)
		api.GET("/balance/:user_id", walletHandler.GetBalance)
	}

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Idempotency-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}