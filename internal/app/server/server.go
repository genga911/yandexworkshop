package server

import (
	"fmt"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/handlers"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	store := storages.CreateLinkStorage()
	getHandlers := handlers.GetHandlers{Storage: store}
	postHandlers := handlers.PostHandlers{Storage: store}

	router := gin.Default()
	router.GET("/:code", getHandlers.GetHandler)
	router.POST("/", postHandlers.PostHandler)

	err := router.Run(config.GetServerAddress())
	if err != nil {
		fmt.Println(err)
	}

	return router
}
