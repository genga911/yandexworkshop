package session

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type session struct {
	UserID string
}

const GuestSession = "guest"

func NewSession(userID string) *session {
	return &session{
		UserID: userID,
	}
}

func GetSession(c *gin.Context) *session {
	s, exist := c.Get("session")

	if !exist {
		fmt.Println("Сессия не найдена, пользователь зашел как гость")
		return &session{
			UserID: GuestSession,
		}
	}

	return s.(*session)
}
