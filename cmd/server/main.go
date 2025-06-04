package main

import (
	"SnickersShopPet1.0/internal/config"
	"SnickersShopPet1.0/internal/handler"
	"SnickersShopPet1.0/internal/logger"
	"SnickersShopPet1.0/internal/repository"
	"SnickersShopPet1.0/pkg/database"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("dev.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	zaplog, err := logger.LoadLogger(cfg)
	if err != nil {
		log.Fatalf("Error loading logger: %v", err)
	}

	if err = database.LoadDatabase(cfg, "./migrations"); err != nil {
		log.Fatalf("Error loading database: %v", err)
	}

	r := chi.NewRouter()

	snickersRepository := repository.NewSnickersRepository(database.ReturnDB(), zaplog)
	snickersHandler := handler.NewSnickersHandler(snickersRepository)

	r.Post("/api/new_snickers", snickersHandler.AddSnickersPOST)

	zaplog.Info("Server started")

	log.Fatal(http.ListenAndServe(cfg.Server.Port, r))

}
