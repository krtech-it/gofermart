package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/krtech-it/gofermart/internal/config"
	"time"
)

func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Request.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		tokenString := token.Value
		userID, err := checkAuthToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}

func checkAuthToken(tokenString string, secretKey string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid token")
	}
	userUUID, err := uuid.Parse(userID)
	return userUUID, nil
}

func GenerateToken(userID uuid.UUID, secretKey string) (string, error) {
	data := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(secretKey))
}
