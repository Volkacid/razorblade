package main

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/server"
	"github.com/Volkacid/razorblade/internal/app/server/middlewares"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	servConf := config.GetServerConfig()
	db := storage.CreateStorage()
	handlers := server.NewHandlersSet(db, ctx)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Use(middlewares.GetUserID)
		router.Use(middlewares.GzipHandle)
		router.Get("/", handlers.MainPage)
		router.Get("/{key}", handlers.GetHandler)
		router.Get("/ping", handlers.PingDB)
		router.Get("/api/user/urls", handlers.UrlsAPIHandler)
		router.Post("/", handlers.PostHandler)
		router.Post("/api/shorten", handlers.PostAPIHandler)
		router.Post("/api/shorten/batch", handlers.BatchHandler)
		router.Delete("/api/user/urls", handlers.DeleteUserURLs)
	})
	log.Fatal(http.ListenAndServe(servConf.ServerAddress, router))
}
