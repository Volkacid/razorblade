package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StorageFile   string `env:"FILE_STORAGE_PATH" envDefault:"internal/app/storage/storage.txt"`
}

var servConf *ServerConfig

func GetServerConfig() *ServerConfig {
	if servConf == nil {
		servConf = &ServerConfig{}
		env.Parse(servConf)
		flag.StringVar(&servConf.ServerAddress, "a", "localhost:8080", "address:port to listen")
		flag.StringVar(&servConf.BaseURL, "b", "http://localhost:8080", "base url of shortener")
		flag.StringVar(&servConf.StorageFile, "f", "internal/app/storage/storage.txt", "address of db file")
		flag.Parse()
	}
	return servConf
}
