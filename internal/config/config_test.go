package config

import (
	"testing"
)

// --- getEnv ---

func TestGetEnv_ReturnsEnvWhenSet(t *testing.T) {
	t.Setenv("TEST_KEY_GOPHERMART", "from_env")

	result := getEnv("TEST_KEY_GOPHERMART", "fallback")

	if result != "from_env" {
		t.Errorf("ожидалось %q, получено %q", "from_env", result)
	}
}

func TestGetEnv_ReturnsFallbackWhenNotSet(t *testing.T) {
	result := getEnv("GOPHERMART_UNSET_KEY_XYZ", "fallback")

	if result != "fallback" {
		t.Errorf("ожидалось %q, получено %q", "fallback", result)
	}
}

// --- checkHostPortAddr ---

func TestCheckHostPortAddr_Valid(t *testing.T) {
	if err := checkHostPortAddr("localhost:8080"); err != nil {
		t.Errorf("неожиданная ошибка для корректного адреса: %v", err)
	}
}

func TestCheckHostPortAddr_Invalid(t *testing.T) {
	if err := checkHostPortAddr("not-an-address"); err == nil {
		t.Error("ожидалась ошибка для некорректного адреса")
	}
}

func TestCheckHostPortAddr_NoPort(t *testing.T) {
	if err := checkHostPortAddr("localhost"); err == nil {
		t.Error("ожидалась ошибка для адреса без порта")
	}
}

func TestCheckHostPortAddr_EmptyPort(t *testing.T) {
	if err := checkHostPortAddr(":"); err != nil {
		t.Errorf("неожиданная ошибка для адреса с пустым портом: %v", err)
	}
}

// --- Load ---
// Load вызывается один раз — повторный вызов паникует из-за повторной регистрации флагов.

func TestLoad_ReturnsValidConfig(t *testing.T) {
	t.Setenv("RUN_ADDRESS", "127.0.0.1:8888")
	t.Setenv("DATABASE_URI", "postgres://localhost/test")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "localhost:9000")
	t.Setenv("JWT_SECRET", "supersecret")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if cfg.RunAddress != "127.0.0.1:8888" {
		t.Errorf("RunAddress: ожидалось %q, получено %q", "127.0.0.1:8888", cfg.RunAddress)
	}
	if cfg.DatabaseURI != "postgres://localhost/test" {
		t.Errorf("DatabaseURI: ожидалось %q, получено %q", "postgres://localhost/test", cfg.DatabaseURI)
	}
	if cfg.AccrualSystemAddress != "http://localhost:9000" {
		t.Errorf("AccrualSystemAddress: ожидалось %q, получено %q", "http://localhost:9000", cfg.AccrualSystemAddress)
	}
	if cfg.JWTSecret != "supersecret" {
		t.Errorf("JWTSecret: ожидалось %q, получено %q", "supersecret", cfg.JWTSecret)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel: ожидалось %q, получено %q", "debug", cfg.LogLevel)
	}
}
