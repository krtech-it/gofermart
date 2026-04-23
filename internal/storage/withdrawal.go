package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
)

// WithdrawalStorage определяет операции хранилища для работы со списаниями баллов.
type WithdrawalStorage interface {
	// GetAllWithdrawalsByUserID возвращает все списания пользователя, отсортированные от новых к старым.
	GetAllWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Withdrawal, error)
	// CreateWithdrawal создаёт новую операцию списания баллов в счёт указанного заказа.
	CreateWithdrawal(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error
	// GetBalance возвращает текущий баланс и суммарное количество списанных баллов пользователя.
	GetBalance(ctx context.Context, userID uuid.UUID) (*model.Balance, error)
}

// GetAllWithdrawalsByUserID возвращает все списания пользователя, отсортированные от новых к старым.
func (p *PostgresStorage) GetAllWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Withdrawal, error) {
	withdrawalsDB := make([]*model.Withdrawal, 0)
	rows, err := p.db.QueryContext(ctx, "select id, user_id, order_number, sum, processed_at from withdrawals where user_id = $1 order by processed_at desc", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var withdrawal = &model.Withdrawal{}
		err := rows.Scan(&withdrawal.ID, &withdrawal.UserId, &withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawalsDB = append(withdrawalsDB, withdrawal)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return withdrawalsDB, nil
}

// CreateWithdrawal создаёт новую операцию списания баллов в счёт указанного заказа.
func (p *PostgresStorage) CreateWithdrawal(ctx context.Context, userID uuid.UUID, orderNumber string, sum float64) error {
	_, err := p.db.ExecContext(ctx, "insert into withdrawals (id, user_id, order_number, sum) values ($1, $2, $3, $4)", uuid.New(), userID, orderNumber, sum)
	if err != nil {
		return err
	}
	return nil
}

// GetBalance возвращает текущий баланс и суммарное количество списанных баллов пользователя.
func (p *PostgresStorage) GetBalance(ctx context.Context, userID uuid.UUID) (*model.Balance, error) {
	balance := &model.Balance{}

	row := p.db.QueryRowContext(ctx, "select "+
		"coalesce((select sum(o.accrual) from orders o where user_id = $1), 0)"+
		"-"+
		"coalesce((select sum(w.sum) from withdrawals w where user_id = $1), 0) "+
		"as current, COALESCE((select sum(w.sum) from withdrawals w where user_id = $1), 0) as withdrawn", userID)
	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
