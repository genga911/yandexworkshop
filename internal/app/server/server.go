package server

import (
	"flag"
	"fmt"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/handlers"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/middleware"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	cfg, cfgError := config.GetConfig()
	if cfgError != nil {
		fmt.Println(cfgError)
		panic(cfgError)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "a  string")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "b  string")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "f  string")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "d  string")
	flag.StringVar(&cfg.DatabaseDSN, "database-dsn", cfg.DatabaseDSN, "d string")
	flag.Parse()

	store, storeError := storages.CreateStorage(&cfg)
	if storeError != nil {
		fmt.Println(storeError)
		panic(storeError)
	}

	getHandlers := handlers.GetHandlers{Storage: store, Config: &cfg}
	postHandlers := handlers.PostHandlers{Storage: store, Config: &cfg}
	postShortenHandlers := handlers.PostShortenHandlers{Storage: store, Config: &cfg}
	userHandlers := handlers.UserHandlers{Storage: store, Config: &cfg}
	dbHandlers := handlers.DBHandlers{Storage: store, Config: &cfg}

	cryptoHelper := heplers.NewHelper([]byte(cfg.CookieKey))

	router := gin.Default()
	router.Use(middleware.Gzip)
	router.GET("/:code", getHandlers.GetHandler)
	router.GET("/ping", dbHandlers.Ping)

	withAuth := router.Group("/").Use(middleware.Auth(cryptoHelper, &cfg))
	{
		withAuth.POST("/", postHandlers.PostHandler)
		withAuth.POST("/api/shorten", postShortenHandlers.PostShortenHandler)
		withAuth.GET("/api/user/urls", userHandlers.Urls)
		withAuth.POST("/api/shorten/batch", postShortenHandlers.PostShortenBatchHandler)
	}

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return router
}
