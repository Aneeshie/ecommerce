package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/identity/handler"
	"github.com/Aneeshie/ecommerce/internal/identity/repository"
	"github.com/Aneeshie/ecommerce/internal/identity/service"
	"github.com/Aneeshie/ecommerce/internal/config"
	"github.com/Aneeshie/ecommerce/internal/database"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	identityRepository := repository.NewRepository(pool)

	identityService := service.NewService(identityRepository)

	identityHandler := handler.NewHandler(identityService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler.RegisterRoutes(r, identityHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
