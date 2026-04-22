package handler

import (
	"github.com/krtech-it/gofermart/internal/service"
)

type Handler struct {
	user       service.UserServiceInterface
	order      service.OrderServiceInterface
	withdrawal service.WithdrawalServiceInterface
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		user:       services.User,
		order:      services.Order,
		withdrawal: services.Withdrawal,
	}
}
