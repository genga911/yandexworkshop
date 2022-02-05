package handlers

import (
	"fmt"
	"github.com/genga911/yandexworkshop/cmd/shortener/storages"
	"io/ioutil"
	"net/http"
)

// мохранение нового урла в хранилище
func Store(storage storages.Repository, w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	link := string(body)

	shortLink := storage.Create(link)

	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte(shortLink))
	if err != nil {
		fmt.Println(err)
		return
	}
}
