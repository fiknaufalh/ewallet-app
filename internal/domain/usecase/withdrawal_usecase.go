package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"ewallet-app/internal/config"
	"ewallet-app/internal/domain/entity"
	"ewallet-app/internal/domain/repository"

	"github.com/google/uuid"
)

type withdrawalUseCase struct {
	db                *sql.DB
	walletRepo        repository.WalletRepository
	transactionRepo   repository.TransactionRepository
	idempotencyRepo   repository.IdempotencyRepository
	config           *config.Config
}

func NewWithdrawalUseCase(
	db *sql.DB,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	idempotencyRepo repository.IdempotencyRepository,
	config *config.Config,
) WithdrawalUseCase {
	return &withdrawalUseCase{
		db:              db,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		idempotencyRepo: idempotencyRepo,
		config:         config,
	}
}

func (uc *withdrawalUseCase) Withdraw(ctx context.Context, input WithdrawalInput) (*TransactionOutput, error) {
	// Validate amount
	if input.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	
	if input.Amount > uc.config.Transaction.MaxWithdrawalAmount {
		return nil, fmt.Errorf("amount exceeds maximum allowed (%f)", uc.config.Transaction.MaxWithdrawalAmount)
	}
	
	if input.Amount < uc.config.Transaction.MinWithdrawalAmount {
		return nil, fmt.Errorf("amount below minimum allowed (%f)", uc.config.Transaction.MinWithdrawalAmount)
	}

	// Validate bank account
	if input.BankAccount == "" {
		return nil, errors.New("bank account is required")
	}

	// Check idempotency
	if resp, err := uc.idempotencyRepo.Get(ctx, input.ReferenceID); err == nil && resp != nil {
		var output TransactionOutput
		if err := json.Unmarshal(resp, &output); err != nil {
			return nil, err
		}
		return &output, nil
	}

	// Start transaction
	tx, err := uc.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get wallet
	wallet, err := uc.walletRepo.GetByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := entity.NewTransaction(
		wallet.ID,
		entity.TransactionTypeWithdrawal,
		input.Amount,
		input.ReferenceID,
		fmt.Sprintf("Withdrawal to bank account %s", input.BankAccount),
	)

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Update wallet balance
	if err := wallet.Debit(input.Amount); err != nil {
		return nil, err
	}

	// Update wallet in database
	if err := uc.walletRepo.UpdateBalance(ctx, wallet); err != nil {
		return nil, err
	}

	// Update transaction status
	transaction.Complete()
	if err := uc.transactionRepo.UpdateStatus(ctx, transaction); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Create output
	output := &TransactionOutput{
		TransactionID: transaction.ID,
		Status:        string(transaction.Status),
		Balance:       wallet.Balance,
	}

	// Save idempotency key
	responseBytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	expiration := time.Now().Add(uc.config.Security.IdempotencyKeyExpiration)
	if err := uc.idempotencyRepo.Save(ctx, input.ReferenceID, responseBytes, expiration); err != nil {
		// Log error but continue since the transaction was successful
		log.Printf("Error saving idempotency key: %v", err)
	}

	return output, nil
}

func (uc *withdrawalUseCase) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	wallet, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}