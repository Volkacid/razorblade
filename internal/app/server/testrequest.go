package server

import (
	"github.com/Volkacid/razortest/internal/app/server/middlewares"
	"github.com/Volkacid/razortest/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http/httptest"
)

func TestRequest(query, method string, bodyReader io.Reader, db *storage.Storage) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Use(middlewares.GetUserID)
		router.Use(middlewares.GzipHandle)
		router.Get("/", MainPage)
		router.Get("/{key}", GetHandler(db))
		router.Get("/api/user/urls", APIUrlsHandler(db))
		router.Post("/", PostHandler(db))
		router.Post("/api/shorten", APIPostHandler(db))
	})
	request := httptest.NewRequest(method, query, bodyReader)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	return writer
}
