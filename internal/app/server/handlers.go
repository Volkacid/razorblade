package server

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
)

type Handlers struct {
	storage  storage.Storage
	servConf *config.ServerConfig
	//URLs are placed in a buffer, from which they are removed every three seconds or when the buffer overflows
	deleteBuffer *service.URLsDeleteBuffer
}

func NewHandlersSet(ctx context.Context, storage storage.Storage) *Handlers {
	return &Handlers{storage: storage, servConf: config.GetServerConfig(), deleteBuffer: service.NewDeleteBuffer(ctx, storage)}
}
