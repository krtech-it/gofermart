package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
)

var ErrOrderAlreadyByOtherUser = errors.New("order already uploaded by another user")
var ErrOrderAlreadyByThisUser = errors.New("user already have the order")
var ErrorInvalidOrderNumber = errors.New("invalid order number")

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
	GetAllOrdersByUser(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
	// UploadOrder загружает новый заказ от пользователя. Проверяет номер заказа алгоритмом Луна.
	UploadOrder(ctx context.Context, userID uuid.UUID, orderNumber string) error
	// UpdateOrder обновляет статус и начисление существующего заказа.
	UpdateOrder(ctx context.Context, userID uuid.UUID, order *model.Order) error
}

// GetAllOrdersByUser возвращает все заказы пользователя, отсортированные от новых к старым.
func (s *OrderService) GetAllOrdersByUser(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	return s.storage.GetAllOrdersByUserID(ctx, userID)
}

// UploadOrder загружает новый заказ от пользователя. Проверяет номер заказа алгоритмом Луна.
func (s *OrderService) UploadOrder(ctx context.Context, userID uuid.UUID, orderNumber string) error {
	if !isValidLuhn(orderNumber) {
		return ErrorInvalidOrderNumber
	}
	orderDB, err := s.storage.GetOrderByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return err
	}
	if orderDB != nil {
		if orderDB.UserId != userID {
			return ErrOrderAlreadyByOtherUser
		} else {
			return ErrOrderAlreadyByThisUser
		}
	}
	orderNew := &model.Order{
		Number: orderNumber,
		UserId: userID,
		Status: model.OrderStatusNew,
	}
	err = s.storage.CreateOrder(ctx, orderNew)
	return err
}

// UpdateOrder обновляет статус и начисление существующего заказа.
func (s *OrderService) UpdateOrder(ctx context.Context, userID uuid.UUID, order *model.Order) error {
	panic("implement me")
}
