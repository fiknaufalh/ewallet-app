package services

import (
	"errors"
	"ewallet-app/internal/models"
	"ewallet-app/internal/repository"
)

type WalletService struct {
	userRepo *repository.UserRepository
}

func NewWalletService(userRepo *repository.UserRepository) *WalletService {
	return &WalletService{
		userRepo: userRepo,
	}
}

func (s *WalletService) CreateUser(user *models.User) error {
	return s.userRepo.SaveUser(user)
}

func (s *WalletService) GetUser(id int64) (*models.User, error) {
	return s.userRepo.GetUser(id)
}

func (s *WalletService) TopUp(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("top-up amount must be positive")
	}

	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}

	newBalance := user.Balance + amount
	return s.userRepo.UpdateBalance(userID, newBalance)
}

func (s *WalletService) Withdraw(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}

	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}

	if user.Balance < amount {
		return errors.New("insufficient balance")
	}

	newBalance := user.Balance - amount
	return s.userRepo.UpdateBalance(userID, newBalance)
}

func (s *WalletService) GetBalance(userID int64) (float64, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return 0, err
	}

	return user.Balance, nil
}