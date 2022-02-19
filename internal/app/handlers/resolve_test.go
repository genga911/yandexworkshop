package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/mocks"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	cfg, _ := config.GetConfig()
	linkWithCode := fmt.Sprintf("%s/%s", cfg.BaseURL, "AaSsDd")
	var emptyStore, _ = storages.CreateLinkStorage(&cfg)
	var notEmptyStore, _ = storages.CreateLinkStorage(&cfg)

	var emptyRouterHandler = GetHandlers{Storage: emptyStore, Config: &cfg}
	var notEmptyRouterHandler = GetHandlers{Storage: notEmptyStore, Config: &cfg}

	code := notEmptyStore.Create(linkWithCode)
	linkWithCode = fmt.Sprintf("%s/%s", cfg.BaseURL, code)

	tests := []struct {
		name string
		want int
		url  string
		code string
		rh   *GetHandlers
	}{
		{
			name: "Нет ссылки в URL",
			want: http.StatusBadRequest,
			url:  cfg.BaseURL,
			code: "",
			rh:   &emptyRouterHandler,
		},
		{
			name: "Ссылка в URL не корректного формата",
			want: http.StatusBadRequest,
			url:  fmt.Sprintf("%s/%s", cfg.BaseURL, "/123456/789"),
			code: "",
			rh:   &emptyRouterHandler,
		},
		{
			name: "Ссылка в URL корректного формата",
			want: http.StatusTemporaryRedirect,
			url:  linkWithCode,
			code: code,
			rh:   &notEmptyRouterHandler,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)

			var params []gin.Param
			params = append(params, gin.Param{
				Key:   "code",
				Value: tt.code,
			})

			c := mocks.MockGinContext(w, r, params)
			Resolve(tt.rh, c)
			res := w.Result()

			_, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Ошибка %v", err)
			}

			errBodyClose := res.Body.Close()
			if err != nil {
				t.Errorf("Ошибка %v", errBodyClose)
			}

			assert.True(t, tt.want == res.StatusCode)
		})
	}
}
