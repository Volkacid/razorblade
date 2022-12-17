package server

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUrlsAPIHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		key        string
		UserID     string
		header     string
		statusCode int
	}{
		{
			name:       "Trying to get values with correct UserID",
			query:      "https://ya.ru",
			key:        "yandex",
			UserID:     "correctid",
			header:     "application/json",
			statusCode: http.StatusOK,
		},
		{
			name:       "Trying to get values with wrong UserID",
			query:      "http://google.com",
			key:        "google",
			UserID:     "wrongid",
			header:     "text/plain; charset=utf-8",
			statusCode: http.StatusNoContent,
		},
	}
	db := storage.CreateTestStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.SaveValue(context.Background(), tt.key, tt.query, "correctid")
			request := TestRequest("/api/user/urls", http.MethodGet, nil, db, tt.UserID)
			response := request.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
			assert.Equal(t, tt.header, response.Header.Get("Content-Type"))
		})
	}
}
