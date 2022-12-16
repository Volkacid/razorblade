package server

import (
	"github.com/Volkacid/razorblade/internal/app/server/middlewares"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http/httptest"
)

func TestRequest(query, method string, bodyReader io.Reader, db storage.Storage) *httptest.ResponseRecorder {
	handlers := NewHandlersSet(db)

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
	})
	request := httptest.NewRequest(method, query, bodyReader)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	return writer
}
