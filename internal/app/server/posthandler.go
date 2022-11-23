package server

import (
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func PostHandler(storage *storage.Storage) http.HandlerFunc {
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

		userID, _ := request.Cookie("UserID")

		for { //На случай, если сгенерированная последовательность уже будет занята
			foundStr := service.GenerateShortString()
			if _, err := storage.GetValue(foundStr); err != nil {
				storage.SaveValue(foundStr, str, userID.Value)
				writer.WriteHeader(http.StatusCreated)
				writer.Write([]byte(config.GetServerConfig().BaseURL + "/" + foundStr))
				break
			}
		}
	}
}
