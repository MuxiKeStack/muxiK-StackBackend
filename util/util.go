package util

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"strings"
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
	t = time.Now()
	return &t
}

func GetTeachersSqStrBySplitting(s string) string {
	sqs := strings.Split(s, ",")
	var teachers []string
	for _, s := range sqs {
		teachers = append(teachers, strings.Split(s, "/")[1])
	}
	return strings.Join(teachers, ",")
}

func HashCourseId(courseNumStr, teachers string) string {
	//fmt.Println(courseNumStr, teachers)
	hash := md5.New()
	hash.Write([]byte(courseNumStr + teachers))
	//id := hash.Sum(nil)

	return fmt.Sprintf("%x", hash.Sum(nil))
}
