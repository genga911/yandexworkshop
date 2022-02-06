package handlers

import (
	"fmt"
	"github.com/genga911/yandexworkshop/internal/app/handlers/mocks"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

func init() {
	// подключаем тестовое окружение
	err := godotenv.Load(os.ExpandEnv("$GOPATH/yandexworkshop/.env.testing"))
	if err != nil {
		fmt.Println("Ошибка загрузки .env: " + err.Error())
	}
}

func TestStore(t *testing.T) {
	var emptyStore = storages.CreateLinkStorage()
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
			url:  heplers.GetMainLink(),
			link: "http://example.com",
			want: want{
				reg:  regexp.MustCompile(heplers.GetMainLink() + `/[a-zA-Z]{8}$`),
				code: http.StatusCreated,
			},
		},
		{
			name: "Тест получения короткой ссылки",
			url:  heplers.GetMainLink(),
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
			r, err := http.NewRequest(http.MethodPost, tt.url, reader)
			if err != nil {
				t.Errorf("Ошибка %v", err)
			}

			c := mocks.MockGinContext(w, r, emptyStore)

			Store(c)
			res := w.Result()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Ошибка %v", err)
			}

			if tt.want.reg != nil {
				assert.Regexp(t, tt.want.reg, string(body))
			}

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
