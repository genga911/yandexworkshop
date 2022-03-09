package mocks

import (
	"net/http"
	"net/http/httptest"

	"github.com/genga911/yandexworkshop/internal/app/session"
	"github.com/gin-gonic/gin"
)

// мокаем контекст для теста
func MockGinContext(userID string, w *httptest.ResponseRecorder, r *http.Request, params []gin.Param) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	if params != nil {
		c.Params = params
	}
	c.Set("session", session.NewSession(userID))
	return c
}
