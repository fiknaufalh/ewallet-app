package repository

import (
	"context"
	"database/sql"
	"errors"

	"ewallet-app/internal/domain/entity"

	"github.com/google/uuid"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		INSERT INTO transactions (id, wallet_id, type, amount, status, reference_id, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.WalletID,
		transaction.Type,
		transaction.Amount,
		transaction.Status,
		transaction.ReferenceID,
		transaction.Description,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	query := `
		SELECT id, wallet_id, type, amount, status, reference_id, description, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`
	
	tx := &entity.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID,
		&tx.WalletID,
		&tx.Type,
		&tx.Amount,
		&tx.Status,
		&tx.ReferenceID,
		&tx.Description,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}
	
	if err != nil {
		return nil, err
	}
	
	return tx, nil
}

func (r *transactionRepository) GetByReferenceID(ctx context.Context, referenceID string) (*entity.Transaction, error) {
	query := `
		SELECT id, wallet_id, type, amount, status, reference_id, description, created_at, updated_at
		FROM transactions
		WHERE reference_id = $1
	`
	
	tx := &entity.Transaction{}
	err := r.db.QueryRowContext(ctx, query, referenceID).Scan(
		&tx.ID,
		&tx.WalletID,
		&tx.Type,
		&tx.Amount,
		&tx.Status,
		&tx.ReferenceID,
		&tx.Description,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}
	
	if err != nil {
		return nil, err
	}
	
	return tx, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	
	result, err := r.db.ExecContext(ctx, query,
		transaction.Status,
		transaction.UpdatedAt,
		transaction.ID,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}
	
	return nil
}

func (r *transactionRepository) GetWalletTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	query := `
		SELECT id, wallet_id, type, amount, status, reference_id, description, created_at, updated_at
		FROM transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, walletID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var transactions []*entity.Transaction
	for rows.Next() {
		tx := &entity.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.WalletID,
			&tx.Type,
			&tx.Amount,
			&tx.Status,
			&tx.ReferenceID,
			&tx.Description,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return transactions, nil
}