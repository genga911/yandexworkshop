package storages

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/session"
)

type LinkStorage struct {
	store      map[string]*LinksArray
	storeMutex sync.Mutex
	file       *os.File
	writer     *bufio.Writer
}

func (l *Link) IsEmpty() bool {
	return l.ShortURL == "" && l.OriginalURL == ""
}

// Создание пустого хранилища
func CreateLinkStorage(cfg *config.Params) (*LinkStorage, error) {
	var ls LinkStorage
	ls.store = make(map[string]*LinksArray)
	ls.file = nil

	if cfg.FileStoragePath != "" {
		var file *os.File
		// проверим существует ли файл
		if _, err := os.Stat(cfg.FileStoragePath); err == nil {
			oFile, openFileError := os.OpenFile(cfg.FileStoragePath, os.O_RDWR, os.ModeAppend)
			if openFileError != nil {
				return nil, openFileError
			}
			file = oFile
		} else if errors.Is(err, os.ErrNotExist) {
			// Создадим файл
			cFile, createFileError := os.Create(cfg.FileStoragePath)
			if createFileError != nil {
				return nil, createFileError
			}
			file = cFile
		} else {
			// в случае другой ошибки
			return nil, err
		}

		ls.file = file
		ls.writer = bufio.NewWriter(ls.file)
		ls.loadFromFile()
	}

	return &ls, nil
}

// возврат ссылки по значению короткой ссылки
func (ls *LinkStorage) FindByValue(value string, userID string) Link {
	ls.storeMutex.Lock()
	defer ls.storeMutex.Unlock()
	for _, link := range ls.GetAll(userID).Links {
		if link.ShortURL == value {
			return link
		}
	}

	return Link{}
}

// Возврат короткой ссылки по длинной
func (ls *LinkStorage) FindByKey(key string, userID string) Link {
	ls.storeMutex.Lock()
	defer ls.storeMutex.Unlock()
	for _, link := range ls.GetAll(userID).Links {
		if link.OriginalURL == key {
			return link
		}
	}

	return Link{}
}

// Создание записи для длинной и короткой ссылок
func (ls *LinkStorage) Create(key string, userID string) (Link, error) {
	var err error
	shortLink := ls.FindByKey(key, userID)
	// Если ключа нет, создадим и сохраним короткую ссылку
	if shortLink.IsEmpty() {
		shortLink.ShortURL = heplers.ShortCode(8)
		shortLink.OriginalURL = key
		ls.storeMutex.Lock()
		defer ls.storeMutex.Unlock()

		ls.store[userID].Links = append(ls.store[userID].Links, shortLink)
		// Если у нас есть файл, допишем ссылку в него
		if ls.file != nil {
			err = ls.appendToFile(shortLink, userID)
		}
	}

	return shortLink, err
}

// Загрузка map из файла
func (ls *LinkStorage) loadFromFile() {
	scanner := bufio.NewScanner(ls.file)
	for scanner.Scan() {
		// данные храним в csv, для простоты разделяем через ,
		s := strings.Split(scanner.Text(), ",")
		ls.store[s[0]].Links = append(ls.store[s[0]].Links, Link{
			ShortURL:    s[0],
			OriginalURL: s[1],
		})
	}
}

// запись в конец файла
func (ls *LinkStorage) appendToFile(link Link, userID string) error {
	// чтобы не нагружать память сервера, будем записывать в файл например по 10 записей
	_, err := ls.file.WriteString(fmt.Sprintf("%s,%s,%s\n", userID, link.ShortURL, link.OriginalURL))
	return err
}

// геттер для стора
func (ls *LinkStorage) GetAll(userID string) *LinksArray {
	// идентификатор сессии гостя
	if userID == session.GuestSession {
		var results []Link
		for _, links := range ls.store {
			results = append(results, links.Links...)
		}

		return &LinksArray{
			Links: results,
		}
	}

	links, exists := ls.store[userID]
	if !exists {
		links = &LinksArray{}
		ls.store[userID] = links
	}

	return links
}

func (ls *LinkStorage) CreateBatch(batch map[string]string, userID string) (map[string]string, error) {
	for key, link := range batch {
		created, err := ls.Create(link, userID)
		if err != nil {
			return nil, err
		}
		batch[key] = created.ShortURL
	}

	return batch, nil
}

func (ls *LinkStorage) Ping() error {
	return nil
}
