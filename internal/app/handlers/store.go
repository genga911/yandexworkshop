package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/session"
	"github.com/gin-gonic/gin"
)

type (
	JSONBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	JSONBatchResult struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

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
	s := session.GetSession(c)
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

	shortLink := ph.Storage.Create(link.String(), s.UserID)

	c.String(http.StatusCreated, heplers.PrepareShortLink(shortLink.ShortURL, ph.Config))
}

func StoreBatchFromJSON(phs *PostShortenHandlers, c *gin.Context) {
	s := session.GetSession(c)
	var body []JSONBatch
	err := c.ShouldBindJSON(&body)

	var result []JSONBatchResult
	c.Header(`Content-type`, gin.MIMEJSON)

	if err != nil {
		encode, _ := json.Marshal(result)
		c.String(http.StatusBadRequest, string(encode))
		return
	}

	shortLinks := make(map[string]string)
	for _, link := range body {
		_, validationError := url.ParseRequestURI(link.OriginalURL)

		if validationError != nil {
			c.JSON(http.StatusBadRequest, validationError)
			return
		}

		shortLinks[link.CorrelationID] = link.OriginalURL
	}

	shortLinks, err = phs.Storage.CreateBatch(shortLinks, s.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	for key, shorted := range shortLinks {
		result = append(result, JSONBatchResult{
			CorrelationID: key,
			ShortURL:      heplers.PrepareShortLink(shorted, phs.Config),
		})
	}

	encode, _ := json.Marshal(result)
	c.String(http.StatusCreated, string(encode))
}

func StoreFromJSON(phs *PostShortenHandlers, c *gin.Context) {
	s := session.GetSession(c)
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

	shortLink := phs.Storage.Create(body.URL, s.UserID)

	shortedLink := heplers.PrepareShortLink(shortLink.ShortURL, phs.Config)

	result.Result = shortedLink
	encode, _ := json.Marshal(result)
	c.String(http.StatusCreated, string(encode))
}
