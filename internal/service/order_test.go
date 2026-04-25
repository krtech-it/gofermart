package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
	"go.uber.org/zap"
)

// mockOrderStorage — ручной мок OrderStorage с настраиваемыми функциями.
type mockOrderStorage struct {
	getAllOrdersByUserID func(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
	getOrderByNumber     func(ctx context.Context, orderNumber string) (*model.Order, error)
	createOrder          func(ctx context.Context, order *model.Order) error
	updateOrder          func(ctx context.Context, order *model.Order) error
	getAllOpenOrders     func(ctx context.Context) ([]*model.Order, error)
}

func (m *mockOrderStorage) GetAllOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	return m.getAllOrdersByUserID(ctx, userID)
}
func (m *mockOrderStorage) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	return m.getOrderByNumber(ctx, orderNumber)
}
func (m *mockOrderStorage) CreateOrder(ctx context.Context, order *model.Order) error {
	return m.createOrder(ctx, order)
}
func (m *mockOrderStorage) UpdateOrder(ctx context.Context, order *model.Order) error {
	return m.updateOrder(ctx, order)
}
func (m *mockOrderStorage) GetAllOpenOrders(ctx context.Context) ([]*model.Order, error) {
	return m.getAllOpenOrders(ctx)
}

// validLuhn — корректный номер заказа по алгоритму Луна.
const validLuhn = "79927398713"

// invalidLuhn — некорректный номер заказа.
const invalidLuhn = "79927398710"

func newTestOrderService(mock *mockOrderStorage) OrderServiceInterface {
	return NewOrderService(mock, zap.NewNop())
}

// --- UploadOrder ---

func TestUploadOrder_InvalidLuhn(t *testing.T) {
	svc := newTestOrderService(&mockOrderStorage{})

	err := svc.UploadOrder(context.Background(), uuid.New(), invalidLuhn)

	if !errors.Is(err, ErrorInvalidOrderNumber) {
		t.Errorf("ожидалась ошибка ErrorInvalidOrderNumber, получена: %v", err)
	}
}

func TestUploadOrder_NewOrder_Success(t *testing.T) {
	svc := newTestOrderService(&mockOrderStorage{
		getOrderByNumber: func(_ context.Context, _ string) (*model.Order, error) {
			return nil, storage.ErrNotFound
		},
		createOrder: func(_ context.Context, order *model.Order) error {
			if order.Status != model.OrderStatusNew {
				t.Errorf("ожидался статус NEW, получен: %v", order.Status)
			}
			return nil
		},
	})

	err := svc.UploadOrder(context.Background(), uuid.New(), validLuhn)

	if err != nil {
		t.Errorf("ожидался nil, получена ошибка: %v", err)
	}
}

func TestUploadOrder_OrderAlreadyByThisUser(t *testing.T) {
	userID := uuid.New()
	svc := newTestOrderService(&mockOrderStorage{
		getOrderByNumber: func(_ context.Context, _ string) (*model.Order, error) {
			return &model.Order{Number: validLuhn, UserId: userID}, nil
		},
	})

	err := svc.UploadOrder(context.Background(), userID, validLuhn)

	if !errors.Is(err, ErrOrderAlreadyByThisUser) {
		t.Errorf("ожидалась ошибка ErrOrderAlreadyByThisUser, получена: %v", err)
	}
}

func TestUploadOrder_OrderAlreadyByOtherUser(t *testing.T) {
	svc := newTestOrderService(&mockOrderStorage{
		getOrderByNumber: func(_ context.Context, _ string) (*model.Order, error) {
			return &model.Order{Number: validLuhn, UserId: uuid.New()}, nil
		},
	})

	err := svc.UploadOrder(context.Background(), uuid.New(), validLuhn)

	if !errors.Is(err, ErrOrderAlreadyByOtherUser) {
		t.Errorf("ожидалась ошибка ErrOrderAlreadyByOtherUser, получена: %v", err)
	}
}

func TestUploadOrder_StorageError_OnGet(t *testing.T) {
	storageErr := errors.New("db error")
	svc := newTestOrderService(&mockOrderStorage{
		getOrderByNumber: func(_ context.Context, _ string) (*model.Order, error) {
			return nil, storageErr
		},
	})

	err := svc.UploadOrder(context.Background(), uuid.New(), validLuhn)

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}

func TestUploadOrder_StorageError_OnCreate(t *testing.T) {
	storageErr := errors.New("insert failed")
	svc := newTestOrderService(&mockOrderStorage{
		getOrderByNumber: func(_ context.Context, _ string) (*model.Order, error) {
			return nil, storage.ErrNotFound
		},
		createOrder: func(_ context.Context, _ *model.Order) error {
			return storageErr
		},
	})

	err := svc.UploadOrder(context.Background(), uuid.New(), validLuhn)

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}

// --- GetAllOrdersByUser ---

func TestGetAllOrdersByUser_ReturnsList(t *testing.T) {
	userID := uuid.New()
	expected := []*model.Order{
		{Number: validLuhn, UserId: userID, Status: model.OrderStatusNew},
	}
	svc := newTestOrderService(&mockOrderStorage{
		getAllOrdersByUserID: func(_ context.Context, _ uuid.UUID) ([]*model.Order, error) {
			return expected, nil
		},
	})

	result, err := svc.GetAllOrdersByUser(context.Background(), userID)

	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(result) != len(expected) {
		t.Errorf("ожидалось %d заказов, получено %d", len(expected), len(result))
	}
}

func TestGetAllOrdersByUser_StorageError(t *testing.T) {
	storageErr := errors.New("db error")
	svc := newTestOrderService(&mockOrderStorage{
		getAllOrdersByUserID: func(_ context.Context, _ uuid.UUID) ([]*model.Order, error) {
			return nil, storageErr
		},
	})

	_, err := svc.GetAllOrdersByUser(context.Background(), uuid.New())

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}
