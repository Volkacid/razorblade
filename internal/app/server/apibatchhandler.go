package server

import (
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/service"
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

func (handlers *Handlers) BatchHandler(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, "Please make a correct request", http.StatusBadRequest)
		return
	}

	var query []BatchQuery
	err = json.Unmarshal(body, &query)
	if err != nil {
		http.Error(writer, "Cannot parse provided URLs", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	userID := ctx.Value(config.UserID).(string)
	response := make([]BatchResponse, len(query))
	batchValues := make(map[string]string, len(query))

	for i, q := range query {
		if !service.ValidateURL(q.OriginalURL) {
			http.Error(writer, "Incorrect URL", http.StatusBadRequest)
			return
		}
		foundStr := service.GenerateShortString(q.OriginalURL)
		batchValues[foundStr] = q.OriginalURL
		response[i] = BatchResponse{CorrelationID: q.CorrelationID, ShortURL: handlers.servConf.BaseURL + "/" + foundStr}
	}
	err = handlers.storage.BatchSave(request.Context(), batchValues, userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	marshalledResponse, _ := json.Marshal(response)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(marshalledResponse)
}
