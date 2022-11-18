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
		flag.StringVar(&servConf.ServerAddress, "a", servConf.ServerAddress, "address:port to listen")
		flag.StringVar(&servConf.BaseURL, "b", servConf.BaseURL, "base url of shortener")
		flag.StringVar(&servConf.StorageFile, "f", servConf.StorageFile, "address of db file")
		flag.Parse()
	}
	return servConf
}
