package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ewallet-app/internal/domain/usecase"
)

type WalletHandler struct {
	topUpUseCase     usecase.TopUpUseCase
	withdrawalUseCase usecase.WithdrawalUseCase
	balanceUseCase   usecase.BalanceUseCase
}

func NewWalletHandler(
	topUpUseCase usecase.TopUpUseCase,
	withdrawalUseCase usecase.WithdrawalUseCase,
	balanceUseCase usecase.BalanceUseCase,
) *WalletHandler {
	return &WalletHandler{
		topUpUseCase:     topUpUseCase,
		withdrawalUseCase: withdrawalUseCase,
		balanceUseCase:   balanceUseCase,
	}
}

type TopUpRequest struct {
	UserID      string  `json:"user_id" binding:"required,uuid"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	ReferenceID string  `json:"reference_id" binding:"required"`
}

type WithdrawalRequest struct {
	UserID      string  `json:"user_id" binding:"required,uuid"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	ReferenceID string  `json:"reference_id" binding:"required"`
	BankAccount string  `json:"bank_account" binding:"required"`
}

func (h *WalletHandler) TopUp(c *gin.Context) {
	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	input := usecase.TopUpInput{
		UserID:      userID,
		Amount:      req.Amount,
		ReferenceID: req.ReferenceID,
	}

	output, err := h.topUpUseCase.TopUp(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	var req WithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	input := usecase.WithdrawalInput{
		UserID:      userID,
		Amount:      req.Amount,
		ReferenceID: req.ReferenceID,
		BankAccount: req.BankAccount,
	}

	output, err := h.withdrawalUseCase.Withdraw(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	balance, err := h.balanceUseCase.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}