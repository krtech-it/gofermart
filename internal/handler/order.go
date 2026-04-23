package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/handler/dto"
	"github.com/krtech-it/gofermart/internal/service"
	"go.uber.org/zap"
)

func (h *Handler) UploadOrder(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Debug("UploadOrder: ошибка чтения тела запроса", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	orderNumber := strings.TrimSpace(string(body))
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
	err = h.order.UploadOrder(c.Request.Context(), userUUID, orderNumber)
	if err != nil {
		if errors.Is(err, service.ErrOrderAlreadyByThisUser) {
			h.logger.Debug("UploadOrder: заказ уже загружен этим пользователем", zap.String("order", orderNumber))
			c.Status(200)
			return
		}
		if errors.Is(err, service.ErrOrderAlreadyByOtherUser) {
			h.logger.Debug("UploadOrder: заказ принадлежит другому пользователю", zap.String("order", orderNumber))
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrorInvalidOrderNumber) {
			h.logger.Debug("UploadOrder: невалидный номер заказа", zap.String("order", orderNumber))
			c.JSON(422, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("UploadOrder: внутренняя ошибка", zap.Error(err))
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	h.logger.Debug("UploadOrder: заказ принят", zap.String("order", orderNumber))
	c.Status(202)
}

func (h *Handler) GetOrders(c *gin.Context) {
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
	orders, err := h.order.GetAllOrdersByUser(c.Request.Context(), userUUID)
	if err != nil {
		h.logger.Error("GetOrders: внутренняя ошибка", zap.Error(err))
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if len(orders) == 0 {
		c.Status(204)
		return
	}
	ordersDTO := make([]*dto.OrderResponse, len(orders))
	for i, order := range orders {
		ordersDTO[i] = &dto.OrderResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		}
	}
	c.JSON(200, ordersDTO)
}
