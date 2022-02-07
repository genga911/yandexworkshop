package server

import (
	"fmt"

	"github.com/genga911/yandexworkshop/internal/app/handlers"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/middleware"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.APIMiddleware(storages.Links))

	router.GET("/:shortLink", handlers.Resolve)
	router.POST("/", handlers.Store)

	err := router.Run(heplers.GetServerAddress())
	if err != nil {
		fmt.Println(err)
	}

	return router
}
