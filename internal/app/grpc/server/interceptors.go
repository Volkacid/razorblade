package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RazorbladeInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "Metadata retrieving error")
	}
	idArr := md.Get("UserID")
	if len(idArr) < 1 {
		return nil, status.Error(codes.Unauthenticated, "UserID not provided")
	}

	h, err := handler(ctx, req)
	return h, err
}
