package http

import (
	"github.com/gin-gonic/gin"
	"github.com/krtech-it/gofermart/internal/config"
	"github.com/krtech-it/gofermart/internal/handler"
	"github.com/krtech-it/gofermart/internal/middleware"
)

func NewRouter(h *handler.Handler, cfg config.Config) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/user")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)

	auth := api.Group("/", middleware.AuthMiddleware(cfg))
	auth.POST("/orders", h.UploadOrder)
	auth.GET("/orders", h.GetOrders)
	auth.GET("/balance", h.GetBalance)
	auth.POST("/balance/withdraw", h.WithdrawProcess)
	auth.GET("/withdrawals", h.GetWithdrawals)
	return r
}
