package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/service"
	"go.uber.org/zap"
)

// --- mock UserServiceInterface ---

type mockUserService struct {
	createUser func(ctx context.Context, login, password string) (string, error)
	login      func(ctx context.Context, login, password string) (string, error)
}

func (m *mockUserService) CreateUser(ctx context.Context, login, password string) (string, error) {
	return m.createUser(ctx, login, password)
}
func (m *mockUserService) Login(ctx context.Context, login, password string) (string, error) {
	return m.login(ctx, login, password)
}

// --- mock OrderServiceInterface ---

type mockOrderService struct {
	getAllOrdersByUser func(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
	uploadOrder        func(ctx context.Context, userID uuid.UUID, orderNumber string) error
	updateOrder        func(ctx context.Context, userID uuid.UUID, order *model.Order) error
}

func (m *mockOrderService) GetAllOrdersByUser(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	return m.getAllOrdersByUser(ctx, userID)
}
func (m *mockOrderService) UploadOrder(ctx context.Context, userID uuid.UUID, orderNumber string) error {
	return m.uploadOrder(ctx, userID, orderNumber)
}
func (m *mockOrderService) UpdateOrder(ctx context.Context, userID uuid.UUID, order *model.Order) error {
	return m.updateOrder(ctx, userID, order)
}

// --- mock WithdrawalServiceInterface ---

type mockWithdrawalService struct {
	withdrawalProcess func(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error
	getWithdrawals    func(ctx context.Context, userID uuid.UUID) ([]*model.Withdrawal, error)
	getBalance        func(ctx context.Context, userID uuid.UUID) (*model.Balance, error)
}

func (m *mockWithdrawalService) WithdrawalProcess(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error {
	return m.withdrawalProcess(ctx, userID, orderNumber, sum)
}
func (m *mockWithdrawalService) GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]*model.Withdrawal, error) {
	return m.getWithdrawals(ctx, userID)
}
func (m *mockWithdrawalService) GetBalance(ctx context.Context, userID uuid.UUID) (*model.Balance, error) {
	return m.getBalance(ctx, userID)
}

// newTestHandler создаёт Handler с мок-сервисами.
func newTestHandler(u service.UserServiceInterface, o service.OrderServiceInterface, w service.WithdrawalServiceInterface) *Handler {
	return &Handler{
		user:       u,
		order:      o,
		withdrawal: w,
		logger:     zap.NewNop(),
	}
}
