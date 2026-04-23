package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/krtech-it/gofermart/internal/accrual"
	"github.com/krtech-it/gofermart/internal/model"
	"github.com/krtech-it/gofermart/internal/storage"
)

type Worker struct {
	storage storage.OrderStorage
	accrual *accrual.Client
}

func NewWorker(storage storage.OrderStorage, accrual *accrual.Client) *Worker {
	return &Worker{storage: storage, accrual: accrual}
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
		log.Printf("error getting all open orders: %v", err)
		return
	}
	for _, order := range orders {
		resp, err := w.accrual.GetOrderAccrual(ctx, order.Number)
		if err != nil {
			var rateLimitErr *accrual.RateLimitError
			if errors.As(err, &rateLimitErr) {
				log.Printf("rate limit error: %v", err)
				time.Sleep(time.Duration(rateLimitErr.RetryAfter) * time.Second)
				return
			}
			log.Printf("error getting order accrual: %v", err)
			continue
		}
		if resp == nil {
			continue
		}
		order.Accrual = &resp.Accrual
		order.Status = model.OrderStatus(resp.Status)
		err = w.storage.UpdateOrder(ctx, order)
		if err != nil {
			log.Printf("error updating order: %v", err)
			continue
		}
	}

}
