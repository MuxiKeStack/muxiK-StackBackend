package script

/*
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var (
	DB     gorm.DB
	DBAddr string
	DBUser string
	DBPwd  string
)

func init() {
	// 配置环境变量
	// export MUXIKSTACK_DB_ADDR=127.0.0.1:3306
	// export MUXIKSTACK_DB_USERNAME=muxi
	// export MUXIKSTACK_DB_PASSWORD=muxi
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MUXIKSTACK")
	DBAddr = viper.GetString("DB_ADDR")
	DBUser = viper.GetString("DB_USERNAME")
	DBPwd = viper.GetString("DB_PASSWORD")
}

func fill(courseId string) string {
	for i := len(courseId); i <= 8; i++ {
		courseId = courseId + "0"
	}
	return courseId
}

func judge1(c string) uint8 {
	switch c {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "5":
		return 5
	default:
		return 4
	}
}

type Outside struct {
	Code float64 `json:"code"`
	Msg  string  `json:"msg"`
	Data Middle  `json:"data"`
}

type Middle struct {
	PageNum  float64      `json:"pagenum"`
	PageSize float64      `json:"pageSize"`
	Size     float64      `json:"size"`
	StartRow float64      `json:"startRow"`
	EndRow   float64      `json:"endRow"`
	Total    float64      `json:"total"`
	Pages    float64      `json:"pages"`
	List     [1200]Inside `json:"list"`
}

type Inside struct {
	SiteId   string `json:"siteId"`
	CourseId string `json:"courseCode"`
	Teacher  string `json:"teacherName"`
	Name     string `json:"courseName"`
	Team     string `json:"termName"`
	Domain   string `json:"domainName"`
}

func main() {
	if DBAddr == "" || DBUser == "" {
		fmt.Println("Database config error, required env settings")
		return
	}
	dbOpenCmd := fmt.Sprintf("%s:%s@(%s)/muxikstack?charset=utf8&parseTime=True", DBUser, DBPwd, DBAddr)
	db, err := gorm.Open("mysql", dbOpenCmd)
	//db, err := gorm.Open("mysql", "*:*@(*.*.*.*:*)/muxikstack?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	fmt.Println("connection succeed")
	fmt.Println("Start crawling and importing...")

	var count int  // 计数
	for i1 := 0; i1 < 28; i1++ {
		strI := strconv.Itoa(i1 + 1)
		resp, err := http.PostForm("http://spoc.ccnu.edu.cn/courseCenterController/fuzzyQuerySitesByConditions", url.Values{"pageNum": {strI}, "pageSize": {"1200"}})
		if err != nil {
			fmt.Println(err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var b Outside
		if err := json.Unmarshal(body, &b); err != nil {
			fmt.Println(err)
			continue
		}
		// fmt.Printf("%+v\n",b)
		// fmt.Println(b.Data.List[5].Name)
		//defer resp.Body.Close()
		for i2 := 0; i2 < 1200; i2++ {
			count++
			fmt.Printf("正在载入第  %d  个课程...\r", count)
			//fmt.Print(i2, " ")
			courseId := fill(b.Data.List[i2].CourseId)
			//fmt.Println(courseId)
			teacher := b.Data.List[i2].Teacher
			name := b.Data.List[i2].Name
			key := util.HashCourseId(courseId, teacher)

			oneCourse := &model.HistoryCourseModel{
				Hash:    key,
				Name:    name,
				Teacher: teacher,
				Type:    judge1(courseId[3:4]),
			}
			d := db.Where("hash = ?", key).First(&oneCourse)
			if d.RecordNotFound() {
				db.Create(oneCourse)
			}
		}
		fmt.Println((i1+1)*1200, 33600-(i1+1)*1200)
		time.Sleep(time.Duration(2) * time.Second)
	}
}*/
