package server

import (
	"fmt"
	handlers2 "github.com/genga911/yandexworkshop/internal/app/handlers"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/middleware"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ApiMiddleware(storages.Links))

	router.GET("/:shortLink", handlers2.Resolve)
	router.POST("/", handlers2.Store)

	err := router.Run(heplers.GetServerAddress())
	if err != nil {
		fmt.Println(err)
	}

	return router
}
