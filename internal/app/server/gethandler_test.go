package server

import (
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		want       string
		statusCode int
	}{
		{
			name:       "Trying to get original link",
			query:      "/testlink",
			want:       "https://yandex.com",
			statusCode: 307,
		},
		{
			name:       "Trying to get link that is not in map",
			query:      "/someotherlink",
			want:       "",
			statusCode: 404,
		},
	}
	db := storage.CreateStorage(false)
	db.SaveValue("testlink", "https://yandex.com")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := TestRequest(tt.query, http.MethodGet, nil, db)
			response := recorder.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.want, response.Header.Get("Location"))
			assert.Equal(t, tt.statusCode, response.StatusCode)
		})
	}
}
