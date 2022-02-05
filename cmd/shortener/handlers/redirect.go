package handlers

import (
	"fmt"
	"github.com/genga911/yandexworkshop/cmd/shortener/storages"
	"net/http"
	"path"
	"regexp"
	"strconv"
)

// редирект с короткого урла на длинный
func Redirect(storage storages.Repository, w http.ResponseWriter, req *http.Request, code int) {
	shortLink := getShortLink(req)
	fmt.Println(shortLink)
	if shortLink != "" {
		// проверим что ссылка есть в Storage
		link := storage.FindByValue(shortLink)

		if link != "" {
			w.WriteHeader(code)
			http.Redirect(w, req, link, code)
			fmt.Println(strconv.Itoa(code) + " redirected")
		}
	}

	// Если короткий урл не нашелся, вернем 404
	w.WriteHeader(http.StatusNotFound)
	http.Redirect(w, req, req.Host, http.StatusNotFound)
}

// получим короткую ссылку из урл
func getShortLink(req *http.Request) string {
	// провалидируем урл, ожидаем только буквы как в константе пакета codehelper
	url := req.URL.RequestURI()
	matched, err := regexp.MatchString(`^/[a-zA-Z]+$`, url)
	if err != nil || !matched {
		return ""
	}

	return path.Base(url)
}
