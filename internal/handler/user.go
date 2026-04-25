package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krtech-it/gofermart/internal/handler/dto"
	"github.com/krtech-it/gofermart/internal/service"
	"go.uber.org/zap"
)

func (h *Handler) Login(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Debug("Login: невалидное тело запроса", zap.Error(err))
		c.AbortWithStatus(400)
		return
	}
	tokenStr, err := h.user.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrorInvalidLoginPassword) {
			h.logger.Debug("Login: неверный логин или пароль", zap.String("login", req.Login))
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Login: внутренняя ошибка", zap.Error(err))
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	h.logger.Debug("Login: успешно", zap.String("login", req.Login))
	c.SetCookie("Authorization", tokenStr, 3600, "/", "", false, true)
	c.Status(http.StatusOK)
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Debug("Register: невалидное тело запроса", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}
	tokenStr, err := h.user.CreateUser(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrorLoginAlreadyExists) {
			h.logger.Debug("Register: логин уже занят", zap.String("login", req.Login))
			c.JSON(409, gin.H{"error": "login already exists"})
			return
		}
		h.logger.Error("Register: внутренняя ошибка", zap.Error(err))
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	h.logger.Debug("Register: пользователь создан", zap.String("login", req.Login))
	c.SetCookie("Authorization", tokenStr, 3600, "/", "", false, true)
	c.Status(http.StatusOK)
}
