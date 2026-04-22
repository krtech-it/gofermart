package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
)

// WithdrawalService реализует бизнес-логику работы со списаниями баллов.
type WithdrawalService struct {
	storage storage.WithdrawalStorage
}

// NewWithdrawalService создаёт новый WithdrawalService с переданным хранилищем.
func NewWithdrawalService(storage storage.WithdrawalStorage) WithdrawalServiceInterface {
	return &WithdrawalService{
		storage: storage,
	}
}

// WithdrawalServiceInterface определяет операции бизнес-логики для работы со списаниями баллов.
type WithdrawalServiceInterface interface {
	// WithdrawalProcess выполняет списание баллов в счёт указанного заказа.
	// Возвращает ошибку, если баланса недостаточно.
	WithdrawalProcess(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error
	// GetWithdrawals возвращает все операции списания пользователя, отсортированные от новых к старым.
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error)
	// GetBalance возвращает текущий баланс и суммарное количество списанных баллов пользователя.
	GetBalance(ctx context.Context, userID uuid.UUID) (model.Balance, error)
}

// WithdrawalProcess выполняет списание баллов в счёт указанного заказа.
// Возвращает ошибку, если баланса недостаточно.
func (s *WithdrawalService) WithdrawalProcess(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error {
	panic("implement me")
}

// GetWithdrawals возвращает все операции списания пользователя, отсортированные от новых к старым.
func (s *WithdrawalService) GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	panic("implement me")
}

// GetBalance возвращает текущий баланс и суммарное количество списанных баллов пользователя.
func (s *WithdrawalService) GetBalance(ctx context.Context, userID uuid.UUID) (model.Balance, error) {
	panic("implement me")
}
