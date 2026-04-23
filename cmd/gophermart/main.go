package main

import (
	"github.com/krtech-it/gofermart/internal/config"
	"github.com/krtech-it/gofermart/internal/delivery/http"
	"github.com/krtech-it/gofermart/internal/handler"
	"github.com/krtech-it/gofermart/internal/service"
	"github.com/krtech-it/gofermart/internal/storage"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := storage.NewPostgresStorage(cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrate("file://./migrations")
	if err != nil {
		log.Fatal(err)
	}
	services := service.NewServices(db, db, db, cfg)
	handlers := handler.NewHandler(services)
	router := http.NewRouter(handlers, cfg)
	err = router.Run(cfg.RunAddress)
	if err != nil {
		log.Fatal(err)
	}
}
