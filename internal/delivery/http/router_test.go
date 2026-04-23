package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/config"
	"github.com/krtech-it/gofermart/internal/handler"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/service"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- заглушки сервисов ---

type stubUser struct{}

func (s *stubUser) CreateUser(_ context.Context, _, _ string) (string, error) { return "", nil }
func (s *stubUser) Login(_ context.Context, _, _ string) (string, error)      { return "", nil }

type stubOrder struct{}

func (s *stubOrder) GetAllOrdersByUser(_ context.Context, _ uuid.UUID) ([]*model.Order, error) {
	return nil, nil
}
func (s *stubOrder) UploadOrder(_ context.Context, _ uuid.UUID, _ string) error { return nil }
func (s *stubOrder) UpdateOrder(_ context.Context, _ uuid.UUID, _ *model.Order) error {
	return nil
}

type stubWithdrawal struct{}

func (s *stubWithdrawal) WithdrawalProcess(_ context.Context, _ uuid.UUID, _ string, _ float64) error {
	return nil
}
func (s *stubWithdrawal) GetWithdrawals(_ context.Context, _ uuid.UUID) ([]*model.Withdrawal, error) {
	return nil, nil
}
func (s *stubWithdrawal) GetBalance(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
	return &model.Balance{}, nil
}

func buildTestRouter() *gin.Engine {
	services := &service.Services{
		User:       &stubUser{},
		Order:      &stubOrder{},
		Withdrawal: &stubWithdrawal{},
	}
	h := handler.NewHandler(services, zap.NewNop())
	cfg := config.Config{
		RunAddress: "localhost:8080",
		JWTSecret:  "test-secret",
	}
	return NewRouter(h, cfg)
}

// --- тесты маршрутов ---

func TestRouter_RegisterRoute_Exists(t *testing.T) {
	r := buildTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", nil)

	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Error("маршрут POST /api/user/register не зарегистрирован")
	}
}

func TestRouter_LoginRoute_Exists(t *testing.T) {
	r := buildTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", nil)

	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Error("маршрут POST /api/user/login не зарегистрирован")
	}
}

func TestRouter_ProtectedOrders_RequiresAuth(t *testing.T) {
	r := buildTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401 без авторизации, получен %d", w.Code)
	}
}

func TestRouter_ProtectedGetOrders_RequiresAuth(t *testing.T) {
	r := buildTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401 без авторизации, получен %d", w.Code)
	}
}

func TestRouter_ProtectedBalance_RequiresAuth(t *testing.T) {
	r := buildTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401 без авторизации, получен %d", w.Code)
	}
}

func TestRouter_ProtectedWithdraw_RequiresAuth(t *testing.T) {
	r := buildTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401 без авторизации, получен %d", w.Code)
	}
}

func TestRouter_ProtectedWithdrawals_RequiresAuth(t *testing.T) {
	r := buildTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401 без авторизации, получен %d", w.Code)
	}
}
