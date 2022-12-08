package api

import (
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"io"
	"net/http"
)

type BatchQuery struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func BatchHandler(storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			http.Error(writer, "Please make a correct request!", http.StatusBadRequest)
			return
		}

		var query []BatchQuery
		err = json.Unmarshal(body, &query)
		if err != nil {
			http.Error(writer, "Cannot parse provided URLs", http.StatusInternalServerError)
			return
		}

		ctx := request.Context()
		userID := ctx.Value(config.UserID{}).(string)
		var response []BatchResponse
		batchValues := make(map[string]string)

		for _, q := range query {
			foundStr := service.GenerateShortString(storage)
			//storage.SaveValue(foundStr, q.OriginalURL, userID)
			batchValues[foundStr] = q.OriginalURL
			response = append(response, BatchResponse{CorrelationID: q.CorrelationID, ShortURL: config.GetServerConfig().BaseURL + "/" + foundStr})
		}
		err = storage.BatchSave(batchValues, userID)
		if err != nil {
			http.Error(writer, "Unknown error", http.StatusInternalServerError)
			return
		}
		marshalledResponse, _ := json.Marshal(response)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(marshalledResponse)
	}
}
