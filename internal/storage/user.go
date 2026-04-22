package storage

import (
	"context"
	"github.com/krtech-it/gofermart/internal/model"
)

// UserStorage определяет операции хранилища для работы с пользователями.
type UserStorage interface {
	// CreateUser сохраняет нового пользователя в базе данных.
	CreateUser(ctx context.Context, user model.User) error
	// GetUserByLogin возвращает пользователя по логину.
	// Возвращает ошибку, если пользователь не найден.
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
}

// CreateUser сохраняет нового пользователя в базе данных.
func (p *PostgresStorage) CreateUser(ctx context.Context, user model.User) error {
	panic("implement me")
}

// GetUserByLogin возвращает пользователя по логину.
// Возвращает ошибку, если пользователь не найден.
func (p *PostgresStorage) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	panic("implement me")
}
