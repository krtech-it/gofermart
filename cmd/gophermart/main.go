package main

import (
	"context"
	"github.com/krtech-it/gofermart/internal/accrual"
	"github.com/krtech-it/gofermart/internal/config"
	"github.com/krtech-it/gofermart/internal/delivery/http"
	"github.com/krtech-it/gofermart/internal/handler"
	"github.com/krtech-it/gofermart/internal/logger"
	"github.com/krtech-it/gofermart/internal/service"
	"github.com/krtech-it/gofermart/internal/storage"
	"github.com/krtech-it/gofermart/internal/worker"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
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
		customLogger.Error("failed to connect to database", zap.Error(err))
	}
	err = db.Migrate("file://./migrations")
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress, customLogger)
	w := worker.NewWorker(db, accrualClient)
	go w.Start(ctx)

	services := service.NewServices(db, db, db, cfg)
	handlers := handler.NewHandler(services)
	router := http.NewRouter(handlers, cfg)
	err = router.Run(cfg.RunAddress)
	if err != nil {
		log.Fatal(err)
	}
}
