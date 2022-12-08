package main

import (
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/server"
	"github.com/Volkacid/razorblade/internal/app/server/api"
	"github.com/Volkacid/razorblade/internal/app/server/middlewares"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

func main() {
	db := storage.CreateStorage()

	service.SetCreatorSeed(time.Now().Unix())

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Use(middlewares.GetUserID)
		router.Use(middlewares.GzipHandle)
		router.Get("/", server.MainPage)
		router.Get("/{key}", server.GetHandler(db))
		router.Get("/ping", server.PingDB())
		router.Get("/api/user/urls", api.UrlsAPIHandler(db))
		router.Post("/", server.PostHandler(db))
		router.Post("/api/shorten", api.PostAPIHandler(db))
		router.Post("/api/shorten/batch", api.BatchHandler(db))
	})
	log.Fatal(http.ListenAndServe(config.GetServerConfig().ServerAddress, router))
}
