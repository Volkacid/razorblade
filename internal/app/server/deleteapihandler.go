package server

import (
	"encoding/json"
	"github.com/Volkacid/razorblade/internal/app/config"
	"io"
	"net/http"
)

func (handlers *Handlers) DeleteUserURLs(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil || len(body) == 0 {
		http.Error(writer, "Please make a correct request", http.StatusBadRequest)
		return
	}
	ctx := request.Context()
	userID := ctx.Value(config.UserID{}).(string)

	var urls []string
	if err := json.Unmarshal(body, &urls); err != nil {
		http.Error(writer, "Please make a correct request", http.StatusBadRequest)
		return
	}
	go handlers.deleteBuffer.AddKeys(urls, userID)

	writer.WriteHeader(http.StatusAccepted)
	writer.Write([]byte("Success!"))
}
