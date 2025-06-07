package main

import (
	"SnickersShopPet1.0/internal/config"
	"SnickersShopPet1.0/internal/handler"
	"SnickersShopPet1.0/internal/logger"
	"SnickersShopPet1.0/internal/repository"
	"SnickersShopPet1.0/pkg/database"
	middleware2 "SnickersShopPet1.0/pkg/middleware"
	redis2 "SnickersShopPet1.0/pkg/redis"
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

	redisClient := redis2.NewRedisClient(cfg)

	snickersRepository := repository.NewSnickersRepository(database.ReturnDB(), zaplog, redisClient)
	userRepository := repository.NewUserRepository(database.ReturnDB(), zaplog)

	snickersHandler := handler.NewSnickersHandler(snickersRepository)
	userHandler := handler.NewUserHandler(userRepository)

	r.Post("/api/new_user", userHandler.NewUserPOST)
	r.Post("/api/login", userHandler.LogInPOST)
	r.Get("/api/all_snickers", snickersHandler.AllSnickersGET)
	r.Get("/api/get_snickers/{id}", snickersHandler.SnickersByIDGET)
	r.Get("/api/get_by_cost", snickersHandler.SnickersByCostGET)

	r.Group(func(r chi.Router) {
		r.Use(middleware2.AuthMiddleware)

		r.Post("/api/new_snickers", snickersHandler.AddSnickersPOST)
	})

	zaplog.Info("Server started")

	log.Fatal(http.ListenAndServe(cfg.Server.Port, r))

}

//TODO: refine repositories(redis), add paginated name search methods, add password encryption, cover the code with tests
