package heplers

import (
	"errors"
	"fmt"
	"github.com/genga911/yandexworkshop/internal/app/constants"
	"net/http"
	"path"
	"regexp"
)

// Возврат корня сайта
func GetMainLink() string {
	return fmt.Sprintf(
		"%s://%s:%s",
		constants.PROTOCOL,
		constants.HOST,
		constants.PORT,
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
		constants.HOST,
		constants.PORT,
	)
}
