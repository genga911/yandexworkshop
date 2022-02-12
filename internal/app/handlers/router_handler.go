package handlers

import (
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

type RouterHandlers struct {
	Storage storages.Repository
}

func (rh *RouterHandlers) GetHandler(c *gin.Context) {
	Resolve(rh, c)
}

func (rh *RouterHandlers) PostHandler(c *gin.Context) {
	Store(rh, c)
}
