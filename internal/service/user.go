package service

import (
	"context"
	"github.com/krtech-it/gofermart/internal/storage"
)

// UserService реализует бизнес-логику работы с пользователями.
type UserService struct {
	storage storage.UserStorage
}

// NewUserService создаёт новый UserService с переданным хранилищем.
func NewUserService(storage storage.UserStorage) UserServiceInterface {
	return &UserService{
		storage: storage,
	}
}

// UserServiceInterface определяет операции бизнес-логики для работы с пользователями.
type UserServiceInterface interface {
	// CreateUser регистрирует нового пользователя и возвращает токен авторизации.
	CreateUser(ctx context.Context, login, password string) (string, error)
	// Login аутентифицирует пользователя и возвращает токен авторизации.
	Login(ctx context.Context, login, password string) (string, error)
}

// CreateUser регистрирует нового пользователя и возвращает токен авторизации.
func (s *UserService) CreateUser(ctx context.Context, login, password string) (string, error) {
	panic("implement me")
}

// Login аутентифицирует пользователя и возвращает токен авторизации.
func (s *UserService) Login(ctx context.Context, login, password string) (string, error) {
	panic("implement me")
}
