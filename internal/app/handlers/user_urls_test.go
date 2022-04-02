package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/mocks"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/stretchr/testify/assert"
)

func TestUrls(t *testing.T) {
	userID := "test"
	anotherUserID := "test2"

	cfg, _ := config.GetConfig()
	var emptyStore, _ = storages.CreateStorage(&cfg)
	uh1 := UserHandlers{
		Storage: emptyStore,
		Config:  &cfg,
	}

	var notEmptyStore, _ = storages.CreateStorage(&cfg)
	linkWithCode := fmt.Sprintf("%s/%s", cfg.BaseURL, "AaSsDd")
	link, _ := notEmptyStore.Create(linkWithCode, userID)
	// добавим еще ссылку в нагрузку
	notEmptyStore.Create(linkWithCode, anotherUserID)
	uh2 := UserHandlers{
		Storage: notEmptyStore,
		Config:  &cfg,
	}

	tests := []struct {
		name   string
		uh     *UserHandlers
		want   int
		result []storages.Link
	}{
		{
			name:   "Запрс при пустом storage",
			want:   http.StatusNoContent,
			result: nil,
			uh:     &uh1,
		},
		{
			name: "Запрс при заполненном storage",
			want: http.StatusOK,
			result: []storages.Link{
				{
					ShortURL:    heplers.PrepareShortLink(link.ShortURL, &cfg),
					OriginalURL: linkWithCode,
				},
			},

			uh: &uh2,
		},
	}

	apiLink := fmt.Sprintf("%s%s", cfg.BaseURL, "/api/user/urls")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reader := strings.NewReader(apiLink)
			r, errWrite := http.NewRequest(http.MethodGet, cfg.BaseURL, reader)
			if errWrite != nil {
				t.Errorf("Ошибка %v", errWrite)
			}

			c := mocks.MockGinContext(userID, w, r, nil)

			Urls(tt.uh, c)
			res := w.Result()
			body, errRead := ioutil.ReadAll(res.Body)
			if errRead != nil {
				t.Errorf("Ошибка %v", errRead)
			}

			errBodyClose := res.Body.Close()
			if errBodyClose != nil {
				t.Errorf("Ошибка %v", errBodyClose)
			}

			if tt.result != nil {
				var parsedResult []storages.Link
				if unmError := json.Unmarshal(body, &parsedResult); unmError != nil {
					t.Errorf("Ошибка %v", unmError)
				}

				assert.Equal(t, parsedResult, tt.result)
			}

			assert.Equal(t, tt.want, res.StatusCode)
		})
	}
}
