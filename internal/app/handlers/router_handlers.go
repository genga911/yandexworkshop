package handlers

import (
	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

type (
	GetHandlers struct {
		Storage *storages.LinkStorage
		Config  *config.Params
	}

	PostHandlers struct {
		Storage *storages.LinkStorage
		Config  *config.Params
	}

	PostShortenHandlers struct {
		Storage *storages.LinkStorage
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
