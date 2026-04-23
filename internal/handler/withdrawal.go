package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/handler/dto"
	"github.com/krtech-it/gofermart/internal/service"
	"net/http"
)

func (h *Handler) GetBalance(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	balance, err := h.withdrawal.GetBalance(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	balanceDTO := dto.BalanceDto{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}
	c.JSON(200, balanceDTO)
}

func (h *Handler) WithdrawProcess(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var req dto.WithdrawProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err := h.withdrawal.WithdrawalProcess(c.Request.Context(), userUUID, req.Order, req.Sum)
	if err != nil {
		if errors.Is(err, service.ErrorInvalidOrderNumber) {
			c.JSON(422, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrorBalanceInsufficientFunds) {
			c.JSON(402, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.Status(200)
}

func (h *Handler) GetWithdrawals(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	withdraws, err := h.withdrawal.GetWithdrawals(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if len(withdraws) == 0 {
		c.Status(204)
		return
	}
	withdrawsDTO := make([]*dto.AllWithdrawResponse, len(withdraws))
	for i, withdraw := range withdraws {
		withdrawsDTO[i] = &dto.AllWithdrawResponse{
			Order:       withdraw.Order,
			Sum:         withdraw.Sum,
			ProcessedAt: withdraw.ProcessedAt,
		}
	}
	c.JSON(200, withdrawsDTO)
}
