package main

import (
	"context"
	"fmt"
	"github.com/Volkacid/razorblade/internal/app/config"
	pb "github.com/Volkacid/razorblade/internal/app/grpc/proto"
	server2 "github.com/Volkacid/razorblade/internal/app/grpc/server"
	"github.com/Volkacid/razorblade/internal/app/server"
	"github.com/Volkacid/razorblade/internal/app/server/middlewares"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	servConf := config.GetServerConfig()
	db := storage.CreateStorage()

	go gRPCService(db, ctx)

	handlers := server.NewHandlersSet(db, ctx)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Use(middlewares.GetUserID)
		router.Use(middlewares.GzipHandle)
		router.Get("/", handlers.MainPage)
		router.Get("/{key}", handlers.GetHandler)
		router.Get("/ping", handlers.PingDB)
		router.Get("/api/user/urls", handlers.UrlsAPIHandler)
		router.Post("/", handlers.PostHandler)
		router.Post("/api/shorten", handlers.PostAPIHandler)
		router.Post("/api/shorten/batch", handlers.BatchHandler)
		router.Delete("/api/user/urls", handlers.DeleteUserURLs)
	})
	log.Fatal(http.ListenAndServe(servConf.ServerAddress, router))
}

func gRPCService(db storage.Storage, ctx context.Context) {
	gRPCListener, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterUsersServer(rpcServer, &server2.RazorbladeService{DB: db, DeleteBuffer: service.NewDeleteBuffer(db, ctx)})
	fmt.Println("Starting a gRPC server")
	if err = rpcServer.Serve(gRPCListener); err != nil {
		log.Fatal(err)
	}
}
