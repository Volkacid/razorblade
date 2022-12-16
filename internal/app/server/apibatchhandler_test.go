package server

import (
	"bytes"
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBatchHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      []BatchQuery
		statusCode int
	}{
		{
			name:       "Trying to short a correct URLs",
			query:      []BatchQuery{{CorrelationID: "id0", OriginalURL: "https://google.com"}, {CorrelationID: "id1", OriginalURL: "https://ya.ru"}},
			statusCode: 201,
		},
		{
			name:       "Trying to short an incorrect URL",
			query:      []BatchQuery{{CorrelationID: "id0", OriginalURL: "https://google"}, {CorrelationID: "id1", OriginalURL: "https://ya.ru"}},
			statusCode: 400,
		},
	}
	db := storage.CreateTestStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshalledQuery, _ := json.Marshal(tt.query)
			reader := bytes.NewReader(marshalledQuery)
			recorder := TestRequest("/api/shorten/batch", http.MethodPost, reader, db)
			response := recorder.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
			assert.NotNil(t, response.Body)
		})
	}
}
