package server

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/Volkacid/razorblade/internal/app/grpc/proto"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RazorbladeService struct {
	pb.UnimplementedUsersServer
	DB           storage.Storage
	DeleteBuffer *service.URLsDeleteBuffer
}

func (rbs *RazorbladeService) GetValue(ctx context.Context, in *pb.GetValueRequest) (*pb.GetValueResponse, error) {
	var response pb.GetValueResponse
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
	response.Value = value
	return &response, nil
}

func (rbs *RazorbladeService) GetValuesByID(ctx context.Context, in *pb.GetValuesByIDRequest) (*pb.GetValuesByIDResponse, error) {
	var response pb.GetValuesByIDResponse
	userValues, err := rbs.DB.GetValuesByID(ctx, in.UserID)
	if err != nil {
		var nfError *storage.NFError
		if errors.As(err, &nfError) {
			return &response, status.Error(codes.NotFound, "Values not found")
		}
		return &response, status.Error(codes.Internal, "Server error")
	}
	strVal := make([]string, len(userValues))
	for i, v := range userValues {
		strVal[i] = fmt.Sprintf("ShortURL: %v, OrigURL: %v", v.ShortURL, v.OriginalURL)
	}
	response.FoundValues = strVal
	return &response, nil
}

func (rbs *RazorbladeService) SaveValue(ctx context.Context, in *pb.SaveValueRequest) (*pb.SaveValueResponse, error) {
	var response pb.SaveValueResponse
	userID, err := getUserID(ctx)
	if err != nil {
		return &response, status.Error(codes.PermissionDenied, "UserID not provided")
	}
	key := service.GenerateShortString(in.Value)
	err = rbs.DB.SaveValue(ctx, key, in.Value, userID)
	if err != nil {
		return &response, status.Error(codes.Internal, "Server error")
	}
	response.Key = key
	return &response, nil
}

func (rbs *RazorbladeService) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	var response pb.DeleteURLsResponse
	userID, err := getUserID(ctx)
	if err != nil {
		return &response, status.Error(codes.PermissionDenied, "UserID not provided")
	}
	go rbs.DeleteBuffer.AddKeys(in.Urls, userID)
	return &response, nil
}

func getUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		idArr := md.Get("UserID")
		if len(idArr) > 0 {
			return idArr[0], nil
		}
	}
	return "", errors.New("not found")
}
