package server

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/server/middlewares"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
)

// TestRequest For test purposes only
func TestRequest(query, method string, bodyReader io.Reader, db storage.Storage, userID string) *httptest.ResponseRecorder {
	handlers := NewHandlersSet(context.Background(), db)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
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
	request := httptest.NewRequest(method, query, bodyReader)
	writer := httptest.NewRecorder()
	http.SetCookie(writer, &http.Cookie{Name: "UserID", Value: userID})
	ctx := context.WithValue(request.Context(), config.UserID{}, userID)
	router.ServeHTTP(writer, request.WithContext(ctx))
	return writer
}
