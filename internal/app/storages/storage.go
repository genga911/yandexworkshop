package storages

import (
	"fmt"

	"github.com/genga911/yandexworkshop/internal/app/config"
)

type (
	Repository interface {
		FindByValue(value string, userID string) Link
		FindByKey(key string, userID string) Link
		Create(key string, userID string) (Link, error)
		GetAll(userID string) *LinksArray
		Ping() error
		CreateBatch(batch map[string]string, userID string) (map[string]string, error)
		Delete(IDS []string, userID string) error
	}

	Link struct {
		ShortURL      string `json:"short_url"`
		OriginalURL   string `json:"original_url"`
		CorrelationID string `json:"correlation_id"`
		IsDeleted     bool
	}

	LinksArray struct {
		Links []Link
	}
)

func CreateStorage(cfg *config.Params) (Repository, error) {
	if cfg.DatabaseDSN != "" {
		storage, err := CreateDBStorage(cfg)
		if err == nil {
			fmt.Println("DB storage was choosen")
			return storage, nil
		} else {
			fmt.Println(err)
		}
	}

	fmt.Println("Memory/file storage was choosen")
	return CreateLinkStorage(cfg)
}
