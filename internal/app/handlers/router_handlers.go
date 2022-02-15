package handlers

import (
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

type GetHandlers struct {
	Storage storages.Repository
}

type PostHandlers struct {
	Storage storages.Repository
}

func (gh *GetHandlers) GetHandler(c *gin.Context) {
	Resolve(gh, c)
}

func (ph *PostHandlers) PostHandler(c *gin.Context) {
	Store(ph, c)
}
