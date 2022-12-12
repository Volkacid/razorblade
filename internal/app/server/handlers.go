package server

import (
	"github.com/Volkacid/razorblade/internal/app/config"
	"github.com/Volkacid/razorblade/internal/app/storage"
)

type Handlers struct {
	storage  storage.Storage
	servConf *config.ServerConfig
}

func NewHandlersSet(storage storage.Storage) *Handlers {
	return &Handlers{storage: storage, servConf: config.GetServerConfig()}
}
