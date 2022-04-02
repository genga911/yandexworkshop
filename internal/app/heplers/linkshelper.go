package heplers

import (
	"errors"
	"fmt"
	"path"
	"regexp"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/gin-gonic/gin"
)

// получим короткую ссылку из урл
func GetShortLink(c *gin.Context) (string, error) {
	// провалидируем урл, ожидаем только буквы как в константе пакета codehelper
	url := c.Param("code")
	matched, err := regexp.MatchString(`^[a-zA-Z]+$`, url)
	if err != nil || !matched {
		return "", errors.New("validation error")
	}

	return path.Base(url), nil
}

func PrepareShortLink(code string, cfg *config.Params) string {
	return fmt.Sprintf(
		"%s/%s",
		cfg.BaseURL,
		code,
	)
}
