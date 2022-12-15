package server

import (
	"encoding/json"
	"errors"
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

func (handlers *Handlers) PostAPIHandler(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Please make a correct request!", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	receivedURL := URL{}
	err = json.Unmarshal(body, &receivedURL)
	if err != nil || receivedURL.URL == "" || !service.ValidateURL(receivedURL.URL) {
		http.Error(writer, "Please provide a correct URL!", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	userID := ctx.Value(config.UserID{}).(string)

	var duplicateErr *storage.DuplicateError
	if key, err := handlers.storage.FindDuplicate(ctx, receivedURL.URL); errors.As(err, &duplicateErr) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		duplicateResult := Result{URL: handlers.servConf.BaseURL + "/" + key}
		marshaledResult, _ := json.Marshal(duplicateResult)
		writer.Write(marshaledResult)
		return
	}

	foundStr := service.GenerateShortString(receivedURL.URL)
	err = handlers.storage.SaveValue(ctx, foundStr, receivedURL.URL, userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	result := Result{URL: handlers.servConf.BaseURL + "/" + foundStr}
	marshaledResult, _ := json.Marshal(result)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(marshaledResult)
}
