package util

import (
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"time"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
	}
	return ""
}

func GetCurrentTime() *time.Time {
	var t time.Time
	t =  time.Now()
	return &t
}
