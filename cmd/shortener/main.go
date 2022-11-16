package main

import (
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/server"
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

	servConf := config.GetServerConfig()

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", server.MainPage)
		router.Get("/{key}", server.GetHandler(db))
		router.Post("/", server.PostHandler(db))
		router.Post("/api/shorten", server.APIPostHandler(db))
	})
	log.Fatal(http.ListenAndServe(servConf.ServerAddress, router))
}
