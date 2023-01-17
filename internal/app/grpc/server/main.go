package server

import (
	"context"
	"errors"
	pb "github.com/Volkacid/razorblade/internal/app/grpc/proto"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RazorbladeService struct {
	pb.UnimplementedRazorbladeServiceServer
	DB           storage.Storage
	DeleteBuffer *service.URLsDeleteBuffer
}

func (rbs *RazorbladeService) GetOriginalURL(ctx context.Context, in *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	var response pb.GetOriginalURLResponse
	value, err := rbs.DB.GetValue(ctx, in.Key)
	if err != nil {
		var nfError *storage.NFError
		if errors.As(err, &nfError) {
			return &response, status.Error(codes.NotFound, "Value not found")
		}
		var deletedErr *storage.DeletedError
		if errors.As(err, &deletedErr) {
			return &response, status.Error(codes.PermissionDenied, "Value deleted")
		}
		return &response, status.Error(codes.Internal, "Server error")
	}
	response.ShortUrl = value
	return &response, nil
}

func (rbs *RazorbladeService) ListURLsByUserID(ctx context.Context, in *pb.ListURLsByUserIDRequest) (*pb.ListURLsByUserIDResponse, error) {
	var response pb.ListURLsByUserIDResponse
	userValues, err := rbs.DB.GetValuesByID(ctx, in.UserId)
	if err != nil {
		var nfError *storage.NFError
		if errors.As(err, &nfError) {
			return &response, status.Error(codes.NotFound, "Values not found")
		}
		return &response, status.Error(codes.Internal, "Server error")
	}
	pbVal := make([]*pb.UserURL, len(userValues))
	for i, v := range userValues {
		pbVal[i] = &pb.UserURL{Key: v.ShortURL, OriginalUrl: v.OriginalURL, UserId: in.UserId}
	}
	response.UserUrls = pbVal
	return &response, nil
}

func (rbs *RazorbladeService) CreateShortURL(ctx context.Context, in *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	var response pb.CreateShortURLResponse
	userID := metadata.ValueFromIncomingContext(ctx, "userid")
	key := service.GenerateShortString(in.OriginalUrl.OriginalUrl)
	err := rbs.DB.SaveValue(ctx, key, in.OriginalUrl.OriginalUrl, userID[0])
	if err != nil {
		return &response, status.Error(codes.Internal, "Server error")
	}
	response.ShortenedUrl = &pb.UserURL{UserId: userID[0], OriginalUrl: in.OriginalUrl.OriginalUrl, Key: key}
	return &response, nil
}

// DeleteShortURLs URLs are placed in a buffer, from which they are removed every three seconds or when the buffer overflows
func (rbs *RazorbladeService) DeleteShortURLs(ctx context.Context, in *pb.DeleteShortURLsRequest) (*emptypb.Empty, error) {
	userID := metadata.ValueFromIncomingContext(ctx, "userid")
	go rbs.DeleteBuffer.AddKeys(in.Keys, userID[0])
	return &emptypb.Empty{}, nil
}
