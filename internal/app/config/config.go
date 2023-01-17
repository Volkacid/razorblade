package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StorageFile   string `env:"FILE_STORAGE_PATH" envDefault:"internal/app/storage/storage.txt"`
	DBAddress     string `env:"DATABASE_DSN" envDefault:"postgres://postgres:practicum@localhost:5432/postgres"`
}

var servConf *ServerConfig

func GetServerConfig() *ServerConfig {
	if servConf == nil {
		servConf = &ServerConfig{}
		err := env.Parse(servConf)
		if err != nil {
			return nil
		}
		flag.StringVar(&servConf.ServerAddress, "a", servConf.ServerAddress, "address:port to listen")
		flag.StringVar(&servConf.BaseURL, "b", servConf.BaseURL, "base url of shortener")
		flag.StringVar(&servConf.StorageFile, "f", servConf.StorageFile, "address of db file")
		flag.StringVar(&servConf.DBAddress, "d", servConf.DBAddress, "user:pass@address:port of PostgreSQL db")
		flag.Parse()
	}
	return servConf
}
