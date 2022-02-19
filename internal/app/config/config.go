package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Params struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
}

func GetConfig() Params {
	var cfg = Params{}

	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
	}

	return cfg
}
