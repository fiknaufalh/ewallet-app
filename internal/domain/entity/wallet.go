package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount      = errors.New("invalid amount")
)

type Wallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   float64   `json:"balance"`
	Version   int       `json:"version"` // For optimistic locking
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewWallet(userID uuid.UUID) *Wallet {
	return &Wallet{
		ID:        uuid.New(),
		UserID:    userID,
		Balance:   0,
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (w *Wallet) Credit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	w.Balance += amount
	w.Version++
	w.UpdatedAt = time.Now()
	return nil
}

func (w *Wallet) Debit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if w.Balance < amount {
		return ErrInsufficientBalance
	}
	w.Balance -= amount
	w.Version++
	w.UpdatedAt = time.Now()
	return nil
}