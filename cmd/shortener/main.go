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
	//Storage initializing. See storage configuration at internal/app/config/config.go
	db := storage.CreateStorage()

	go gRPCService(ctx, db)

	//Initiating an object to access handlers
	handlers := server.NewHandlersSet(ctx, db)

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

	//http server launch
	log.Fatal(http.ListenAndServe(servConf.ServerAddress, router))
}

func gRPCService(ctx context.Context, db storage.Storage) {
	gRPCListener, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer := grpc.NewServer(withServerUnaryInterceptor())
	pb.RegisterRazorbladeServiceServer(rpcServer, &server2.RazorbladeService{DB: db, DeleteBuffer: service.NewDeleteBuffer(ctx, db)})

	//gRPC server launch
	fmt.Println("Starting a gRPC server")
	log.Fatal(rpcServer.Serve(gRPCListener))
}

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(server2.RazorbladeInterceptor)
}
