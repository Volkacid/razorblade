package server

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

func PingDB( /*storage *storage.Storage*/ ) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		conn, err := pgx.Connect(ctx, config.GetServerConfig().DBAddress)
		if err != nil {
			panic(err)
		}
		defer conn.Close(ctx)

		if err = conn.Ping(ctx); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Write([]byte("Success!"))
	}
}
