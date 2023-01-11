package server

import (
	"encoding/json"
	"errors"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"net/http"
)

func (handlers *Handlers) UrlsAPIHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userID := ctx.Value(config.UserID{}).(string)
	writer.Header().Set("Content-Type", "application/json")

	values, err := handlers.storage.GetValuesByID(ctx, userID)
	if err != nil {
		var nfError *storage.NFError
		if errors.As(err, &nfError) {
			http.Error(writer, "", http.StatusNoContent)
			return
		}
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	marshalledResponse, _ := json.Marshal(values)
	writer.Write(marshalledResponse)
}
