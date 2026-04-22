package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/krtech-it/gofermart/internal/handler/dto"
	"github.com/krtech-it/gofermart/internal/service"
	"net/http"
)

func (h *Handler) Login(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(400)
		return
	}
	tokenStr, err := h.user.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrorInvalidLoginPassword) {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.SetCookie("Authorization", tokenStr, 3600, "/", "", false, true)
	c.Status(http.StatusOK)

}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	tokenStr, err := h.user.CreateUser(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrorLoginAlreadyExists) {
			c.JSON(409, gin.H{"error": "login already exists"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.SetCookie("Authorization", tokenStr, 3600, "/", "", false, true)
	c.Status(http.StatusOK)
}
