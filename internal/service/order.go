package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
)

// OrderService реализует бизнес-логику работы с заказами.
type OrderService struct {
	storage storage.OrderStorage
}

// NewOrderService создаёт новый OrderService с переданным хранилищем.
func NewOrderService(storage storage.OrderStorage) OrderServiceInterface {
	return &OrderService{
		storage: storage,
	}
}

// OrderServiceInterface определяет операции бизнес-логики для работы с заказами.
type OrderServiceInterface interface {
	// GetAllOrdersByUser возвращает все заказы пользователя, отсортированные от новых к старым.
	GetAllOrdersByUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
	// UploadOrder загружает новый заказ от пользователя. Проверяет номер заказа алгоритмом Луна.
	UploadOrder(ctx context.Context, userID uuid.UUID, orderNumber string) error
	// UpdateOrder обновляет статус и начисление существующего заказа.
	UpdateOrder(ctx context.Context, userID uuid.UUID, order *model.Order) error
}

// GetAllOrdersByUser возвращает все заказы пользователя, отсортированные от новых к старым.
func (s *OrderService) GetAllOrdersByUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	panic("implement me")
}

// UploadOrder загружает новый заказ от пользователя. Проверяет номер заказа алгоритмом Луна.
func (s *OrderService) UploadOrder(ctx context.Context, userID uuid.UUID, orderNumber string) error {
	panic("implement me")
}

// UpdateOrder обновляет статус и начисление существующего заказа.
func (s *OrderService) UpdateOrder(ctx context.Context, userID uuid.UUID, order *model.Order) error {
	panic("implement me")
}
