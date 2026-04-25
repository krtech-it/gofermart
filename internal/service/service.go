// Пакет service содержит бизнес-логику системы лояльности Gophermart.
package service

import (
	"github.com/krtech-it/gofermart/internal/config"
	"github.com/krtech-it/gofermart/internal/storage"
	"go.uber.org/zap"
)

// Services объединяет все сервисы приложения в одну структуру.
type Services struct {
	// User — сервис для работы с пользователями.
	User UserServiceInterface
	// Order — сервис для работы с заказами.
	Order OrderServiceInterface
	// Withdrawal — сервис для работы со списаниями баллов.
	Withdrawal WithdrawalServiceInterface
}

// NewServices создаёт и возвращает Services с инициализированными зависимостями.
func NewServices(userStorage storage.UserStorage, orderStorage storage.OrderStorage, withdrawalStorage storage.WithdrawalStorage, cfg config.Config, logger *zap.Logger) *Services {
	return &Services{
		User:       NewUserService(userStorage, cfg.JWTSecret, logger),
		Order:      NewOrderService(orderStorage, logger),
		Withdrawal: NewWithdrawalService(withdrawalStorage, logger),
	}
}
