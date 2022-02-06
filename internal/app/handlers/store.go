package handlers

import (
	"fmt"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

// мохранение нового урла в хранилище
func Store(c *gin.Context) {
	store := c.MustGet("Store").(storages.Repository)

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

	shortLink := store.Create(link.String())

	shortedLink := fmt.Sprintf(
		"%s/%s",
		heplers.GetMainLink(),
		shortLink,
	)

	c.String(http.StatusCreated, shortedLink)
}
