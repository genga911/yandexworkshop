package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/gin-gonic/gin"
)

type (
	JSONBody struct {
		URL string `from:"url" json:"url" binding:"required"`
	}

	JSONResult struct {
		Result string `json:"result,omitempty"`
		Error  string `json:"error,omitempty"`
	}
)

// cохранение нового урла в хранилище
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

func StoreFromJSON(ph *PostShortenHandlers, c *gin.Context) {
	body := JSONBody{}
	err := c.ShouldBindJSON(&body)
	result := JSONResult{}
	c.Header(`Content-type`, gin.MIMEJSON)

	if err != nil {
		encode, _ := json.Marshal(result)
		c.String(http.StatusBadRequest, string(encode))
		return
	}

	_, validationError := url.ParseRequestURI(body.URL)

	if validationError != nil {
		encode, _ := json.Marshal(result)
		c.String(http.StatusBadRequest, string(encode))
		return
	}

	shortLink := ph.Storage.Create(body.URL)

	shortedLink := fmt.Sprintf(
		"%s/%s",
		config.GetMainLink(),
		shortLink,
	)

	result.Result = shortedLink
	encode, _ := json.Marshal(result)
	c.String(http.StatusCreated, string(encode))
}
