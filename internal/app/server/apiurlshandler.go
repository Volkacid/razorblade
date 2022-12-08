package server

import (
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"net/http"
)

func UrlsAPIHandler(storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userID := ctx.Value(config.UserID{}).(string)

		values, err := storage.GetValuesByID(userID)
		if err != nil {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		marshalledResponse, _ := json.Marshal(values)
		writer.Write(marshalledResponse)
	}
}
