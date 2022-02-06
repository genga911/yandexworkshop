package middleware

import (
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
)

func ApiMiddleware(store storages.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Store", store)
		c.Next()
	}
}
