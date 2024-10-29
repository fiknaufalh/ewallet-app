package usecase

import (
	"context"
	"database/sql"

	"ewallet-app/internal/domain/entity"
	"ewallet-app/internal/domain/repository"

	"github.com/google/uuid"
)

type userUseCase struct {
	db         *sql.DB
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewUserUseCase(
	db *sql.DB,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
) UserUseCase {
	return &userUseCase{
		db:         db,
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (uc *userUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	// Start transaction
	tx, err := uc.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create user
	user := entity.NewUser(input.Username, input.Email)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create wallet for user
	wallet := entity.NewWallet(user.ID)
	if err := uc.walletRepo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &UserOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (uc *userUseCase) GetUser(ctx context.Context, userID uuid.UUID) (*UserOutput, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}