// Пакет config загружает и валидирует конфигурацию приложения из флагов CLI
// и переменных окружения. Переменные окружения имеют приоритет над флагами.
package config

import (
	"net"
	"os"
)

const JwtKey = "secret_jwt_key"

// Config хранит конфигурацию приложения.
type Config struct {
	// RunAddress — адрес host:port, на котором слушает HTTP-сервер (env: RUN_ADDRESS, флаг: -a).
	RunAddress string
	// DatabaseURI — строка подключения к PostgreSQL (env: DATABASE_URI, флаг: -d).
	DatabaseURI string
	// AccrualSystemAddress — базовый URL внешнего сервиса начислений (env: ACCRUAL_SYSTEM_ADDRESS, флаг: -r).
	AccrualSystemAddress string
	JWTSecret            string
	LogLevel             string
}

// Load собирает Config из флагов CLI с переопределением через переменные окружения
// и проверяет, что RunAddress является корректной парой host:port.
// Возвращает ошибку, если формат адреса некорректен.
func Load() (Config, error) {
	configFlag := ParseFlag()
	runAddress := getEnv("RUN_ADDRESS", configFlag.RunAddress)
	databaseURI := getEnv("DATABASE_URI", configFlag.DatabaseURI)
	accrualSystemAddress := getEnv("ACCRUAL_SYSTEM_ADDRESS", configFlag.AccrualSystemAddress)
	if err := checkHostPortAddr(runAddress); err != nil {
		return Config{}, err
	}
	return Config{
		RunAddress:           runAddress,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrualSystemAddress,
		JWTSecret:            getEnv("JWT_SECRET", JwtKey),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func checkHostPortAddr(addr string) error {
	_, _, err := net.SplitHostPort(addr)
	return err
}
