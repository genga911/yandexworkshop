package storages

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
)

type Repository interface {
	FindByValue(value string) string
	FindByKey(key string) string
	Create(key string) string
}

type LinkStorage struct {
	store  map[string]string
	file   *os.File
	writer *bufio.Writer
}

// Создание пустого хранилища
func CreateLinkStorage(cfg *config.Params) (*LinkStorage, error) {
	var ls LinkStorage
	ls.store = make(map[string]string)
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
		// Если у нас есть файл, допишем ссылку в него
		if ls.file != nil {
			appendError := ls.appendToFile(key, shortLink)
			if appendError != nil {
				fmt.Println(appendError)
			}
		}
	}

	return shortLink
}

// Загрузка map из файла
func (ls *LinkStorage) loadFromFile() {
	scanner := bufio.NewScanner(ls.file)
	for scanner.Scan() {
		// данные храним в csv, для простоты разделяем через ,
		s := strings.Split(scanner.Text(), ",")
		ls.store[s[0]] = s[1]
	}
}

// запись в конец файла
func (ls *LinkStorage) appendToFile(key string, value string) error {
	// чтобы не нагружать память сервера, будем записывать в файл например по 10 записей
	_, err := ls.file.WriteString(fmt.Sprintf("%s,%s\n", key, value))
	return err
}
