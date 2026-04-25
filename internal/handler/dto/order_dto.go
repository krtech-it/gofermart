package dto

import (
	"github.com/krtech-it/gofermart/internal/model"
	"time"
)

type OrderResponse struct {
	Number     string            `json:"number"`
	Status     model.OrderStatus `json:"status"`
	Accrual    *float64          `json:"accrual,omitempty"`
	UploadedAt time.Time         `json:"uploaded_at"`
}
