package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/genga911/yandexworkshop/internal/app/handlers/mocks"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/storages"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	linkWithID := fmt.Sprintf("%s/%s", heplers.GetMainLink(), "AaSsDd")
	var emptyStore = storages.CreateLinkStorage()
	var notEmptyStore = storages.CreateLinkStorage()
	linkWithID = fmt.Sprintf("%s/%s", heplers.GetMainLink(), notEmptyStore.Create(linkWithID))

	tests := []struct {
		name  string
		want  int
		url   string
		store *storages.LinkStorage
	}{
		{
			name:  "Нет ссылки в URL",
			want:  http.StatusBadRequest,
			url:   heplers.GetMainLink(),
			store: emptyStore,
		},
		{
			name: "Ссылка в URL не корректного формата",
			want: http.StatusBadRequest,
			url:  fmt.Sprintf("%s/%s", heplers.GetMainLink(), "123456/789"),
		},
		{
			name:  "Ссылка в URL корректного формата",
			want:  http.StatusTemporaryRedirect,
			url:   linkWithID,
			store: notEmptyStore,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			c := mocks.MockGinContext(w, r, tt.store)
			Resolve(c)
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
