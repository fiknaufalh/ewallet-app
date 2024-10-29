package usecase

import (
	"context"

	"github.com/google/uuid"
)

type TopUpInput struct {
	UserID      uuid.UUID `json:"user_id"`
	Amount      float64   `json:"amount"`
	ReferenceID string    `json:"reference_id"`
}

type WithdrawalInput struct {
	UserID      uuid.UUID `json:"user_id"`
	Amount      float64   `json:"amount"`
	ReferenceID string    `json:"reference_id"`
	BankAccount string    `json:"bank_account"`
}

type TransactionOutput struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	Balance       float64   `json:"balance"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserOutput struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserUseCase interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*UserOutput, error)
}

type TopUpUseCase interface {
	TopUp(ctx context.Context, input TopUpInput) (*TransactionOutput, error)
}

type WithdrawalUseCase interface {
	Withdraw(ctx context.Context, input WithdrawalInput) (*TransactionOutput, error)
}

type BalanceUseCase interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (float64, error)
}