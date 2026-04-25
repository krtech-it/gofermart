package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/service"
)

func init() {
	gin.SetMode(gin.TestMode)
}

var errInternal = errors.New("internal error")

func TestLogin_BadRequest(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString("not json"))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался 400, получен %d", w.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	h := newTestHandler(&mockUserService{
		login: func(_ context.Context, _, _ string) (string, error) {
			return "", service.ErrorInvalidLoginPassword
		},
	}, &mockOrderService{}, &mockWithdrawalService{})

	body, _ := json.Marshal(map[string]string{"login": "user", "password": "wrong"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	h := newTestHandler(&mockUserService{
		login: func(_ context.Context, _, _ string) (string, error) {
			return "token123", nil
		},
	}, &mockOrderService{}, &mockWithdrawalService{})

	body, _ := json.Marshal(map[string]string{"login": "user", "password": "pass"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

func TestLogin_InternalError(t *testing.T) {
	h := newTestHandler(&mockUserService{
		login: func(_ context.Context, _, _ string) (string, error) {
			return "", errInternal
		},
	}, &mockOrderService{}, &mockWithdrawalService{})

	body, _ := json.Marshal(map[string]string{"login": "user", "password": "pass"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ожидался 500, получен %d", w.Code)
	}
}

func TestRegister_BadRequest(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString("not json"))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался 400, получен %d", w.Code)
	}
}

func TestRegister_LoginAlreadyExists(t *testing.T) {
	h := newTestHandler(&mockUserService{
		createUser: func(_ context.Context, _, _ string) (string, error) {
			return "", service.ErrorLoginAlreadyExists
		},
	}, &mockOrderService{}, &mockWithdrawalService{})

	body, _ := json.Marshal(map[string]string{"login": "user", "password": "pass"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)

	if w.Code != http.StatusConflict {
		t.Errorf("ожидался 409, получен %d", w.Code)
	}
}

func TestRegister_Success(t *testing.T) {
	h := newTestHandler(&mockUserService{
		createUser: func(_ context.Context, _, _ string) (string, error) {
			return "token123", nil
		},
	}, &mockOrderService{}, &mockWithdrawalService{})

	body, _ := json.Marshal(map[string]string{"login": "user", "password": "pass"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

func TestRegister_InternalError(t *testing.T) {
	h := newTestHandler(&mockUserService{
		createUser: func(_ context.Context, _, _ string) (string, error) {
			return "", errInternal
		},
	}, &mockOrderService{}, &mockWithdrawalService{})

	body, _ := json.Marshal(map[string]string{"login": "user", "password": "pass"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ожидался 500, получен %d", w.Code)
	}
}

// setUserID добавляет userID в gin.Context (имитирует прошедший middleware).
func setUserID(c *gin.Context, id uuid.UUID) {
	c.Set("userID", id)
}

// setOrders добавляет в контекст список заказов через мок-сервис.
func ordersFromModels(orders []*model.Order) *mockOrderService {
	return &mockOrderService{
		getAllOrdersByUser: func(_ context.Context, _ uuid.UUID) ([]*model.Order, error) {
			return orders, nil
		},
	}
}
