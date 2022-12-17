package server

import (
	"errors"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (handlers *Handlers) GetHandler(writer http.ResponseWriter, request *http.Request) {
	key := chi.URLParam(request, "key")
	receivedValue, err := handlers.storage.GetValue(request.Context(), key)
	if err != nil {
		var nfError *storage.NFError
		if errors.As(err, &nfError) {
			http.Error(writer, "Not found", http.StatusNotFound)
			return
		}
		var deletedErr *storage.DeletedError
		if errors.As(err, &deletedErr) {
			http.Error(writer, "Value deleted", http.StatusGone)
			return
		}
		http.Error(writer, "Unknown error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Location", receivedValue)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}
