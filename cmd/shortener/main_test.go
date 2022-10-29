package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
			request := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			writer := httptest.NewRecorder()
			handler := http.HandlerFunc(URLHandler(tt.db))
			handler.ServeHTTP(writer, request)
			response := writer.Result()
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
			request := httptest.NewRequest(http.MethodGet, tt.query, nil)
			writer := httptest.NewRecorder()
			handler := http.HandlerFunc(URLHandler(tt.db))
			handler.ServeHTTP(writer, request)
			response := writer.Result()
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
			postRequest := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			writer := httptest.NewRecorder()
			handler := http.HandlerFunc(URLHandler(tt.db))
			handler.ServeHTTP(writer, postRequest)
			postResponse := writer.Body.String()
			getRequest := httptest.NewRequest(http.MethodGet, postResponse, nil)

			writer = httptest.NewRecorder()
			handler.ServeHTTP(writer, getRequest)
			response := writer.Result()
			defer response.Body.Close()
			assert.Equal(t, tt.want, response.Header.Get("Location"))
			assert.Equal(t, tt.statusCode, response.StatusCode)
		})
	}
}
