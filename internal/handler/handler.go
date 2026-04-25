// Пакет handler содержит HTTP-обработчики системы Gophermart.
package handler

import (
	"github.com/krtech-it/gofermart/internal/service"
	"go.uber.org/zap"
)

// Handler содержит зависимости HTTP-обработчиков.
type Handler struct {
	user       service.UserServiceInterface
	order      service.OrderServiceInterface
	withdrawal service.WithdrawalServiceInterface
	logger     *zap.Logger
}

// NewHandler создаёт новый Handler с переданными сервисами и логгером.
func NewHandler(services *service.Services, logger *zap.Logger) *Handler {
	return &Handler{
		user:       services.User,
		order:      services.Order,
		withdrawal: services.Withdrawal,
		logger:     logger,
	}
}
