package storages

import (
	"fmt"

	"github.com/genga911/yandexworkshop/internal/app/config"
)

type (
	Repository interface {
		FindByValue(value string, userID string) Link
		FindByKey(key string, userID string) Link
		Create(key string, userID string) Link
		GetAll(userID string) *LinksArray
		Ping() error
	}

	Link struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	LinksArray struct {
		Links []Link
	}
)

func CreateStorage(cfg *config.Params) (Repository, error) {
	if cfg.DatabaseDSN != "" {
		fmt.Println("DB storage was choosen")
		return CreateDBStorage(cfg)
	}
	fmt.Println("Memory/file storage was choosen")
	return CreateLinkStorage(cfg)
}
