package handlers

import (
	"fmt"
	"net/http"

	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/gin-gonic/gin"
)

// псевдоредирект с короткого урла на длинный
func Resolve(gh *GetHandlers, c *gin.Context) {
	// ссылка поумолчанию на корень сайта
	link := gh.Config.BaseURL
	// поумолчанию выставим код 404
	code := http.StatusNotFound
	shortLink, err := heplers.GetShortLink(c)
	fmt.Println(err)
	// в случае ошибки отправим пользователю 400 код и редиректим его на главную
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if shortLink != "" {
		// проверим что ссылка есть в Storage
		link = gh.Storage.FindByValue(shortLink)
		if link != "" {
			code = http.StatusTemporaryRedirect
		}
	}

	c.Redirect(code, link)
}
