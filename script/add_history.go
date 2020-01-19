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
)

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
	}
	return 4
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
	db, err := gorm.Open("mysql", "**:*****@(**.**.**.***:****)/muxikstack?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connection succeed")

	var i1, i2 int
	for i1 = 0; i1 < 28; i1++ {
		stri := strconv.Itoa(i1 + 1)
		resp, err := http.PostForm("http://spoc.ccnu.edu.cn/courseCenterController/fuzzyQuerySitesByConditions", url.Values{"pageNum": {stri}, "pageSize": {"1200"}})
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var b Outside
		if err := json.Unmarshal([]byte(body), &b); err != nil {
			fmt.Println("json unmarshal error", err)
			return
		}
		// fmt.Printf("%+v\n",b)
		// fmt.Println(b.Data.List[5].Name)
		//defer resp.Body.Close()
		for i2 = 0; i2 < 1200; i2++ {
			fmt.Print(i2, " ")
			courseId := b.Data.List[i2].CourseId
			//fmt.Println(courseId)
			teacher := b.Data.List[i2].Teacher
			name := b.Data.List[i2].Name
			key := util.HashCourseId(courseId, teacher)

			onecourse := &model.HistoryCourseModel{
				Hash:    key,
				Name:    name,
				Teacher: teacher,
				Type:    judge1(courseId[3:4]),
			}
			d := db.Where("hash = ?", key).First(&onecourse)
			if d.RecordNotFound() {
				db.Create(onecourse)
			} else {
				continue
			}
		}
		fmt.Println((i1+1)*1200, 33600-(i1+1)*1200)
		time.Sleep(time.Duration(2) * time.Second)
	}
}
*/