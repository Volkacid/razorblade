package server

import (
	"errors"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetHandler(db *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		key := chi.URLParam(request, "key")
		receivedValue, err := db.GetValue(key)
		if err != nil {
			var nfError *storage.NFError
			if errors.As(err, &nfError) {
				http.Error(writer, "Not found", http.StatusNotFound)
				return
			}
			http.Error(writer, "Unknown error", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Location", receivedValue)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}
