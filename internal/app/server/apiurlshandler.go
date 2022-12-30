package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"net/http"
	"time"
)

func (handlers *Handlers) UrlsAPIHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userID := ctx.Value(config.UserID{}).(string)

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
	fmt.Println("Urls handler request at: ", time.Now())
	writer.Header().Set("Content-Type", "application/json")
	marshalledResponse, _ := json.Marshal(values)
	writer.Write(marshalledResponse)
}
