package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"golang.org/x/net/publicsuffix"
)

type OriginalCourses struct {
	Items *[]OriginalCourseItem `json:"items" binding:"required"`
}

// 课程数据信息
type OriginalCourseItem struct {
	Kch    string `json:"kch" binding:"required"`  // 课程号
	Kcmc   string `json:"kcmc" binding:"required"` // 课程名称
	Jsxx   string `json:"jsxx" binding:"required"` // 教师信息，格式如：2008980036/宋冰玉/讲师
	Xnm    string `json:"xnm" binding:"required"`  // 学年名，如 2019
	Xqmc   string `json:"xqmc" binding:"required"` // 学期名称，如 1/2/3
	Kkxymc string `json:"kkxymc"`                  // 开课学院
	// Kclbmc string `json:"kclbmc"`                  // 课程类别名称，如公共课/专业课
	// Kcxzmc string `json:"kcxzmc"`                  // 课程性质，如专业主干课程/通识必修课
}

// 获取个人已上过的课程
func GetSelfCoursesFromXK(sid, password string, year, term string) (*OriginalCourses, error) {
	params, err := MakeAccountPreflightRequest()
	if err != nil {
		log.Error("MakeAccountPreflightRequest function error", err)
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}

	if err := MakeAccountRequest(sid, password, params, client); err != nil {
		log.Error("MakeAccountRequest function err", err)
		return nil, err
	}

	if err := MakeXKLogin(client); err != nil {
		log.Error("MakeXKLogin function error", err)
		return nil, err
	}

	courses, err := MakeCoursesGetRequest(client, "http://xk.ccnu.edu.cn/jwglxt/xkcx/xkmdcx_cxXkmdcxIndex.html?doType=query&gnmkdm=N255010&su="+sid, year, term)
	if err != nil {
		log.Error("MakeCoursesGetRequest function error", err)
		return nil, err
	}

	// 排除选课期间刚选的课
	if courses.Items != nil {
		var list []OriginalCourseItem

		m := "0"
		switch true {
		case time.Now().Month() > 5 && time.Now().Month() < 10:
			m = "1"
		case time.Now().Month() < 4 || time.Now().Month() > 11:
			m = "2"
		}

		for _, course := range *courses.Items {
			if course.Xnm != strconv.Itoa(time.Now().Year()) || course.Xqmc != m {
				list = append(list, course)
			}
		}
		courses.Items = &list
	}

	return courses, err
}

// 请求获取课程列表
func MakeCoursesGetRequest(client *http.Client, sid, year, term string) (*OriginalCourses, error) {
	var termMap = map[string]string{"1": "3", "2": "12", "3": "16"} // 学期参数
	if year == "0" {
		year = ""
	}

	formData := url.Values{}
	formData.Set("xnm", year)          // 学年名
	formData.Set("xqm", termMap[term]) // 学期名
	formData.Set("_search", "false")
	formData.Set("nd", string(time.Now().UnixNano()))
	formData.Set("queryModel.showCount", "1000")
	formData.Set("queryModel.currentPage", "1")
	formData.Set("queryModel.sortName", "")
	formData.Set("queryModel.sortOrder", "asc")
	formData.Set("time", "5")

	requestUrl := "http://xk.ccnu.edu.cn/jwglxt/xkcx/xkmdcx_cxXkmdcxIndex.html?doType=query&gnmkdm=N255010&su=" + sid
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Origin", "http://xk.ccnu.edu.cn")
	req.Header.Set("Host", "xk.ccnu.edu.cn")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var data = OriginalCourses{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Json unmarshal failed", err)
		return nil, err
	}

	return &data, nil
}
