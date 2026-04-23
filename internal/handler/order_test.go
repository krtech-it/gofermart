package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/service"
)

// --- UploadOrder ---

func TestUploadOrder_NoUserID(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("79927398713"))

	h.UploadOrder(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestUploadOrder_InvalidLuhn(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{
		uploadOrder: func(_ context.Context, _ uuid.UUID, _ string) error {
			return service.ErrorInvalidOrderNumber
		},
	}, &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("79927398710"))
	setUserID(c, uuid.New())

	h.UploadOrder(c)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("ожидался 422, получен %d", w.Code)
	}
}

func TestUploadOrder_AlreadyByThisUser(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{
		uploadOrder: func(_ context.Context, _ uuid.UUID, _ string) error {
			return service.ErrOrderAlreadyByThisUser
		},
	}, &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("79927398713"))
	setUserID(c, uuid.New())

	h.UploadOrder(c)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

func TestUploadOrder_AlreadyByOtherUser(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{
		uploadOrder: func(_ context.Context, _ uuid.UUID, _ string) error {
			return service.ErrOrderAlreadyByOtherUser
		},
	}, &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("79927398713"))
	setUserID(c, uuid.New())

	h.UploadOrder(c)

	if w.Code != http.StatusConflict {
		t.Errorf("ожидался 409, получен %d", w.Code)
	}
}

func TestUploadOrder_Success(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{
		uploadOrder: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}, &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("79927398713"))
	setUserID(c, uuid.New())

	h.UploadOrder(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusAccepted {
		t.Errorf("ожидался 202, получен %d", w.Code)
	}
}

func TestUploadOrder_InternalError(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{
		uploadOrder: func(_ context.Context, _ uuid.UUID, _ string) error {
			return errInternal
		},
	}, &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("79927398713"))
	setUserID(c, uuid.New())

	h.UploadOrder(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ожидался 500, получен %d", w.Code)
	}
}

// --- GetOrders ---

func TestGetOrders_NoUserID(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{}, &mockWithdrawalService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)

	h.GetOrders(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestGetOrders_Empty(t *testing.T) {
	h := newTestHandler(&mockUserService{}, ordersFromModels(nil), &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	setUserID(c, uuid.New())

	h.GetOrders(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusNoContent {
		t.Errorf("ожидался 204, получен %d", w.Code)
	}
}

func TestGetOrders_ReturnsList(t *testing.T) {
	orders := []*model.Order{
		{Number: "79927398713", Status: model.OrderStatusNew, UploadedAt: time.Now()},
	}
	h := newTestHandler(&mockUserService{}, ordersFromModels(orders), &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	setUserID(c, uuid.New())

	h.GetOrders(c)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

func TestGetOrders_InternalError(t *testing.T) {
	h := newTestHandler(&mockUserService{}, &mockOrderService{
		getAllOrdersByUser: func(_ context.Context, _ uuid.UUID) ([]*model.Order, error) {
			return nil, errInternal
		},
	}, &mockWithdrawalService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	setUserID(c, uuid.New())

	h.GetOrders(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ожидался 500, получен %d", w.Code)
	}
}
