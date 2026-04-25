package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"go.uber.org/zap"
)

// mockWithdrawalStorage — ручной мок WithdrawalStorage с настраиваемыми функциями.
type mockWithdrawalStorage struct {
	getAllWithdrawalsByUserID func(ctx context.Context, userID uuid.UUID) ([]*model.Withdrawal, error)
	createWithdrawal          func(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error
	getBalance                func(ctx context.Context, userID uuid.UUID) (*model.Balance, error)
}

func (m *mockWithdrawalStorage) GetAllWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Withdrawal, error) {
	return m.getAllWithdrawalsByUserID(ctx, userID)
}

func (m *mockWithdrawalStorage) CreateWithdrawal(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error {
	return m.createWithdrawal(ctx, userID, orderNumber, sum)
}

func (m *mockWithdrawalStorage) GetBalance(ctx context.Context, userID uuid.UUID) (*model.Balance, error) {
	return m.getBalance(ctx, userID)
}

func newTestWithdrawalService(mock *mockWithdrawalStorage) WithdrawalServiceInterface {
	return NewWithdrawalService(mock, zap.NewNop())
}

// --- WithdrawalProcess ---

func TestWithdrawalProcess_InvalidLuhn(t *testing.T) {
	svc := newTestWithdrawalService(&mockWithdrawalStorage{})

	err := svc.WithdrawalProcess(context.Background(), uuid.New(), invalidLuhn, 100)

	if !errors.Is(err, ErrorInvalidOrderNumber) {
		t.Errorf("ожидалась ошибка ErrorInvalidOrderNumber, получена: %v", err)
	}
}

func TestWithdrawalProcess_InsufficientFunds(t *testing.T) {
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return &model.Balance{Current: 50, Withdrawn: 0}, nil
		},
	})

	err := svc.WithdrawalProcess(context.Background(), uuid.New(), validLuhn, 100)

	if !errors.Is(err, ErrorBalanceInsufficientFunds) {
		t.Errorf("ожидалась ошибка ErrorBalanceInsufficientFunds, получена: %v", err)
	}
}

func TestWithdrawalProcess_Success(t *testing.T) {
	var capturedSum float64
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return &model.Balance{Current: 200, Withdrawn: 0}, nil
		},
		createWithdrawal: func(_ context.Context, _ uuid.UUID, _ string, sum float64) error {
			capturedSum = sum
			return nil
		},
	})

	err := svc.WithdrawalProcess(context.Background(), uuid.New(), validLuhn, 100)

	if err != nil {
		t.Fatalf("ожидался nil, получена ошибка: %v", err)
	}
	if capturedSum != 100 {
		t.Errorf("ожидалась сумма 100, получена: %v", capturedSum)
	}
}

func TestWithdrawalProcess_StorageError_OnGetBalance(t *testing.T) {
	storageErr := errors.New("db error")
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return nil, storageErr
		},
	})

	err := svc.WithdrawalProcess(context.Background(), uuid.New(), validLuhn, 100)

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}

func TestWithdrawalProcess_StorageError_OnCreate(t *testing.T) {
	storageErr := errors.New("insert failed")
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return &model.Balance{Current: 200}, nil
		},
		createWithdrawal: func(_ context.Context, _ uuid.UUID, _ string, _ float64) error {
			return storageErr
		},
	})

	err := svc.WithdrawalProcess(context.Background(), uuid.New(), validLuhn, 100)

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}

// --- GetWithdrawals ---

func TestGetWithdrawals_ReturnsList(t *testing.T) {
	userID := uuid.New()
	expected := []*model.Withdrawal{
		{ID: uuid.New(), UserId: userID, Order: validLuhn, Sum: 50},
	}
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getAllWithdrawalsByUserID: func(_ context.Context, _ uuid.UUID) ([]*model.Withdrawal, error) {
			return expected, nil
		},
	})

	result, err := svc.GetWithdrawals(context.Background(), userID)

	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(result) != len(expected) {
		t.Errorf("ожидалось %d записей, получено %d", len(expected), len(result))
	}
}

func TestGetWithdrawals_StorageError(t *testing.T) {
	storageErr := errors.New("db error")
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getAllWithdrawalsByUserID: func(_ context.Context, _ uuid.UUID) ([]*model.Withdrawal, error) {
			return nil, storageErr
		},
	})

	_, err := svc.GetWithdrawals(context.Background(), uuid.New())

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}

// --- GetBalance ---

func TestGetBalance_ReturnsBalance(t *testing.T) {
	expected := &model.Balance{Current: 150, Withdrawn: 50}
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return expected, nil
		},
	})

	result, err := svc.GetBalance(context.Background(), uuid.New())

	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if result.Current != expected.Current || result.Withdrawn != expected.Withdrawn {
		t.Errorf("ожидался баланс %+v, получен %+v", expected, result)
	}
}

func TestGetBalance_StorageError(t *testing.T) {
	storageErr := errors.New("db error")
	svc := newTestWithdrawalService(&mockWithdrawalStorage{
		getBalance: func(_ context.Context, _ uuid.UUID) (*model.Balance, error) {
			return nil, storageErr
		},
	})

	_, err := svc.GetBalance(context.Background(), uuid.New())

	if !errors.Is(err, storageErr) {
		t.Errorf("ожидалась ошибка хранилища, получена: %v", err)
	}
}
