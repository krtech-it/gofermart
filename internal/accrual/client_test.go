package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

// newTestClient создаёт Client, направленный на тестовый HTTP-сервер.
func newTestClient(server *httptest.Server) *Client {
	return NewClient(server.URL, zap.NewNop())
}

// --- GetOrderAccrual ---

func TestGetOrderAccrual_200_ValidJSON(t *testing.T) {
	expected := AccrualResponse{Order: "12345", Status: "PROCESSED", Accrual: 99.5}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := newTestClient(server)
	result, err := client.GetOrderAccrual(context.Background(), "12345")

	if err != nil {
		t.Fatalf("ожидался nil, получена ошибка: %v", err)
	}
	if result == nil {
		t.Fatal("ожидался непустой ответ")
	}
	if result.Order != expected.Order {
		t.Errorf("Order: ожидалось %q, получено %q", expected.Order, result.Order)
	}
	if result.Status != expected.Status {
		t.Errorf("Status: ожидалось %q, получено %q", expected.Status, result.Status)
	}
	if result.Accrual != expected.Accrual {
		t.Errorf("Accrual: ожидалось %v, получено %v", expected.Accrual, result.Accrual)
	}
}

func TestGetOrderAccrual_200_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.GetOrderAccrual(context.Background(), "12345")

	if err == nil {
		t.Error("ожидалась ошибка при невалидном JSON")
	}
}

func TestGetOrderAccrual_204(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server)
	result, err := client.GetOrderAccrual(context.Background(), "12345")

	if err != nil {
		t.Fatalf("ожидался nil, получена ошибка: %v", err)
	}
	if result != nil {
		t.Errorf("ожидался nil-результат, получен: %+v", result)
	}
}

func TestGetOrderAccrual_429_WithRetryAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.GetOrderAccrual(context.Background(), "12345")

	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("ожидалась ошибка RateLimitError, получена: %v", err)
	}
	if rateLimitErr.RetryAfter != 30 {
		t.Errorf("RetryAfter: ожидалось 30, получено %d", rateLimitErr.RetryAfter)
	}
}

func TestGetOrderAccrual_429_WithoutRetryAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.GetOrderAccrual(context.Background(), "12345")

	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("ожидалась ошибка RateLimitError, получена: %v", err)
	}
	if rateLimitErr.RetryAfter != 5 {
		t.Errorf("RetryAfter: ожидалось 5 (дефолт), получено %d", rateLimitErr.RetryAfter)
	}
}

func TestGetOrderAccrual_429_NonNumericRetryAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "invalid")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.GetOrderAccrual(context.Background(), "12345")

	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("ожидалась ошибка RateLimitError, получена: %v", err)
	}
	if rateLimitErr.RetryAfter != 5 {
		t.Errorf("RetryAfter: ожидалось 5 (дефолт), получено %d", rateLimitErr.RetryAfter)
	}
}

func TestGetOrderAccrual_UnexpectedStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.GetOrderAccrual(context.Background(), "12345")

	if err == nil {
		t.Error("ожидалась ошибка при неожиданном статус-коде")
	}
}

func TestGetOrderAccrual_CorrectURL(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server)
	client.GetOrderAccrual(context.Background(), "79927398713")

	expected := "/api/orders/79927398713"
	if capturedPath != expected {
		t.Errorf("путь запроса: ожидался %q, получен %q", expected, capturedPath)
	}
}

// --- RateLimitError ---

func TestRateLimitError_ErrorMessage(t *testing.T) {
	err := &RateLimitError{RetryAfter: 10}
	if err.Error() != "rate limit exceeded" {
		t.Errorf("неожиданное сообщение ошибки: %q", err.Error())
	}
}
