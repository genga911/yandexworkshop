package storages

import "github.com/genga911/yandexworkshop/cmd/shortener/heplers"

type Repository interface {
	FindByValue(value string) string
	FindByKey(key string) string
	Create(key string) string
}

type LinkStorage struct {
	Store map[string]string
}

func (ls *LinkStorage) Init() {
	ls.Store = make(map[string]string)
}

// возврат длинной ссылке по значению короткой
func (ls *LinkStorage) FindByValue(value string) string {
	for link, shortLink := range ls.Store {
		if shortLink == value {
			return link
		}
	}

	return ""
}

// Возврат короткой ссылки по ключу
func (ls *LinkStorage) FindByKey(key string) string {
	return ls.Store[key]
}

// Создание записи для длинной и короткой ссылок
func (ls *LinkStorage) Create(key string) string {
	shortLink, ok := ls.Store[key]

	// Если ключа нет, создадим и сохраним короткую ссылку
	if !ok {
		shortLink = heplers.ShortCode(8)
		ls.Store[key] = shortLink
	}

	return shortLink
}
