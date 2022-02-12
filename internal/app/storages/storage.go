package storages

import (
	"github.com/genga911/yandexworkshop/internal/app/heplers"
)

type Repository interface {
	FindByValue(value string) string
	FindByKey(key string) string
	Create(key string) string
}

type LinkStorage struct {
	store map[string]string
}

// Создание пустого хранилища
func CreateLinkStorage() *LinkStorage {
	var ls LinkStorage
	ls.store = make(map[string]string)
	return &ls
}

// возврат длинной ссылке по значению короткой
func (ls *LinkStorage) FindByValue(value string) string {
	for link, shortLink := range ls.store {
		if shortLink == value {
			return link
		}
	}

	return ""
}

// Возврат короткой ссылки по ключу
func (ls *LinkStorage) FindByKey(key string) string {
	return ls.store[key]
}

// Создание записи для длинной и короткой ссылок
func (ls *LinkStorage) Create(key string) string {
	shortLink, ok := ls.store[key]

	// Если ключа нет, создадим и сохраним короткую ссылку
	if !ok {
		shortLink = heplers.ShortCode(8)
		ls.store[key] = shortLink
	}

	return shortLink
}
