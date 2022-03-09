package handlers

import (
	"fmt"
	"net/http"

	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/session"
	"github.com/gin-gonic/gin"
)

// псевдоредирект с короткого урла на длинный
func Resolve(gh *GetHandlers, c *gin.Context) {
	s := session.GetSession(c)
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
		finded := gh.Storage.FindByValue(shortLink, s.UserID)
		if !finded.IsEmpty() {
			code = http.StatusTemporaryRedirect
			link = finded.OriginalURL
		}
	}

	if code != http.StatusNotFound {
		c.Redirect(code, link)
	} else {
		c.Status(http.StatusNotFound)
	}
}
