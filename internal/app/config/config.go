package config

import (
	"fmt"
	"math/rand"

	"github.com/caarlos0/env/v6"
)

type Params struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	CookieTTL       int    `env:"COOKIE_TTL" envDefault:"300"`
	CookieKey       string `env:"COOKIE_KEY" envDefault:"abcasdfghjklqwer"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`
}

func GetConfig() (Params, error) {
	var cfg = Params{}

	err := env.Parse(&cfg)
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
