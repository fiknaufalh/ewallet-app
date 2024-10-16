package repository

import (
	"errors"
	"ewallet-app/internal/db"
	"ewallet-app/internal/models"
)

type UserRepository struct {
	db *db.Database
}

func NewUserRepository(db *db.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(id int64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, name, balance FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Balance)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) SaveUser(user *models.User) error {
	result, err := r.db.ExecuteQuery("INSERT INTO users (name, balance) VALUES (?, ?)", user.Name, user.Balance)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *UserRepository) UpdateBalance(id int64, newBalance float64) error {
	_, err := r.db.ExecuteQuery("UPDATE users SET balance = ? WHERE id = ?", newBalance, id)
	return err
}