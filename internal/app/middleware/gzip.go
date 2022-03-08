package middleware

import (
	"compress/gzip"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func Gzip(c *gin.Context) {
	encoding := ""
	encodingFromHeaders := c.Request.Header["Content-Encoding"]
	if encodingFromHeaders != nil {
		encoding = c.Request.Header["Content-Encoding"][0]
	}

	var globError error

	if strings.Contains(encoding, "gzip") {
		body, bodyErr := ioutil.ReadAll(c.Request.Body)
		if bodyErr != nil {
			globError = bodyErr
		} else {
			reader, newReaderErr := gzip.NewReader(strings.NewReader(string(body)))
			if newReaderErr != nil {
				globError = newReaderErr
			}
			if reader != nil {
				c.Request.Body = reader
			}
		}
	}

	if globError == nil {
		c.Next()
	}
}
