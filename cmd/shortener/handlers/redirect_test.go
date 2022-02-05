package handlers

import (
	"github.com/genga911/yandexworkshop/cmd/shortener/storages"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	linkWithId := "http://example.com/123456"
	var emptyStore = storages.CreateLinkStorage()
	var notEmptyStore = storages.CreateLinkStorage()
	notEmptyStore.Create(linkWithId)

	tests := []struct {
		name  string
		want  int
		url   string
		store *storages.LinkStorage
	}{
		{
			name:  "Нет ссылки в URL",
			want:  http.StatusNotFound,
			url:   "http://example.com",
			store: emptyStore,
		},
		{
			name: "Ссылка в URL не корректного формата",
			want: http.StatusNotFound,
			url:  "http://example.com/123456/789",
		},
		{
			name:  "Ссылка в URL корректного формата",
			want:  http.StatusNotFound,
			url:   linkWithId,
			store: notEmptyStore,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			Redirect(tt.store, w, request, 0)
			res := w.Result()
			_, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Ошибка %v", err)
			}
			assert.True(t, tt.want == res.StatusCode)
		})
	}
}

func Test_getShortLink(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Пустой url",
			url:  "http://example.com",
			want: "",
		},
		{
			name: "Корректный url c id",
			url:  "http://example.com/aAbBcC",
			want: "aAbBcC",
		},
		{
			name: "Не корректный url c id из цифр",
			url:  "http://example.com/123456",
			want: "",
		},
		{
			name: "Не корректный url",
			url:  "http://example.com/123456/789",
			want: "",
		},
		{
			name: "Не корректный url без цифр",
			url:  "http://example.com/aAbBcC/abs",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.url, nil)
			res := getShortLink(request)
			assert.Equal(t, tt.want, res)
		})
	}
}
