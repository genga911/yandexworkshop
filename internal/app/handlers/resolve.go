package handlers

import (
	"net/http"

	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

// псевдоредирект с короткого урла на длинный
func Resolve(c *gin.Context) {
	store := c.MustGet("Store").(storages.Repository)
	// ссылка поумолчанию на корень сайта
	link := heplers.GetMainLink()
	// поумолчанию выставим код 404
	code := http.StatusNotFound
	shortLink, err := heplers.GetShortLink(c.Request)

	// в случае ошибки отправим пользователю 400 код и редиректим его на главную
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if shortLink != "" {
		// проверим что ссылка есть в Storage
		link = store.FindByValue(shortLink)
		if link != "" {
			code = http.StatusTemporaryRedirect
		}
	}

	c.String(code, link)
}
