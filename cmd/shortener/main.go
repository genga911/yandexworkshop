package main

import (
	"fmt"
	"github.com/genga911/yandexworkshop/cmd/shortener/handlers"
	"github.com/genga911/yandexworkshop/cmd/shortener/storages"
	"net/http"
)

var Links = storages.LinkStorage{}

// фабрика для разбора запросов
func resolve(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		handlers.Redirect(&Links, w, req, http.StatusTemporaryRedirect)
	case "POST":
		handlers.Store(&Links, w, req)
	default:
		handlers.Redirect(&Links, w, req, http.StatusBadRequest)
	}
}

func main() {
	Links.Init()

	http.HandleFunc("/", resolve)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
