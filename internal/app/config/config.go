package config

import "github.com/caarlos0/env/v6"

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	StorageFile   string `env:"FILE_STORAGE_PATH"`
}

func GetServerConfig() *ServerConfig {
	servConf := &ServerConfig{}
	env.Parse(servConf)
	servConf.StorageFile = "internal/app/storage/storage.txt"
	return servConf
}
