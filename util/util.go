package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// func GenShortId() (string, error) {
// 	return shortid.Generate()
// }

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

// 获取当前时间，东八区
func GetCurrentTime() *time.Time {
	// loc, _ := time.LoadLocation("Asia/Shanghai")
	// loc := time.FixedZone("CST", 8*3600)
	t := time.Now().UTC().Add(8 * time.Hour)
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

type ParseClassTimeItem struct {
	Start     int8
	End       int8
	Day       int8
	Weeks     string
	WeekState int8
}

// 解析上课时间详情
// time 格式：1-2#1 ==> 周一的第一到第二节，#后面的数字代表周几(1-7)
// week 格式：2-17#0 ==> 2-17周，全周；0为全周，1为单周，2为双周
func ParseClassTime(timeString, weekString string) ([]*ParseClassTimeItem, error) {
	// 分隔节次和星期
	timeSplits := strings.Split(timeString, "#")
	if len(timeSplits) < 2 {
		err := errors.New(fmt.Sprintf("split timeString %s failed", timeString))
		log.Error("split timeString failed, no #", err)
		return nil, err
	}

	// 上课星期
	day, err := strconv.Atoi(timeSplits[1])
	if err != nil {
		log.Error("strconv.Atoi error when parsing day", err)
		return nil, err
	}

	// 解析上课周次和单双周状态
	weekSplits := strings.Split(weekString, "#") // 分割周次和单双周状态
	weekState, err := strconv.Atoi(weekSplits[1])
	if err != nil {
		log.Error("strconv.Atoi error when parsing weekState", err)
		return nil, err
	}

	// 课堂节次
	// 可能有复数个节次区间，如 3-4,9-10
	// 2020-8-2 fix: time 出现个例 3-4,9-10#4
	multiSections := strings.Split(timeSplits[0], ",")

	var result []*ParseClassTimeItem

	for _, section := range multiSections {
		// section: 3-4
		split := strings.Split(section, "-")

		start, err := strconv.Atoi(split[0])
		if err != nil {
			log.Error("strconv.Atoi error when parsing start", err)
			return nil, err
		}

		end, err := strconv.Atoi(split[1])
		if err != nil {
			log.Error("strconv.Atoi error when parsing end", err)
			return nil, err
		}

		result = append(result, &ParseClassTimeItem{
			Start:     int8(start),
			End:       int8(end),
			Day:       int8(day),
			Weeks:     weekSplits[0],
			WeekState: int8(weekState),
		})
	}
	return result, nil
}
