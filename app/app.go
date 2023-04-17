package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"

	"refactoring/config"
	"refactoring/internal/repository/storage"
	"refactoring/internal/service"
)

func Run() error {
	// инициализируем конфиги
	config.Configure()

	// инициализируем роутер
	root := chi.NewRouter()

	// используем middlewares
	root.Use(middleware.RequestID)
	root.Use(middleware.RealIP)
	root.Use(middleware.Logger)
	root.Use(middleware.Recoverer)
	root.Use(middleware.Timeout(viper.GetDuration(config.CtxTimeout)))

	// инициализируем storage
	store, err := storage.NewJsonStorage(viper.GetString(config.Store))
	if err != nil {
		return fmt.Errorf("error init storage: %w", err)
	}

	// инициализируем service
	srv := service.NewService(store, root)

	// инициализируем routes
	srv.InitRoutes()

	log.Println("start server on:", viper.GetString(config.ServerAddr))
	// открываем сервер
	return http.ListenAndServe(viper.GetString(config.ServerAddr), root)
}
