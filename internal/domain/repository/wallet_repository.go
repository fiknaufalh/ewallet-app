package repository

import (
	"context"
	"database/sql"
	"errors"

	"ewallet-app/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, balance, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		wallet.ID,
		wallet.UserID,
		wallet.Balance,
		wallet.Version,
		wallet.CreatedAt,
		wallet.UpdatedAt,
	)
	
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return errors.New("wallet already exists for user")
			}
		}
		return err
	}
	
	return nil
}

func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Wallet, error) {
	query := `
		SELECT id, user_id, balance, version, created_at, updated_at
		FROM wallets
		WHERE id = $1
		FOR UPDATE
	`
	
	wallet := &entity.Wallet{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Version,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}
	
	if err != nil {
		return nil, err
	}
	
	return wallet, nil
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Wallet, error) {
	query := `
		SELECT id, user_id, balance, version, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
		FOR UPDATE
	`
	
	wallet := &entity.Wallet{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Version,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}
	
	if err != nil {
		return nil, err
	}
	
	return wallet, nil
}

func (r *walletRepository) UpdateBalance(ctx context.Context, wallet *entity.Wallet) error {
	query := `
		UPDATE wallets
		SET balance = $1, version = $2, updated_at = $3
		WHERE id = $4 AND version = $5
	`
	
	result, err := r.db.ExecContext(ctx, query,
		wallet.Balance,
		wallet.Version,
		wallet.UpdatedAt,
		wallet.ID,
		wallet.Version-1, // Check previous version for optimistic locking
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("concurrent update detected")
	}
	
	return nil
}