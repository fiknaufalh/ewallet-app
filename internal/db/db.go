package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"
)

type Database struct {
	*sql.DB
	cb *gobreaker.CircuitBreaker
}

func NewDatabase(dsn string) (*Database, error) {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "database",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})

	var db *sql.DB
	operation := func() error {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		return db.Ping()
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{DB: db, cb: cb}, nil
}

func (d *Database) ExecuteQuery(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.cb.Execute(func() (interface{}, error) {
		return d.DB.Exec(query, args...)
	})

	if err != nil {
		return nil, fmt.Errorf("circuit breaker error: %w", err)
	}

	return result.(sql.Result), nil
}

func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	row, _ := d.cb.Execute(func() (interface{}, error) {
		return d.DB.QueryRow(query, args...), nil
	})

	return row.(*sql.Row)
}

// Add other methods as needed (e.g., Query, Prepare, etc.)