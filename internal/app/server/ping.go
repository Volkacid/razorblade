package server

import (
	"github.com/Volkacid/razorblade/internal/app/storage"
	"net/http"
)

func (handlers *Handlers) PingDB(writer http.ResponseWriter, request *http.Request) {
	if !storage.CheckDBConnection() {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write([]byte("Success!"))
}
