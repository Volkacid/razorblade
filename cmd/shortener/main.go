package main

import (
	"github.com/Volkacid/razortest/internal/app/config"
	"github.com/Volkacid/razortest/internal/app/server"
	"github.com/Volkacid/razortest/internal/app/server/middlewares"
	"github.com/Volkacid/razortest/internal/app/service"
	"github.com/Volkacid/razortest/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

func main() {
	servConf := config.GetServerConfig()
	isStorageFileExist := servConf.StorageFile != ""
	db := storage.CreateStorage(isStorageFileExist)
	service.SetCreatorSeed(time.Now().Unix())

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Use(middlewares.GetUserID)
		router.Use(middlewares.GzipHandle)
		router.Get("/", server.MainPage)
		router.Get("/{key}", server.GetHandler(db))
		router.Get("/api/user/urls", server.APIUrlsHandler(db))
		router.Post("/", server.PostHandler(db))
		router.Post("/api/shorten", server.APIPostHandler(db))
	})
	log.Fatal(http.ListenAndServe(servConf.ServerAddress, router))
}
