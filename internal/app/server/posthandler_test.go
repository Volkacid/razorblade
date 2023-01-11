package server

import (
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestPostHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		statusCode int
	}{
		{
			name:       "Trying to short a correct link",
			query:      "https://github.com/Volkacid/razorblade",
			statusCode: 201,
		},
		{
			name:       "Trying to short incorrect link",
			query:      "http://somestring",
			statusCode: 400,
		},
	}
	db := storage.CreateTestStorage()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.query)
			recorder := TestRequest("/", http.MethodPost, bodyReader, db, "someid")
			response := recorder.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
			assert.NotNil(t, response.Body)
		})
	}
}
