package server

import (
	"fmt"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/handlers"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	cfg := config.GetConfig()
	store := storages.CreateLinkStorage()
	getHandlers := handlers.GetHandlers{Storage: store, Config: &cfg}
	postHandlers := handlers.PostHandlers{Storage: store, Config: &cfg}
	postShortenHandlers := handlers.PostShortenHandlers{Storage: store, Config: &cfg}

	router := gin.Default()
	router.GET("/:code", getHandlers.GetHandler)
	router.POST("/", postHandlers.PostHandler)
	router.POST("/api/shorten", postShortenHandlers.PostShortenHandler)

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		fmt.Println(err)
	}

	return router
}
