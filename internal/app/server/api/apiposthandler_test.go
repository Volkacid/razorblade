package api

import (
	"bytes"
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/server"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPostAPIHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      URL
		statusCode int
	}{
		{
			name:       "Trying to short a correct link",
			query:      URL{URL: "https://github.com/Volkacid/razorblade"},
			statusCode: 201,
		},
		{
			name:       "Trying to short incorrect link",
			query:      URL{URL: "http://somestring"},
			statusCode: 400,
		},
	}
	db := storage.CreateTestStorage()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshalledQuery, _ := json.Marshal(tt.query)
			reader := bytes.NewReader(marshalledQuery)
			recorder := server.TestRequest("/api/shorten", http.MethodPost, reader, db)
			response := recorder.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
			assert.NotNil(t, response.Body)
		})
	}
}
