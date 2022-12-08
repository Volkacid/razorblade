package server

import (
	"fmt"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetHandler(storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		key := chi.URLParam(request, "key")
		if receivedValue, err := storage.GetValue(key); err != nil {
			fmt.Println(err)
			http.Error(writer, "Not found", http.StatusNotFound)
		} else {
			writer.Header().Set("Location", receivedValue)
			writer.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
}
