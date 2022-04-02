package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/session"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

// получение всех урл пользователем
func Urls(uh *UserHandlers, c *gin.Context) {
	s := session.GetSession(c)
	links := uh.Storage.GetAll(s.UserID)
	status := http.StatusNoContent
	var result []storages.Link

	if len(links.Links) != 0 {
		status = http.StatusOK
		for _, link := range links.Links {
			result = append(result, storages.Link{
				ShortURL:    heplers.PrepareShortLink(link.ShortURL, uh.Config),
				OriginalURL: link.OriginalURL,
			})
		}
	}

	if status == http.StatusNoContent {
		// это тупо, но тест в гите ожидает именно такого ответа, при отсутствии урлов
		c.JSON(status, "{}")
		return
	}

	c.JSON(status, result)
}

func Delete(dh *DeleteHandlers, c *gin.Context) {
	s := session.GetSession(c)
	var IDS []string

	// получаем тело
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// декодируем в массив строк
	err = json.Unmarshal(body, &IDS)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// выполняем soft-delete
	err = dh.Storage.Delete(IDS, s.UserID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusAccepted, "")
}
