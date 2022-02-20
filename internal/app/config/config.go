package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Params struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

func GetConfig() (Params, error) {
	var cfg = Params{}

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	// теперь обработаем флаги, заменим значения в конфиге
	// это костыль для инкремента 2, так как при запуске тестов флаг "а" используется, вызов StringVar вызывает ошибку
	// panic: /tmp/go-build3078548901/b184/handlers.test flag redefined: a [recovered]
	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "a string")
	}
	if flag.Lookup("b") == nil {
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "a string")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "a string")
	}

	flag.Parse()

	return cfg, nil
}
