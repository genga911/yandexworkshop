package config

import (
	"crypto/aes"
	"flag"
	"fmt"
	"math/rand"

	"github.com/caarlos0/env/v6"
)

type Params struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	CookieTTL       int    `env:"COOKIE_TTL" envDefault:"300"`
	CookieKey       []byte
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`
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
	if flag.Lookup("d") == nil {
		flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "a string")
	}

	flag.Parse()

	cfg.CookieKey, err = generateRandom(aes.BlockSize) // ключ шифрования
	if err != nil {
		fmt.Println(err)
		return cfg, err
	}

	return cfg, nil
}

// функция взята из обучающих материалов
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
