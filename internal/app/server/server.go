package server

import (
	"fmt"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/handlers"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	rooterHandlers := handlers.RouterHandlers{Storage: storages.CreateLinkStorage()}
	router := gin.Default()
	router.GET("/:code", rooterHandlers.GetHandler)
	router.POST("/", rooterHandlers.PostHandler)

	err := router.Run(config.GetServerAddress())
	if err != nil {
		fmt.Println(err)
	}

	return router
}
