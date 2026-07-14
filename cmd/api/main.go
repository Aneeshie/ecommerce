package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Aneeshie/ecommerce/internal/config"
	"github.com/Aneeshie/ecommerce/internal/database"
	"github.com/Aneeshie/ecommerce/internal/identity/handler"
	"github.com/Aneeshie/ecommerce/internal/identity/service"
	"github.com/Aneeshie/ecommerce/internal/identity/token"
	md "github.com/Aneeshie/ecommerce/internal/middleware"
	productHandle "github.com/Aneeshie/ecommerce/internal/product/handler"
	productServices "github.com/Aneeshie/ecommerce/internal/product/service"
	"github.com/Aneeshie/ecommerce/internal/store"
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

	store := store.NewStore(pool)

	identityManager := token.NewManager(cfg.JwtSecret)

	identityService := service.NewService(store, identityManager)

	identityHandler := handler.NewHandler(identityService)

	authMiddleware := md.NewAuthMiddleware(identityManager)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler.RegisterRoutes(r, identityHandler, authMiddleware)

	productService := productServices.NewService(store)

	productHandler := productHandle.NewHandler(productService)

	productHandle.RegisterRoutes(r, productHandler, authMiddleware)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
