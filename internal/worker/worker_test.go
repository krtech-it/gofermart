package worker

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/accrual"
	"github.com/krtech-it/gofermart/internal/model"
	"go.uber.org/zap"
)

// --- мок OrderStorage ---

type mockOrderStorage struct {
	getAllOpenOrders func(ctx context.Context) ([]*model.Order, error)
	updateOrder      func(ctx context.Context, order *model.Order) error
	// остальные методы интерфейса — заглушки
}

func (m *mockOrderStorage) GetAllOpenOrders(ctx context.Context) ([]*model.Order, error) {
	return m.getAllOpenOrders(ctx)
}
func (m *mockOrderStorage) UpdateOrder(ctx context.Context, order *model.Order) error {
	return m.updateOrder(ctx, order)
}
func (m *mockOrderStorage) GetAllOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	return nil, nil
}
func (m *mockOrderStorage) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	return nil, nil
}
func (m *mockOrderStorage) CreateOrder(ctx context.Context, order *model.Order) error { return nil }

// newTestWorker создаёт Worker с мок-хранилищем и accrual-клиентом на тестовый сервер.
func newTestWorker(storage *mockOrderStorage, server *httptest.Server) *Worker {
	hostPort := strings.TrimPrefix(server.URL, "http://")
	client := accrual.NewClient(hostPort, zap.NewNop())
	return NewWorker(storage, client, zap.NewNop())
}

// accrualServer создаёт тестовый HTTP-сервер, возвращающий заданный статус и тело.
func accrualServer(status int, body interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if body != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(body)
		} else {
			w.WriteHeader(status)
		}
	}))
}

// --- processOrders ---

func TestProcessOrders_StorageError(t *testing.T) {
	storageErr := errors.New("db error")
	server := accrualServer(http.StatusNoContent, nil)
	defer server.Close()

	w := newTestWorker(&mockOrderStorage{
		getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
			return nil, storageErr
		},
	}, server)

	// Не должно паниковать, должно молча вернуться.
	w.processOrders(context.Background())
}

func TestProcessOrders_EmptyOrders(t *testing.T) {
	server := accrualServer(http.StatusNoContent, nil)
	defer server.Close()

	updateCalled := false
	w := newTestWorker(&mockOrderStorage{
		getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
			return []*model.Order{}, nil
		},
		updateOrder: func(_ context.Context, _ *model.Order) error {
			updateCalled = true
			return nil
		},
	}, server)

	w.processOrders(context.Background())

	if updateCalled {
		t.Error("UpdateOrder не должен вызываться при пустом списке заказов")
	}
}

func TestProcessOrders_AccrualReturns204_NoUpdate(t *testing.T) {
	server := accrualServer(http.StatusNoContent, nil)
	defer server.Close()

	updateCalled := false
	w := newTestWorker(&mockOrderStorage{
		getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
			return []*model.Order{{Number: "79927398713", Status: model.OrderStatusNew}}, nil
		},
		updateOrder: func(_ context.Context, _ *model.Order) error {
			updateCalled = true
			return nil
		},
	}, server)

	w.processOrders(context.Background())

	if updateCalled {
		t.Error("UpdateOrder не должен вызываться при ответе 204 от accrual")
	}
}

func TestProcessOrders_AccrualReturnsProcessed_UpdatesCalled(t *testing.T) {
	accrualResp := accrual.AccrualResponse{
		Order:   "79927398713",
		Status:  "PROCESSED",
		Accrual: 99.5,
	}
	server := accrualServer(http.StatusOK, accrualResp)
	defer server.Close()

	var updatedOrder *model.Order
	w := newTestWorker(&mockOrderStorage{
		getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
			return []*model.Order{{Number: "79927398713", Status: model.OrderStatusNew}}, nil
		},
		updateOrder: func(_ context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}, server)

	w.processOrders(context.Background())

	if updatedOrder == nil {
		t.Fatal("ожидался вызов UpdateOrder")
	}
	if updatedOrder.Status != model.OrderStatusProcessed {
		t.Errorf("статус: ожидался PROCESSED, получен %v", updatedOrder.Status)
	}
	if updatedOrder.Accrual == nil || *updatedOrder.Accrual != 99.5 {
		t.Errorf("начисление: ожидалось 99.5, получено %v", updatedOrder.Accrual)
	}
}

func TestProcessOrders_AccrualReturnsRegistered_SetsProcessing(t *testing.T) {
	accrualResp := accrual.AccrualResponse{
		Order:  "79927398713",
		Status: "REGISTERED",
	}
	server := accrualServer(http.StatusOK, accrualResp)
	defer server.Close()

	var updatedOrder *model.Order
	w := newTestWorker(&mockOrderStorage{
		getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
			return []*model.Order{{Number: "79927398713", Status: model.OrderStatusNew}}, nil
		},
		updateOrder: func(_ context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}, server)

	w.processOrders(context.Background())

	if updatedOrder == nil {
		t.Fatal("ожидался вызов UpdateOrder")
	}
	if updatedOrder.Status != model.OrderStatusProcessing {
		t.Errorf("статус: ожидался PROCESSING, получен %v", updatedOrder.Status)
	}
}

func TestProcessOrders_UpdateOrderError_Continues(t *testing.T) {
	accrualResp := accrual.AccrualResponse{Order: "79927398713", Status: "PROCESSED", Accrual: 10}
	server := accrualServer(http.StatusOK, accrualResp)
	defer server.Close()

	updateCallCount := 0
	orders := []*model.Order{
		{Number: "79927398713", Status: model.OrderStatusNew},
		{Number: "79927398713", Status: model.OrderStatusNew},
	}
	w := newTestWorker(&mockOrderStorage{
		getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
			return orders, nil
		},
		updateOrder: func(_ context.Context, _ *model.Order) error {
			updateCallCount++
			return errors.New("update failed")
		},
	}, server)

	// Не должно останавливаться на ошибке — должно обработать оба заказа.
	w.processOrders(context.Background())

	if updateCallCount != 2 {
		t.Errorf("ожидалось 2 вызова UpdateOrder, получено %d", updateCallCount)
	}
}

func TestProcessOrders_RateLimit_Returns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	updateCalled := false
	orders := []*model.Order{
		{Number: "79927398713", Status: model.OrderStatusNew},
		{Number: "79927398713", Status: model.OrderStatusNew},
	}
	hostPort := strings.TrimPrefix(server.URL, "http://")
	client := accrual.NewClient(hostPort, zap.NewNop())
	wkr := &Worker{
		storage: &mockOrderStorage{
			getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
				return orders, nil
			},
			updateOrder: func(_ context.Context, _ *model.Order) error {
				updateCalled = true
				return nil
			},
		},
		accrual: client,
		logger:  zap.NewNop(),
	}

	// При 429 worker должен вернуться, не обработав второй заказ.
	wkr.processOrders(context.Background())

	if updateCalled {
		t.Error("UpdateOrder не должен вызываться при rate limit")
	}
}

// --- Start ---

func TestStart_StopsOnContextCancel(t *testing.T) {
	server := accrualServer(http.StatusNoContent, nil)
	defer server.Close()

	w := newTestWorker(&mockOrderStorage{
		getAllOpenOrders: func(_ context.Context) ([]*model.Order, error) {
			return nil, nil
		},
	}, server)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Start(ctx)
		close(done)
	}()

	cancel()
	<-done // горутина должна завершиться после отмены контекста
}
