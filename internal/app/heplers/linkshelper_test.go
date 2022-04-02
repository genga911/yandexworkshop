package heplers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/genga911/yandexworkshop/internal/app/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetShortLink(t *testing.T) {
	userID := "test"
	tests := []struct {
		name string
		url  string
		code string
		want string
	}{
		{
			name: "Пустой url",
			url:  "http://example.com",
			code: "",
			want: "",
		},
		{
			name: "Корректный url c id",
			url:  "http://example.com/aAbBcC",
			code: "aAbBcC",
			want: "aAbBcC",
		},
		{
			name: "Не корректный url c id из цифр",
			url:  "http://example.com/123456",
			code: "123456",
			want: "",
		},
		{
			name: "Не корректный url",
			url:  "http://example.com/123456/789",
			code: "123456",
			want: "",
		},
		{
			name: "Корректный url без цифр, главное что параметр читаем",
			url:  "http://example.com/aAbBcC/abs",
			code: "aAbBcC",
			want: "aAbBcC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			var params []gin.Param
			params = append(params, gin.Param{
				Key:   "code",
				Value: tt.code,
			})

			c := mocks.MockGinContext(userID, w, r, params)

			res, err := GetShortLink(c)

			// проверяем что если пришла ошибка, её текст говоит о валидации
			if err != nil {
				assert.Equal(t, err.Error(), "validation error")
			} else {
				// иначе сравниваем результат
				assert.Equal(t, tt.want, res)
			}
		})
	}
}
