package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
)

// OrderStorage определяет операции хранилища для работы с заказами.
type OrderStorage interface {
	// GetAllOrdersByUserID возвращает все заказы пользователя, отсортированные от новых к старым.
	GetAllOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
	// GetOrderByNumber возвращает заказ по его номеру.
	// Возвращает ошибку, если заказ не найден.
	GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error)
	// CreateOrder сохраняет новый заказ и возвращает его с заполненными полями.
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	// UpdateOrder обновляет статус и начисление заказа и возвращает обновлённую запись.
	UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
}

// GetAllOrdersByUserID возвращает все заказы пользователя, отсортированные от новых к старым.
func (p *PostgresStorage) GetAllOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	panic("implement me")
}

// GetOrderByNumber возвращает заказ по его номеру.
// Возвращает ошибку, если заказ не найден.
func (p *PostgresStorage) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	panic("implement me")
}

// CreateOrder сохраняет новый заказ и возвращает его с заполненными полями.
func (p *PostgresStorage) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	panic("implement me")
}

// UpdateOrder обновляет статус и начисление заказа и возвращает обновлённую запись.
func (p *PostgresStorage) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	panic("implement me")
}
