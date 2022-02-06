package heplers

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	// подключаем тестовое окружение
	err := godotenv.Load(os.ExpandEnv("$GOPATH/yandexworkshop/.env.testing"))
	if err != nil {
		fmt.Println("Ошибка загрузки .env: " + err.Error())
	}
}

func TestGetShortLink(t *testing.T) {
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
			res, err := GetShortLink(request)

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

func TestGetServerAddress(t *testing.T) {

	tests := []struct {
		name    string
		want    string
		envHost string
		envPort string
	}{
		{
			name: "Корректно заданный env",
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := GetServerAddress()
			assert.Equal(t, address, tt.want)
		})
	}
}

func TestGetMainLink(t *testing.T) {
	tests := []struct {
		name        string
		want        string
		envHost     string
		envPort     string
		envProtocol string
	}{
		{
			name: "Корректно заданный env",
			want: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := GetMainLink()
			assert.Equal(t, address, tt.want)
		})
	}
}
