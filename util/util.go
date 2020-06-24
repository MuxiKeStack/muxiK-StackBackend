package util

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
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
	var loc, _ = time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc)
	return &t
}

func ParseTime(t *time.Time) (string, string) {
	return t.Format("2006-01-02"), t.Format("15:04:05")
}

// 根据获取的教师字段提取教师名，[2006982005/张立荣,2006982022/费军](或2006982005/张立荣/教授,2006982022/费军/讲师)=>张立荣,费军
func GetTeachersSqStrBySplitting(s string) string {
	sqs := strings.Split(s, ",")
	var teachers []string
	for _, s := range sqs {
		teachers = append(teachers, strings.Split(s, "/")[1])
	}
	return strings.Join(teachers, ",")
}

// 根据课程号和教师名字符串hash
func HashCourseId(courseNumStr, teachers string) string {
	hash := md5.New()
	hash.Write([]byte(courseNumStr + teachers))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
