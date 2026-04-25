package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/krtech-it/gofermart/internal/model"
)

// UserStorage определяет операции хранилища для работы с пользователями.
type UserStorage interface {
	// CreateUser сохраняет нового пользователя в базе данных.
	CreateUser(ctx context.Context, user *model.User) error
	// GetUserByLogin возвращает пользователя по логину.
	// Возвращает ошибку, если пользователь не найден.
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

// CreateUser сохраняет нового пользователя в базе данных.
func (p *PostgresStorage) CreateUser(ctx context.Context, user *model.User) error {
	_, err := p.db.ExecContext(ctx, "insert into users (id, login, password_hash) values ($1, $2, $3)", user.ID, user.Login, user.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByLogin возвращает пользователя по логину.
// Возвращает ошибку, если пользователь не найден.
func (p *PostgresStorage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	userDB := &model.User{}
	row := p.db.QueryRowContext(ctx, "select id, login, password_hash from users where login = $1", login)
	err := row.Scan(&userDB.ID, &userDB.Login, &userDB.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return userDB, nil
}
