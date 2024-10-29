package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type idempotencyRepository struct {
	db *sql.DB
}

func NewIdempotencyRepository(db *sql.DB) IdempotencyRepository {
	return &idempotencyRepository{db: db}
}

func (r *idempotencyRepository) Save(ctx context.Context, key string, response []byte, expiration time.Time) error {
	query := `
		INSERT INTO idempotency_keys (key, response, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key) DO NOTHING
	`
	
	result, err := r.db.ExecContext(ctx, query, key, response, time.Now(), expiration)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("idempotency key already exists")
	}
	
	return nil
}

func (r *idempotencyRepository) Get(ctx context.Context, key string) ([]byte, error) {
	query := `
		SELECT response
		FROM idempotency_keys
		WHERE key = $1 AND expires_at > NOW()
	`
	
	var response []byte
	err := r.db.QueryRowContext(ctx, query, key).Scan(&response)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return response, nil
}

func (r *idempotencyRepository) Delete(ctx context.Context, key string) error {
	query := `
		DELETE FROM idempotency_keys
		WHERE key = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, key)
	return err
}