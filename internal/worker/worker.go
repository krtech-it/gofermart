package worker

import (
	"context"
	"errors"
	"time"

	"github.com/krtech-it/gofermart/internal/accrual"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
	"go.uber.org/zap"
)

// Worker опрашивает хранилище заказов и обновляет их статусы через сервис начислений.
type Worker struct {
	storage storage.OrderStorage
	accrual *accrual.Client
	logger  *zap.Logger
}

// NewWorker создаёт новый Worker с переданными зависимостями.
func NewWorker(storage storage.OrderStorage, accrual *accrual.Client, logger *zap.Logger) *Worker {
	return &Worker{storage: storage, accrual: accrual, logger: logger}
}

func (w *Worker) Start(ctx context.Context) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			w.processOrders(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) processOrders(ctx context.Context) {
	orders, err := w.storage.GetAllOpenOrders(ctx)
	if err != nil {
		w.logger.Error("processOrders: ошибка получения заказов", zap.Error(err))
		return
	}
	for _, order := range orders {
		resp, err := w.accrual.GetOrderAccrual(ctx, order.Number)
		if err != nil {
			var rateLimitErr *accrual.RateLimitError
			if errors.As(err, &rateLimitErr) {
				w.logger.Warn("processOrders: превышен лимит запросов, ожидание", zap.Int("retry_after", rateLimitErr.RetryAfter))
				time.Sleep(time.Duration(rateLimitErr.RetryAfter) * time.Second)
				return
			}
			w.logger.Error("processOrders: ошибка получения начисления", zap.String("order", order.Number), zap.Error(err))
			continue
		}
		if resp == nil {
			continue
		}
		order.Accrual = &resp.Accrual
		switch resp.Status {
		case "REGISTERED":
			order.Status = model.OrderStatusProcessing
		default:
			order.Status = model.OrderStatus(resp.Status)
		}
		err = w.storage.UpdateOrder(ctx, order)
		if err != nil {
			w.logger.Error("processOrders: ошибка обновления заказа", zap.String("order", order.Number), zap.Error(err))
			continue
		}
		w.logger.Debug("processOrders: заказ обновлён", zap.String("order", order.Number), zap.String("status", string(order.Status)))
	}
}
