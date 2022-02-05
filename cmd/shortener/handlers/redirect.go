package handlers

import (
	"fmt"
	"github.com/genga911/yandexworkshop/cmd/shortener/storages"
	"net/http"
	"path"
	"strconv"
)

// редирект с короткого урла на длинный
func Redirect(storage storages.Repository, w http.ResponseWriter, req *http.Request, code int) {
	// получим короткую ссылку из урл
	shortLink := path.Base(req.URL.String())
	// проверим что ссылка есть в Storage
	link := storage.FindByValue(shortLink)

	if link != "" {
		w.WriteHeader(code)
		http.Redirect(w, req, link, code)
		fmt.Println(strconv.Itoa(code) + " redirected")
	}

	// Если короткий урл не нашелся, вернем 404
	w.WriteHeader(http.StatusNotFound)
	http.Redirect(w, req, req.Host, http.StatusNotFound)
}
