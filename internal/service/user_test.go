package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// mockUserStorage — ручной мок UserStorage с настраиваемыми функциями.
type mockUserStorage struct {
	createUser     func(ctx context.Context, user *model.User) error
	getUserByLogin func(ctx context.Context, login string) (*model.User, error)
}

func (m *mockUserStorage) CreateUser(ctx context.Context, user *model.User) error {
	return m.createUser(ctx, user)
}

func (m *mockUserStorage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	return m.getUserByLogin(ctx, login)
}

const testJWTSecret = "test-secret"

func newTestUserService(mock *mockUserStorage) UserServiceInterface {
	return NewUserService(mock, testJWTSecret, zap.NewNop())
}

// --- CreateUser ---

func TestCreateUser_Success(t *testing.T) {
	svc := newTestUserService(&mockUserStorage{
		getUserByLogin: func(_ context.Context, _ string) (*model.User, error) {
			return nil, storage.ErrNotFound
		},
		createUser: func(_ context.Context, _ *model.User) error {
			return nil
		},
	})

	token, err := svc.CreateUser(context.Background(), "user1", "pass123")

	if err != nil {
		t.Fatalf("ожидался nil, получена ошибка: %v", err)
	}
	if token == "" {
		t.Error("ожидался непустой токен")
	}
}

func TestCreateUser_LoginAlreadyExists(t *testing.T) {
	svc := newTestUserService(&mockUserStorage{
		getUserByLogin: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{ID: uuid.New(), Login: "user1"}, nil
		},
	})

	_, err := svc.CreateUser(context.Background(), "user1", "pass123")

	if !errors.Is(err, ErrorLoginAlreadyExists) {
		t.Errorf("ожидалась ошибка ErrorLoginAlreadyExists, получена: %v", err)
	}
}

func TestCreateUser_StorageError_OnGet(t *testing.T) {
	storageErr := errors.New("db error")
	svc := newTestUserService(&mockUserStorage{
		getUserByLogin: func(_ context.Context, _ string) (*model.User, error) {
			return nil, storageErr
		},
	})

	_, err := svc.CreateUser(context.Background(), "user1", "pass123")

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}

func TestCreateUser_StorageError_OnCreate(t *testing.T) {
	storageErr := errors.New("insert failed")
	svc := newTestUserService(&mockUserStorage{
		getUserByLogin: func(_ context.Context, _ string) (*model.User, error) {
			return nil, storage.ErrNotFound
		},
		createUser: func(_ context.Context, _ *model.User) error {
			return storageErr
		},
	})

	_, err := svc.CreateUser(context.Background(), "user1", "pass123")

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}

// --- Login ---

func TestLogin_Success(t *testing.T) {
	password := "pass123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("не удалось создать хеш пароля: %v", err)
	}
	svc := newTestUserService(&mockUserStorage{
		getUserByLogin: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{ID: uuid.New(), Login: "user1", PasswordHash: string(hash)}, nil
		},
	})

	token, err := svc.Login(context.Background(), "user1", password)

	if err != nil {
		t.Fatalf("ожидался nil, получена ошибка: %v", err)
	}
	if token == "" {
		t.Error("ожидался непустой токен")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc := newTestUserService(&mockUserStorage{
		getUserByLogin: func(_ context.Context, _ string) (*model.User, error) {
			return nil, storage.ErrNotFound
		},
	})

	_, err := svc.Login(context.Background(), "unknown", "pass123")

	if !errors.Is(err, ErrorInvalidLoginPassword) {
		t.Errorf("ожидалась ошибка ErrorInvalidLoginPassword, получена: %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("не удалось создать хеш пароля: %v", err)
	}
	svc := newTestUserService(&mockUserStorage{
		getUserByLogin: func(_ context.Context, _ string) (*model.User, error) {
			return &model.User{ID: uuid.New(), Login: "user1", PasswordHash: string(hash)}, nil
		},
	})

	_, err = svc.Login(context.Background(), "user1", "wrongpassword")

	if !errors.Is(err, ErrorInvalidLoginPassword) {
		t.Errorf("ожидалась ошибка ErrorInvalidLoginPassword, получена: %v", err)
	}
}
