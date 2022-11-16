package config

import "github.com/caarlos0/env/v6"

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}

func GetServerConfig() *ServerConfig {
	servConf := &ServerConfig{}
	env.Parse(servConf)
	return servConf
}
