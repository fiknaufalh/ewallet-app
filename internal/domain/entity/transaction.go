package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeTopup      TransactionType = "topup"
	TransactionTypeWithdrawal TransactionType = "withdrawal"

	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID          uuid.UUID         `json:"id"`
	WalletID    uuid.UUID         `json:"wallet_id"`
	Type        TransactionType   `json:"type"`
	Amount      float64           `json:"amount"`
	Status      TransactionStatus `json:"status"`
	ReferenceID string           `json:"reference_id"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func NewTransaction(walletID uuid.UUID, txType TransactionType, amount float64, referenceID, description string) *Transaction {
	return &Transaction{
		ID:          uuid.New(),
		WalletID:    walletID,
		Type:        txType,
		Amount:      amount,
		Status:      TransactionStatusPending,
		ReferenceID: referenceID,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (t *Transaction) Complete() {
	t.Status = TransactionStatusCompleted
	t.UpdatedAt = time.Now()
}

func (t *Transaction) Fail() {
	t.Status = TransactionStatusFailed
	t.UpdatedAt = time.Now()
}