package storages

import (
	"context"
	"fmt"
	"strings"

	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/session"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"

	"github.com/genga911/yandexworkshop/internal/app/config"
)

const LinksTable = "links"

type DBStorage struct {
	connection *pgx.Conn
}

func (dbs *DBStorage) deleteTable() error {
	query :=
		"DROP TABLE IF EXISTS links;"

	_, err := dbs.connection.Exec(context.Background(), query)

	if err != nil {
		return err
	}

	return nil
}

func (dbs *DBStorage) createTable() error {
	derr := dbs.deleteTable()
	if derr != nil {
		return derr
	}

	query :=
		"CREATE TABLE IF NOT EXISTS links (" +
			"id              SERIAL PRIMARY KEY," +
			"user_id varchar(255) NOT NULL," +
			"short_url varchar(255) NOT NULL," +
			"original_url varchar(255) UNIQUE NOT NULL," +
			"correlation_id varchar(255) UNIQUE, " +
			"is_deleted boolean DEFAULT false NOT NULL" +
			");"

	_, err := dbs.connection.Exec(context.Background(), query)

	return err
}

func (dbs *DBStorage) connect(connStr string) error {
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return err
	}

	dbs.connection = db

	return nil
}

// Создание хранилища
func CreateDBStorage(cfg *config.Params) (*DBStorage, error) {
	storage := DBStorage{}
	// подключимся к БД
	dbe := storage.connect(cfg.DatabaseDSN)
	if dbe != nil {
		fmt.Println(dbe)
		return nil, dbe
	}

	// создадим таблицу
	dbe = storage.createTable()
	if dbe != nil {
		return nil, dbe
	}

	return &storage, nil
}

// возврат ссылки по значению короткой ссылки
func (dbs *DBStorage) FindByValue(value string, userID string) Link {
	link := Link{}
	query := fmt.Sprintf("SELECT short_url, original_url, is_deleted FROM %s WHERE short_url = $1", LinksTable)
	var args []interface{}
	args = append(args, value)
	if userID != session.GuestSession {
		query += " AND user_id = $2"
		args = append(args, userID)
	}

	res := dbs.connection.QueryRow(
		context.Background(),
		query,
		args...,
	)

	err := res.Scan(&link.ShortURL, &link.OriginalURL, &link.IsDeleted)
	if err != nil {
		fmt.Println(err)
		return Link{}
	}

	return link
}

// Возврат короткой ссылки по длинной
func (dbs *DBStorage) FindByKey(key string, userID string) Link {
	link := Link{}
	query := fmt.Sprintf("SELECT short_url, original_url, is_deleted FROM %s WHERE original_url = $1", LinksTable)
	var args []interface{}
	args = append(args, key)
	if userID != session.GuestSession {
		query += " AND user_id = $2"
		args = append(args, userID)
	}

	res := dbs.connection.QueryRow(
		context.Background(),
		query,
		args...,
	)

	err := res.Scan(&link.ShortURL, &link.OriginalURL, &link.IsDeleted)
	if err != nil {
		fmt.Println(err)
		return Link{}
	}

	return link
}

// Создание записи для длинной и короткой ссылок
func (dbs *DBStorage) Create(key string, userID string) (Link, error) {
	shortLink := Link{}
	shortLink.ShortURL = heplers.ShortCode(8)
	shortLink.OriginalURL = key
	_, err := dbs.connection.Exec(
		context.Background(),
		fmt.Sprintf("INSERT INTO %s(id, user_id, short_url, original_url) VALUES(DEFAULT, $1, $2, $3)", LinksTable),
		userID,
		shortLink.ShortURL,
		shortLink.OriginalURL,
	)

	if err != nil {
		return Link{}, err
	}

	return shortLink, nil
}

// геттер для стора
func (dbs *DBStorage) GetAll(userID string) *LinksArray {
	query := fmt.Sprintf("SELECT short_url, original_url, is_deleted FROM %s", LinksTable)
	var args []interface{}
	if userID != session.GuestSession {
		query += " WHERE user_id = $1"
		args = append(args, userID)
	}

	query += " ORDER BY correlation_id ASC"

	var links LinksArray
	results, err := dbs.connection.Query(context.Background(), query, args...)
	if err != nil {
		fmt.Println(err)
		return &LinksArray{}
	}

	defer results.Close()

	for results.Next() {
		var link Link
		serr := results.Scan(&link.ShortURL, &link.OriginalURL, &link.IsDeleted)
		if serr != nil {
			fmt.Println(serr)
			return &LinksArray{}
		}

		links.Links = append(links.Links, link)
	}

	return &links
}

func (dbs *DBStorage) CreateBatch(batch map[string]string, userID string) (map[string]string, error) {
	query := fmt.Sprintf("INSERT INTO %s(id, user_id, short_url, original_url, correlation_id) VALUES", LinksTable)
	var args []interface{}
	result := make(map[string]string)
	for key, link := range batch {
		shortCode := heplers.ShortCode(8)
		args = append(args, userID, shortCode, link, key)
		l := len(args)
		// конструктор не хочет работать с ? придется использовать $n
		query += fmt.Sprintf("(DEFAULT, $%d, $%d, $%d, $%d),", l-3, l-2, l-1, l)
		result[key] = shortCode
	}
	query = strings.Trim(query, ",") + " ON CONFLICT ON CONSTRAINT links_original_url_key DO NOTHING;"

	_, err := dbs.connection.Exec(context.Background(), query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (dbs *DBStorage) Ping() error {
	return dbs.connection.Ping(context.Background())
}

func (dbs *DBStorage) Delete(IDS []string, userID string) error {
	var err error
	preparedIDs := &pgtype.TextArray{}
	err = preparedIDs.Set(IDS)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET is_deleted=true WHERE user_id=$1 AND short_url = any($2)", LinksTable)
	_, err = dbs.connection.Exec(context.Background(), query, userID, preparedIDs)

	return err
}
