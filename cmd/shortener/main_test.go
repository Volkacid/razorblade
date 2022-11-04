package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(query, method string, bodyReader io.Reader, db map[string]string) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/{key}", GetHandler(db))
		router.Post("/", PostHandler(db))
	})
	request := httptest.NewRequest(method, query, bodyReader)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	return writer
}

func TestURLHandlerPOST(t *testing.T) {

	tests := []struct {
		name       string
		query      string
		statusCode int
		db         map[string]string
	}{
		{
			name:       "Trying to short link",
			query:      "https://github.com/Volkacid/razorblade",
			statusCode: 201,
			db:         make(map[string]string),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.query)
			recorder := testRequest("/", http.MethodPost, bodyReader, tt.db)
			response := recorder.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
			assert.NotNil(t, response.Body)
		})
	}
}

func TestURLHandlerGET(t *testing.T) {

	tests := []struct {
		name       string
		query      string
		want       string
		statusCode int
		db         map[string]string
	}{
		{
			name:       "Trying to get original link",
			query:      "/testlink",
			want:       "https://yandex.com",
			statusCode: 307,
			db: map[string]string{
				"testlink": "https://yandex.com",
			},
		},
		{
			name:       "Trying to get link that is not in map",
			query:      "/someotherlink",
			want:       "",
			statusCode: 404,
			db:         make(map[string]string),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := testRequest(tt.query, http.MethodGet, nil, tt.db)
			response := recorder.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.want, response.Header.Get("Location"))
			assert.Equal(t, tt.statusCode, response.StatusCode)
		})
	}
}

func TestPOSTGET(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		want       string
		statusCode int
		db         map[string]string
	}{
		{
			name:       "Trying to short and get original link",
			query:      "https://yandex.by",
			want:       "https://yandex.by",
			statusCode: 307,
			db:         make(map[string]string),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.query)
			postRecorder := testRequest("/", http.MethodPost, bodyReader, tt.db)
			postResponse := postRecorder.Body.String()
			defer postRecorder.Result().Body.Close()
			getRecorder := testRequest(postResponse, http.MethodGet, nil, tt.db)
			getResponse := getRecorder.Result()
			defer getResponse.Body.Close()
			assert.Equal(t, tt.want, getResponse.Header.Get("Location"))
			assert.Equal(t, tt.statusCode, getResponse.StatusCode)
		})
	}
}
