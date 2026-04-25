package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/krtech-it/gofermart/internal/accrual"
	"github.com/krtech-it/gofermart/internal/config"
	delivery "github.com/krtech-it/gofermart/internal/delivery/http"
	"github.com/krtech-it/gofermart/internal/handler"
	"github.com/krtech-it/gofermart/internal/logger"
	"github.com/krtech-it/gofermart/internal/service"
	"github.com/krtech-it/gofermart/internal/storage"
	"github.com/krtech-it/gofermart/internal/worker"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	customLogger, err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	db, err := storage.NewPostgresStorage(cfg.DatabaseURI)
	if err != nil {
		customLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	err = db.Migrate("file://./migrations")
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress, customLogger)
	w := worker.NewWorker(db, accrualClient, customLogger)
	go w.Start(ctx)

	services := service.NewServices(db, db, db, cfg, customLogger)
	handlers := handler.NewHandler(services, customLogger)
	router := delivery.NewRouter(handlers, cfg)

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			customLogger.Fatal("server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		customLogger.Error("shutdown error", zap.Error(err))
	}
}
