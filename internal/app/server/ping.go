package server

import (
	"github.com/Volkacid/razorblade/internal/app/service"
	"net/http"
)

func PingDB() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !service.CheckDBConnection() {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Write([]byte("Success!"))
	}
}
