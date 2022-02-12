package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/mocks"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	var emptyStore = storages.CreateLinkStorage()
	var emptyRouterHandler = RouterHandlers{Storage: emptyStore}

	type want struct {
		reg  *regexp.Regexp
		code int
	}

	tests := []struct {
		name string
		url  string
		link string
		want want
	}{
		{
			name: "Тест получения короткой ссылки",
			url:  config.GetMainLink(),
			link: "http://example.com",
			want: want{
				reg:  regexp.MustCompile(config.GetMainLink() + `/[a-zA-Z]{8}$`),
				code: http.StatusCreated,
			},
		},
		{
			name: "Тест получения короткой ссылки",
			url:  config.GetMainLink(),
			link: "example-com",
			want: want{
				reg:  nil,
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reader := strings.NewReader(tt.link)
			r, errWrite := http.NewRequest(http.MethodPost, tt.url, reader)
			if errWrite != nil {
				t.Errorf("Ошибка %v", errWrite)
			}

			c := mocks.MockGinContext(w, r, nil)

			Store(&emptyRouterHandler, c)
			res := w.Result()
			body, errRead := ioutil.ReadAll(res.Body)
			if errRead != nil {
				t.Errorf("Ошибка %v", errRead)
			}

			errBodyClose := res.Body.Close()
			if errBodyClose != nil {
				t.Errorf("Ошибка %v", errBodyClose)
			}

			if tt.want.reg != nil {
				assert.Regexp(t, tt.want.reg, string(body))
			}

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
