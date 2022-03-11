package handlers

import (
	"fmt"
	"net/http"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

type (
	DBHandlers struct {
		Storage storages.Repository
		Config  *config.Params
	}

	GetHandlers struct {
		Storage storages.Repository
		Config  *config.Params
	}

	PostHandlers struct {
		Storage storages.Repository
		Config  *config.Params
	}

	PostShortenHandlers struct {
		Storage storages.Repository
		Config  *config.Params
	}

	UserHandlers struct {
		Storage storages.Repository
		Config  *config.Params
	}
)

func (gh *GetHandlers) GetHandler(c *gin.Context) {
	Resolve(gh, c)
}

func (ph *PostHandlers) PostHandler(c *gin.Context) {
	Store(ph, c)
}

func (phs *PostShortenHandlers) PostShortenHandler(c *gin.Context) {
	StoreFromJSON(phs, c)
}

func (uh *UserHandlers) Urls(c *gin.Context) {
	Urls(uh, c)
}

func (dbh *DBHandlers) Ping(c *gin.Context) {
	c.Status(http.StatusOK)
	err := dbh.Storage.Ping()
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
}
