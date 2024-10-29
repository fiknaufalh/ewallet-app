package usecase

import (
	"context"

	"ewallet-app/internal/domain/repository"

	"github.com/google/uuid"
)

type balanceUseCase struct {
	walletRepo repository.WalletRepository
}

func NewBalanceUseCase(walletRepo repository.WalletRepository) BalanceUseCase {
	return &balanceUseCase{
		walletRepo: walletRepo,
	}
}

func (uc *balanceUseCase) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	wallet, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}