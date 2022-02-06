package mocks

import (
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

// мокаем контекст для теста
func MockGinContext(w *httptest.ResponseRecorder, r *http.Request, store storages.Repository) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	c.Set("Store", store)

	return c
}
