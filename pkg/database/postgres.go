package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"ewallet-app/internal/config"
)

func NewPostgresConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error

	// Retry logic for database connection
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", cfg.GetDSN())
		if err != nil {
			log.Printf("Failed to open database connection (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(time.Second * 5)
			continue
		}

		// Test connection
		err = db.Ping()
		if err == nil {
			break
		}
		
		log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		db.Close()
		time.Sleep(time.Second * 5)
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to the database after %d attempts: %v", maxRetries, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Println("Successfully connected to database")
	return db, nil
}