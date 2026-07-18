// @title Ecommerce API
// @version 1.0
// @description A simple ecommerce backend written in Go.
// @BasePath /api/v1
package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/Aneeshie/ecommerce/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Aneeshie/ecommerce/internal/config"
	"github.com/Aneeshie/ecommerce/internal/database"
	"github.com/Aneeshie/ecommerce/internal/identity/handler"
	"github.com/Aneeshie/ecommerce/internal/identity/service"
	"github.com/Aneeshie/ecommerce/internal/identity/token"
	md "github.com/Aneeshie/ecommerce/internal/middleware"
	orderHandle "github.com/Aneeshie/ecommerce/internal/order/handler"
	orderServices "github.com/Aneeshie/ecommerce/internal/order/service"
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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	handler.RegisterRoutes(r, identityHandler, authMiddleware)

	productService := productServices.NewService(store)

	productHandler := productHandle.NewHandler(productService)

	productHandle.RegisterRoutes(r, productHandler, authMiddleware)

	orderService := orderServices.NewService(store)

	orderHandler := orderHandle.NewHandler(orderService)

	orderHandle.RegisterRoutes(r, orderHandler, authMiddleware)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
