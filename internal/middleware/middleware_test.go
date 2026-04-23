package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testSecret = "test-secret-key"

// --- GenerateToken ---

func TestGenerateToken_ReturnsNonEmptyToken(t *testing.T) {
	id := uuid.New()
	token, err := GenerateToken(id, testSecret)

	if err != nil {
		t.Fatalf("ожидался nil, получена ошибка: %v", err)
	}
	if token == "" {
		t.Error("ожидался непустой токен")
	}
}

func TestGenerateToken_DifferentIDsProduceDifferentTokens(t *testing.T) {
	token1, _ := GenerateToken(uuid.New(), testSecret)
	token2, _ := GenerateToken(uuid.New(), testSecret)

	if token1 == token2 {
		t.Error("токены для разных userID должны различаться")
	}
}

// --- AuthMiddleware ---

func buildRouter(secret string) *gin.Engine {
	r := gin.New()
	cfg := config.Config{JWTSecret: secret}
	r.GET("/test", AuthMiddleware(cfg), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestAuthMiddleware_NoCookie_Returns401(t *testing.T) {
	r := buildRouter(testSecret)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken_Returns401(t *testing.T) {
	r := buildRouter(testSecret)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: "invalid.token.value"})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken_Passes(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID, testSecret)
	if err != nil {
		t.Fatalf("не удалось создать токен: %v", err)
	}

	r := buildRouter(testSecret)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: token})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}
}

func TestAuthMiddleware_WrongSecret_Returns401(t *testing.T) {
	token, _ := GenerateToken(uuid.New(), "other-secret")

	r := buildRouter(testSecret)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: token})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидался 401, получен %d", w.Code)
	}
}

func TestAuthMiddleware_SetsUserID(t *testing.T) {
	userID := uuid.New()
	token, _ := GenerateToken(userID, testSecret)

	var gotUserID uuid.UUID
	r := gin.New()
	cfg := config.Config{JWTSecret: testSecret}
	r.GET("/test", AuthMiddleware(cfg), func(c *gin.Context) {
		val, _ := c.Get("userID")
		gotUserID, _ = val.(uuid.UUID)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: token})
	r.ServeHTTP(w, req)

	if gotUserID != userID {
		t.Errorf("ожидался userID %v, получен %v", userID, gotUserID)
	}
}
