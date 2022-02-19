package config

import (
	"github.com/caarlos0/env/v6"
)

type Params struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"../tmp/"`
	FileStorageName string `env:"FILE_STORAGE_NAME" envDefault:"links.csv"`
}

func GetConfig() (Params, error) {
	var cfg = Params{}

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
