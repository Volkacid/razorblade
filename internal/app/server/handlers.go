package server

import (
	"context"
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/service"
	"github.com/Volkacid/razorblade/internal/app/storage"
)

type Handlers struct {
	storage      storage.Storage
	servConf     *config.ServerConfig
	deleteBuffer *service.URLsDeleteBuffer
}

func NewHandlersSet(storage storage.Storage, ctx context.Context) *Handlers {
	return &Handlers{storage: storage, servConf: config.GetServerConfig(), deleteBuffer: service.NewDeleteBuffer(storage, ctx)}
}
