package heplers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
)

// Возврат корня сайта
func GetMainLink() string {
	return fmt.Sprintf(
		"%s://%s:%s",
		os.Getenv("PROTOCOL"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
	)
}

// получим короткую ссылку из урл
func GetShortLink(req *http.Request) (string, error) {
	// провалидируем урл, ожидаем только буквы как в константе пакета codehelper
	url := req.URL.RequestURI()
	matched, err := regexp.MatchString(`^/[a-zA-Z]+$`, url)
	if err != nil || !matched {
		return "", errors.New("validation error")
	}

	return path.Base(url), nil
}

// адрес для запуска сервера, без протокола
func GetServerAddress() string {
	return fmt.Sprintf(
		"%s:%s",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
	)
}
