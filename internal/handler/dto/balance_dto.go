package dto

import "time"

type BalanceDto struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawProcessRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type AllWithdrawResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
