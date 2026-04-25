package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/middleware"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrorLoginAlreadyExists = errors.New("login already exists")
var ErrorInvalidLoginPassword = errors.New("invalid login password")

// UserService реализует бизнес-логику работы с пользователями.
type UserService struct {
	storage   storage.UserStorage
	jwtSecret string
	logger    *zap.Logger
}

// NewUserService создаёт новый UserService с переданным хранилищем и логгером.
func NewUserService(storage storage.UserStorage, jwtSecret string, logger *zap.Logger) UserServiceInterface {
	return &UserService{
		storage:   storage,
		jwtSecret: jwtSecret,
		logger:    logger,
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
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return "", err
	}
	if user != nil {
		return "", ErrorLoginAlreadyExists
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	var newUser = model.User{
		ID:           uuid.New(),
		Login:        login,
		PasswordHash: string(hashPassword),
	}
	err = s.storage.CreateUser(ctx, &newUser)
	if err != nil {
		return "", err
	}
	token, err := middleware.GenerateToken(newUser.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Login аутентифицирует пользователя и возвращает токен авторизации.
func (s *UserService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return "", ErrorInvalidLoginPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrorInvalidLoginPassword
	}
	token, err := middleware.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}
