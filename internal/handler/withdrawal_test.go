package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/service"
)

// --- GetBalance ---

func TestGetBalance_NoUserID(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)

	h.GetBalance(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestGetBalance_InternalError(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return nil, errInternal
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	setUserID(c, uuid.New())

	h.GetBalance(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ожидался 500, получен %d", w.Code)
	}
}

func TestGetBalance_Success(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return &model.Balance{Current: 100, Withdrawn: 50}, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	setUserID(c, uuid.New())

	h.GetBalance(c)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

// --- WithdrawProcess ---

func TestWithdrawProcess_NoUserID(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", nil)

	h.WithdrawProcess(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestWithdrawProcess_BadRequest(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString("not json"))
	c.Request.Header.Set("Content-Type", "application/json")
	setUserID(c, uuid.New())

	h.WithdrawProcess(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался 400, получен %d", w.Code)
	}
}

func TestWithdrawProcess_InvalidOrder(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		withdrawalProcess: func(_ context.Context, _ uuid.UUID, _ string, _ float64) error {
			return service.ErrorInvalidOrderNumber
		},
	})

	body, _ := json.Marshal(map[string]interface{}{"order": "12345", "sum": 50.0})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	setUserID(c, uuid.New())

	h.WithdrawProcess(c)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("ожидался 422, получен %d", w.Code)
	}
}

func TestWithdrawProcess_InsufficientFunds(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		withdrawalProcess: func(_ context.Context, _ uuid.UUID, _ string, _ float64) error {
			return service.ErrorBalanceInsufficientFunds
		},
	})

	body, _ := json.Marshal(map[string]interface{}{"order": "12345", "sum": 50.0})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	setUserID(c, uuid.New())

	h.WithdrawProcess(c)

	if w.Code != http.StatusPaymentRequired {
		t.Errorf("ожидался 402, получен %d", w.Code)
	}
}

func TestWithdrawProcess_Success(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		withdrawalProcess: func(_ context.Context, _ uuid.UUID, _ string, _ float64) error {
			return nil
		},
	})

	body, _ := json.Marshal(map[string]interface{}{"order": "12345", "sum": 50.0})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	setUserID(c, uuid.New())

	h.WithdrawProcess(c)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

// --- GetWithdrawals ---

func TestGetWithdrawals_NoUserID(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)

	h.GetWithdrawals(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestGetWithdrawals_Empty(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		getWithdrawals: func(_ context.Context, _ uuid.UUID) ([]*model.Withdrawal, error) {
			return []*model.Withdrawal{}, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	setUserID(c, uuid.New())

	h.GetWithdrawals(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusNoContent {
		t.Errorf("ожидался 204, получен %d", w.Code)
	}
}

func TestGetWithdrawals_ReturnsList(t *testing.T) {
	withdrawals := []*model.Withdrawal{
		{ID: uuid.New(), Order: "79927398713", Sum: 50, ProcessedAt: time.Now()},
	}
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		getWithdrawals: func(_ context.Context, _ uuid.UUID) ([]*model.Withdrawal, error) {
			return withdrawals, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	setUserID(c, uuid.New())

	h.GetWithdrawals(c)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

func TestGetWithdrawals_InternalError(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{
		getWithdrawals: func(_ context.Context, _ uuid.UUID) ([]*model.Withdrawal, error) {
			return nil, errInternal
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	setUserID(c, uuid.New())

	h.GetWithdrawals(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ожидался 500, получен %d", w.Code)
	}
}
