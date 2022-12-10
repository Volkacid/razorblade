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

func PostHandler(db *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var str string
		body, err := io.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			http.Error(writer, "Please provide a correct URL!", http.StatusBadRequest)
			return
		}

		str = string(body)
		if strings.Contains(str, "URL=") {
			_, str, _ = strings.Cut(string(body), "URL=")
		}
		str, _ = url.QueryUnescape(str)
		if !service.ValidateURL(str) {
			http.Error(writer, "Incorrect URL!", http.StatusBadRequest)
			return
		}

		var duplicateErr *storage.DuplicateError
		if key, err := db.FindDuplicate(str); errors.As(err, &duplicateErr) {
			writer.WriteHeader(http.StatusConflict)
			writer.Write([]byte(config.GetServerConfig().BaseURL + "/" + key))
			return
		}

		ctx := request.Context()
		userID := ctx.Value(config.UserID{}).(string)

		foundStr := service.GenerateShortString(db)
		db.SaveValue(foundStr, str, userID)
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(config.GetServerConfig().BaseURL + "/" + foundStr))
	}
}
