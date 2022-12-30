package server

import (
	"errors"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (handlers *Handlers) PostHandler(writer http.ResponseWriter, request *http.Request) {
	var str string
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Please provide a correct URL!", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	str = string(body)
	if strings.Contains(str, "URL=") {
		_, str, _ = strings.Cut(string(body), "URL=")
	}
	str, _ = url.QueryUnescape(str)
	if !service.ValidateURL(str) {
		http.Error(writer, "Incorrect URL!", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	userID := ctx.Value(config.UserID{}).(string)

	var duplicateErr *storage.DuplicateError
	if key, err := handlers.storage.FindDuplicate(ctx, str); errors.As(err, &duplicateErr) {
		writer.WriteHeader(http.StatusConflict)
		writer.Write([]byte(handlers.servConf.BaseURL + "/" + key))
		return
	}

	foundStr := service.GenerateShortString(str)
	err = handlers.storage.SaveValue(ctx, foundStr, str, userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(handlers.servConf.BaseURL + "/" + foundStr))
}
