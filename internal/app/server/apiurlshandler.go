package server

import (
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
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
		userID := ctx.Value("UserID").(string)
		values, err := storage.GetValuesByID(userID)
		if err != nil {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		for _, elem := range values {
			shortURL, originalURL, _ := strings.Cut(elem, ":-:")
			shortURL = config.GetServerConfig().BaseURL + "/" + shortURL
			response := Response{OriginalURL: originalURL, ShortURL: shortURL}
			marshaledResponse, _ := json.Marshal(response)
			writer.Write(marshaledResponse)
		}
	}
}
