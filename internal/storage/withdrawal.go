package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
)

// WithdrawalStorage определяет операции хранилища для работы со списаниями баллов.
type WithdrawalStorage interface {
	// GetAllWithdrawalsByUserID возвращает все списания пользователя, отсортированные от новых к старым.
	GetAllWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error)
	// CreateWithdrawal создаёт новую операцию списания баллов в счёт указанного заказа.
	CreateWithdrawal(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) (model.Withdrawal, error)
	// GetBalance возвращает текущий баланс и суммарное количество списанных баллов пользователя.
	GetBalance(ctx context.Context, userID uuid.UUID) (model.Balance, error)
}

// GetAllWithdrawalsByUserID возвращает все списания пользователя, отсортированные от новых к старым.
func (p *PostgresStorage) GetAllWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	panic("implement me")
}

// CreateWithdrawal создаёт новую операцию списания баллов в счёт указанного заказа.
func (p *PostgresStorage) CreateWithdrawal(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) (model.Withdrawal, error) {
	panic("implement me")
}

// GetBalance возвращает текущий баланс и суммарное количество списанных баллов пользователя.
func (p *PostgresStorage) GetBalance(ctx context.Context, userID uuid.UUID) (model.Balance, error) {
	panic("implement me")
}
