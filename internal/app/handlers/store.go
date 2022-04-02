package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/session"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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
	status := http.StatusCreated
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

	shortLink, createError := ph.Storage.Create(link.String(), s.UserID)
	if createError != nil {
		var pgErr *pgconn.PgError
		if errors.As(createError, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				shortLink = ph.Storage.FindByKey(link.String(), s.UserID)
				status = http.StatusConflict
			}
		} else {
			c.String(http.StatusBadRequest, createError.Error())
			return
		}
	}

	c.String(status, heplers.PrepareShortLink(shortLink.ShortURL, ph.Config))
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
	status := http.StatusCreated
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

	shortLink, createError := phs.Storage.Create(body.URL, s.UserID)
	fmt.Println(createError)
	if createError != nil {
		var pgErr *pgconn.PgError
		if errors.As(createError, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				shortLink = phs.Storage.FindByKey(body.URL, s.UserID)
				status = http.StatusConflict
			}
		} else {
			c.JSON(http.StatusBadRequest, createError)
			return
		}
	}

	shortedLink := heplers.PrepareShortLink(shortLink.ShortURL, phs.Config)

	result.Result = shortedLink
	encode, _ := json.Marshal(result)
	c.String(status, string(encode))
}
