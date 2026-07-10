package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/config"
	"github.com/Aneeshie/ecommerce/internal/database"
	md "github.com/Aneeshie/ecommerce/internal/middleware"
	"github.com/Aneeshie/ecommerce/internal/identity/handler"
	productHandle "github.com/Aneeshie/ecommerce/internal/product/handler"
	"github.com/Aneeshie/ecommerce/internal/identity/repository"
	productRepo "github.com/Aneeshie/ecommerce/internal/product/repository"
	"github.com/Aneeshie/ecommerce/internal/identity/service"
	productServices "github.com/Aneeshie/ecommerce/internal/product/service"
	"github.com/Aneeshie/ecommerce/internal/identity/token"
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

	identityManager := token.NewManager(cfg.JwtSecret)

	identityService := service.NewService(identityRepository, identityManager)

	identityHandler := handler.NewHandler(identityService)

	authMiddleware := md.NewAuthMiddleware(identityManager)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler.RegisterRoutes(r, identityHandler, authMiddleware)

	productRepository := productRepo.NewRepository(pool)

	productService := productServices.NewService(productRepository)

	productHandler := productHandle.NewHandler(productService)

	productHandle.RegisterRoutes(r, productHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
