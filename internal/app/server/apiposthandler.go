package server

import (
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"io"
	"net/http"
)

type URL struct {
	URL string `json:"url"`
}

type Result struct {
	URL string `json:"result"`
}

func APIPostHandler(storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			http.Error(writer, "Please make a correct request!", http.StatusBadRequest)
			return
		}

		receivedURL := URL{}
		err = json.Unmarshal(body, &receivedURL)
		if err != nil || receivedURL.URL == "" || !service.ValidateURL(receivedURL.URL) {
			http.Error(writer, "Please provide a correct URL!", http.StatusBadRequest)
			return
		}

		ctx := request.Context()
		userID := ctx.Value("UserID").(string)

		for { //На случай, если сгенерированная последовательность уже будет занята
			foundStr := service.GenerateShortString()
			if _, err := storage.GetValue(foundStr); err != nil {
				storage.SaveValue(foundStr, receivedURL.URL, userID)
				result := Result{URL: config.GetServerConfig().BaseURL + "/" + foundStr}
				marshaledResult, _ := json.Marshal(result)
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusCreated)
				writer.Write(marshaledResult)
				break
			}
		}
	}
}
