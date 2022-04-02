package handlers

import (
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
	// в случае ошибки отправим пользователю 400 код и редиректим его на главную
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if shortLink != "" {
		// проверим что ссылка есть в Storage
		finded := gh.Storage.FindByValue(shortLink, s.UserID)
		if !finded.IsEmpty() {
			if finded.IsDeleted {
				code = http.StatusGone
			} else {
				code = http.StatusTemporaryRedirect
				link = finded.OriginalURL
			}
		}
	}

	if code == http.StatusTemporaryRedirect {
		c.Redirect(code, link)
	} else {
		c.Status(code)
	}
}
