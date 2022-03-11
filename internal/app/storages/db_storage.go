package storages

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/genga911/yandexworkshop/internal/app/config"
)

type DBStorage struct {
	connection *sql.DB
}

func (dbs *DBStorage) connect(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	dbs.connection = db
	return nil
}

// Создание хранилища
func CreateDBStorage(cfg *config.Params) (*DBStorage, error) {
	storage := DBStorage{}
	dbe := storage.connect(cfg.DatabaseDSN)
	return &storage, dbe
}

// возврат ссылки по значению короткой ссылки
func (dbs *DBStorage) FindByValue(value string, userID string) Link {
	return Link{}
}

// Возврат короткой ссылки по длинной
func (dbs *DBStorage) FindByKey(key string, userID string) Link {
	return Link{}
}

// Создание записи для длинной и короткой ссылок
func (dbs *DBStorage) Create(key string, userID string) Link {
	return Link{}
}

// геттер для стора
func (dbs *DBStorage) GetAll(userID string) *LinksArray {
	return &LinksArray{}
}

func (dbs *DBStorage) Ping() error {
	return dbs.connection.Ping()
}
