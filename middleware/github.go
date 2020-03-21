package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidatePayload(key []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("X-Hub-Signature")
		if signature == "" {
			c.AbortWithStatus(http.StatusTeapot)
			return
		}
		signature = strings.Replace(signature, "sha1=", "", 1)

		s, err := hex.DecodeString(signature)
		if err != nil {
			c.AbortWithStatus(http.StatusTeapot)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)

		mac := hmac.New(sha1.New, key)
		mac.Write(buf.Bytes())
		p := mac.Sum(nil)

		if !hmac.Equal(s, p) {
			c.AbortWithStatus(http.StatusTeapot)
			return
		}

		c.Request.Body = ioutil.NopCloser(buf)

		c.Next()
	}
}
