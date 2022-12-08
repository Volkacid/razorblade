package main

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/server"
	"github.com/Volkacid/razorblade/internal/app/server/middlewares"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

func main() {
	servConf := config.GetServerConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pgPool, _ := pgxpool.New(ctx, servConf.DBAddress)
	db := storage.CreateStorage(pgPool)
	defer pgPool.Close()

	service.SetCreatorSeed(time.Now().Unix())

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Use(middlewares.GetUserID)
		router.Use(middlewares.GzipHandle)
		router.Get("/", server.MainPage)
		router.Get("/{key}", server.GetHandler(db))
		router.Get("/ping", server.PingDB())
		router.Get("/api/user/urls", server.UrlsAPIHandler(db))
		router.Post("/", server.PostHandler(db))
		router.Post("/api/shorten", server.PostAPIHandler(db))
		router.Post("/api/shorten/batch", server.BatchHandler(db))
	})
	log.Fatal(http.ListenAndServe(servConf.ServerAddress, router))
}
