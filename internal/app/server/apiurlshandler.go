package server

import (
	"encoding/json"
	"fmt"
	"github.com/Volkacid/razortest/internal/app/config"
	"github.com/Volkacid/razortest/internal/app/service"
	"github.com/Volkacid/razortest/internal/app/storage"
	"net/http"
	"strings"
)

type Response struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func APIUrlsHandler(storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		//userID := ctx.Value("UserID").(string)
		userID := ctx.Value(service.UserID{}).(string)
		values, err := storage.GetValuesByID(userID)
		if err != nil {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		for _, elem := range values {
			fmt.Println("FOUND: ", elem)
			shortURL, originalURL, _ := strings.Cut(elem, ":-:")
			shortURL = config.GetServerConfig().BaseURL + "/" + shortURL
			response := Response{OriginalURL: originalURL, ShortURL: shortURL}
			marshaledResponse, _ := json.Marshal(response)
			writer.Write(marshaledResponse)
		}
	}
}
