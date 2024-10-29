package repository

import (
	"context"
	"time"

	"ewallet-app/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *entity.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Wallet, error)
	UpdateBalance(ctx context.Context, wallet *entity.Wallet) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetByReferenceID(ctx context.Context, referenceID string) (*entity.Transaction, error)
	UpdateStatus(ctx context.Context, transaction *entity.Transaction) error
	GetWalletTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
}

type IdempotencyRepository interface {
	Save(ctx context.Context, key string, response []byte, expiration time.Time) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}