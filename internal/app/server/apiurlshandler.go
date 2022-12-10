package server

import (
	"encoding/json"
	"errors"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"net/http"
)

func UrlsAPIHandler(db *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userID := ctx.Value(config.UserID{}).(string)

		values, err := db.GetValuesByID(userID)
		if err != nil {
			var nfError *storage.NFError
			if errors.As(err, &nfError) {
				http.Error(writer, "", http.StatusNoContent)
				return
			}
			http.Error(writer, "", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		marshalledResponse, _ := json.Marshal(values)
		writer.Write(marshalledResponse)
	}
}
