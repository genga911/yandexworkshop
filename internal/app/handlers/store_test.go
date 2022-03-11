package handlers

import (
	"encoding/json"
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

type DefaultWant struct {
	reg  *regexp.Regexp
	code int
}

type DefaultStoreTest struct {
	name string
	link string
	want DefaultWant
}

// набор данных повторяется, так что его можно вынести из тестов
func testsProvider() []DefaultStoreTest {
	cfg, _ := config.GetConfig()
	return []DefaultStoreTest{
		{
			name: "Success тест получения короткой ссылки",
			link: "http://example.com",
			want: DefaultWant{
				reg:  regexp.MustCompile(cfg.BaseURL + `/[a-zA-Z]{8}$`),
				code: http.StatusCreated,
			},
		},
		{
			name: "Failt тест получения короткой ссылки",
			link: "example-com",
			want: DefaultWant{
				reg:  nil,
				code: http.StatusBadRequest,
			},
		},
	}
}

func TestStore(t *testing.T) {
	userID := "test"
	cfg, _ := config.GetConfig()
	var emptyStore, _ = storages.CreateStorage(&cfg)
	var emptyRouterHandler = PostHandlers{Storage: emptyStore, Config: &cfg}

	tests := testsProvider()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reader := strings.NewReader(tt.link)
			r, errWrite := http.NewRequest(http.MethodPost, cfg.BaseURL, reader)
			if errWrite != nil {
				t.Errorf("Ошибка %v", errWrite)
			}

			c := mocks.MockGinContext(userID, w, r, nil)

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

func TestStoreFromJson(t *testing.T) {
	userID := "test"
	cfg, _ := config.GetConfig()
	var emptyStore, _ = storages.CreateStorage(&cfg)
	var emptyRouterHandler = PostShortenHandlers{Storage: emptyStore, Config: &cfg}

	tests := testsProvider()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			requestBody := JSONBody{URL: tt.link}
			jsonString, _ := json.Marshal(requestBody)
			reader := strings.NewReader(string(jsonString))

			r, errWrite := http.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten", reader)
			if errWrite != nil {
				t.Errorf("Ошибка %v", errWrite)
			}

			c := mocks.MockGinContext(userID, w, r, nil)

			StoreFromJSON(&emptyRouterHandler, c)
			res := w.Result()
			body, errRead := ioutil.ReadAll(res.Body)
			if errRead != nil {
				t.Errorf("Ошибка %v", errRead)
			}

			parsedResult := JSONResult{}
			if unmError := json.Unmarshal(body, &parsedResult); unmError != nil {
				t.Errorf("Ошибка %v", unmError)
			}

			errBodyClose := res.Body.Close()
			if errBodyClose != nil {
				t.Errorf("Ошибка %v", errBodyClose)
			}

			if tt.want.reg != nil {
				assert.Regexp(t, tt.want.reg, parsedResult.Result)
			}

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
