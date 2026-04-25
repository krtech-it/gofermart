package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/model"
)

// OrderStorage определяет операции хранилища для работы с заказами.
type OrderStorage interface {
	// GetAllOrdersByUserID возвращает все заказы пользователя, отсортированные от новых к старым.
	GetAllOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
	// GetOrderByNumber возвращает заказ по его номеру.
	// Возвращает ошибку, если заказ не найден.
	GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error)
	// CreateOrder сохраняет новый заказ и возвращает его с заполненными полями.
	CreateOrder(ctx context.Context, order *model.Order) error
	// UpdateOrder обновляет статус и начисление заказа и возвращает обновлённую запись.
	UpdateOrder(ctx context.Context, order *model.Order) error
	GetAllOpenOrders(ctx context.Context) ([]*model.Order, error)
}

// GetAllOrdersByUserID возвращает все заказы пользователя, отсортированные от новых к старым.
func (p *PostgresStorage) GetAllOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	ordersDB := make([]*model.Order, 0)
	rows, err := p.db.QueryContext(ctx, "select number, user_id, status, accrual, uploaded_at from orders where user_id = $1 order by uploaded_at desc", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order = &model.Order{}
		var accrual sql.NullFloat64
		err := rows.Scan(&order.Number, &order.UserId, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		if accrual.Valid {
			order.Accrual = &accrual.Float64
		}
		ordersDB = append(ordersDB, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ordersDB, nil
}

// GetOrderByNumber возвращает заказ по его номеру.
// Возвращает ошибку, если заказ не найден.
func (p *PostgresStorage) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	orderDB := &model.Order{}
	var accrual sql.NullFloat64
	row := p.db.QueryRowContext(ctx, "select number, user_id, status, accrual, uploaded_at from orders where number = $1", orderNumber)
	err := row.Scan(&orderDB.Number, &orderDB.UserId, &orderDB.Status, &accrual, &orderDB.UploadedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if accrual.Valid {
		orderDB.Accrual = &accrual.Float64
	}
	return orderDB, nil
}

// CreateOrder сохраняет новый заказ и возвращает его с заполненными полями.
func (p *PostgresStorage) CreateOrder(ctx context.Context, order *model.Order) error {
	_, err := p.db.ExecContext(ctx, "insert into orders (number, user_id, status) values ($1, $2, $3)", order.Number, order.UserId, order.Status)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrder обновляет статус и начисление заказа и возвращает обновлённую запись.
func (p *PostgresStorage) UpdateOrder(ctx context.Context, order *model.Order) error {
	_, err := p.db.ExecContext(ctx, "update orders set status = $1, accrual = $2 where number = $3", order.Status, order.Accrual, order.Number)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) GetAllOpenOrders(ctx context.Context) ([]*model.Order, error) {
	ordersDB := make([]*model.Order, 0)
	rows, err := p.db.QueryContext(ctx, "select number, user_id, status from orders where status in ('NEW', 'PROCESSING')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order = &model.Order{}
		err := rows.Scan(&order.Number, &order.UserId, &order.Status)
		if err != nil {
			return nil, err
		}
		ordersDB = append(ordersDB, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ordersDB, nil
}
