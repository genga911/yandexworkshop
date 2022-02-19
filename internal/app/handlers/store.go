package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/gin-gonic/gin"
)

// мохранение нового урла в хранилище
func Store(ph *PostHandlers, c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	link, validationError := url.ParseRequestURI(string(body))

	if validationError != nil {
		c.String(http.StatusBadRequest, validationError.Error())
		return
	}

	shortLink := ph.Storage.Create(link.String())

	shortedLink := fmt.Sprintf(
		"%s/%s",
		config.GetMainLink(),
		shortLink,
	)

	c.String(http.StatusCreated, shortedLink)
}
