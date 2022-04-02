package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/genga911/yandexworkshop/internal/app/config"
	"github.com/genga911/yandexworkshop/internal/app/heplers"
	"github.com/genga911/yandexworkshop/internal/app/session"
	"github.com/gin-gonic/gin"
)

const AuthIDLength = 16

func Auth(helper heplers.EasyCrypto, params *config.Params) gin.HandlerFunc {
	return func(c *gin.Context) {
		// проверим, есть ли у пользователя кука
		authCookieName := "auth"
		authCookie, err := c.Cookie(authCookieName)

		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			fmt.Println(err)
			return
		}

		var userID string
		if len(authCookie) < AuthIDLength {
			userID = heplers.ShortCode(AuthIDLength)
			encrypted := string(helper.Encode(userID))
			c.SetCookie(authCookieName, encrypted, params.CookieTTL, "/", params.ServerAddress, false, false)
		} else {
			userID = string(helper.Decode(authCookie))
		}

		c.Set("session", session.NewSession(userID))

		c.Next()
	}
}
