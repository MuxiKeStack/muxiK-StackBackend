package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lexkong/log"
	"golang.org/x/net/publicsuffix"
)

type OriginalGrade struct {
	Items []*GradeItem `json:"items"`
}

type GradeItem struct {
	JxbId string `json:"jxb_id"` // 教学id，用于请求平时成绩和期末成绩
	Kch   string `json:"kch"`    // 课程号
	Kcmc  string `json:"kcmc"`   // 课程名
	Jsxm  string `json:"jsxm"`   // 教师名
	Cj    string `json:"cj"`     // 成绩
	Xnm   string `json:"xnm"`
	Xqm   string `json:"xqm"`
}

type ResultGradeItem struct {
	Teacher    string  `json:"teacher"`
	CourseId   string  `json:"course_id"` // 课程号
	CourseName string  `json:"course_name"`
	TotalScore float32 `json:"total_score"` // 总评
	UsualScore float32 `json:"usual_grade"` // 平时成绩
	FinalScore float32 `json:"final_score"` // 期末成绩
}

// 从教务处获取成绩
func GetGradeFromXK(sid, password string, curRecordNum int) ([]*ResultGradeItem, bool, error) {
	params, err := MakeAccountPreflightRequest()
	if err != nil {
		log.Error("MakeAccountPreflightRequest function error", err)
		return nil, false, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, false, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}

	if err := MakeAccountRequest(sid, password, params, client); err != nil {
		log.Error("MakeAccountRequest function err", err)
		return nil, false, err
	}

	if err := MakeXKLogin(client); err != nil {
		log.Error("MakeXKLogin function error", err)
		return nil, false, err
	}

	formData := url.Values{}
	formData.Set("xnm", "")
	formData.Set("xqm", "")
	formData.Set("_search", "false")
	formData.Set("nd", string(time.Now().UnixNano()))
	formData.Set("queryModel.showCount", "1000")
	formData.Set("queryModel.currentPage", "1")
	formData.Set("queryModel.sortName", "")
	formData.Set("queryModel.sortOrder", "asc")
	formData.Set("time", "0")

	requestUrl := "http://xk.ccnu.edu.cn/jwglxt/cjcx/cjcx_cxDgXscj.html?doType=query&gnmkdm=N305005"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, false, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Origin", "http://xk.ccnu.edu.cn")
	req.Header.Set("Host", "xk.ccnu.edu.cn")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error", err)
		return nil, false, err

	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, false, err

	}

	var data = OriginalGrade{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("Json unmarshal failed, maybe login error, need to enter verification code", err)
		return nil, false, err
	}

	// 未发布新成绩成绩，直接返回
	if len(data.Items) <= curRecordNum {
		return nil, false, nil
	}

	var result []*ResultGradeItem

	for _, item := range data.Items {
		// 解析总成绩
		// 解析出现错误值：缓考
		totalScore, err := strconv.ParseFloat(item.Cj, 32)
		if err != nil {
			log.Errorf(err, "parse %s to float error", item.Cj)
			continue
		}

		// 获取平时和期末成绩
		usualScore, finalScore, err := GetUsualAndFinalGradeFromXK(client, item.JxbId, item.Kcmc, item.Xnm, item.Xqm)
		if err != nil {
			log.Errorf(err, "GetUsualAndFinalGrade for (%s, %s, %s) error", item.Kch, item.Kcmc, item.Jsxm)
			continue
		}

		result = append(result, &ResultGradeItem{
			Teacher:    item.Jsxm,
			CourseId:   item.Kch,
			CourseName: item.Kcmc,
			TotalScore: float32(totalScore),
			UsualScore: usualScore,
			FinalScore: finalScore,
		})
	}

	log.Info("Get grades successfully")
	return result, true, nil
}

// 发起请求，获取平时和期末成绩
func GetUsualAndFinalGradeFromXK(client *http.Client, jxbid, kcmc, xnm, xqm string) (float32, float32, error) {
	formData := url.Values{}
	formData.Set("xnm", xnm)
	formData.Set("xqm", xqm)
	formData.Set("jxb_id", jxbid)
	formData.Set("kcmc", kcmc)

	requestUrl := "http://xk.ccnu.edu.cn/jwglxt/cjcx/cjcx_cxCjxq.html?time=" + strconv.Itoa(int(time.Now().UnixNano())) + "&gnmkdm=N305005"

	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return 0, 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Origin", "http://xk.ccnu.edu.cn")
	req.Header.Set("Host", "xk.ccnu.edu.cn")
	req.Header.Set("Accept", "text/html, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error", err)
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	return ParseByRegexp(string(body))
}

// 正则解析HTML网页，获取平时成绩和期末成绩
func ParseByRegexp(bodyStr string) (float32, float32, error) {
	rg, err := regexp.Compile(`<td valign="middle">([0-9|\\.]*)&nbsp;</td>`)
	if err != nil {
		log.Error("regexp err", err)
		return 0, 0, err
	}

	result := rg.FindAllStringSubmatch(bodyStr, 2)
	if len(result) < 2 {
		return 0, 0, nil
	}

	var score [2]float64
	for i, r := range result {
		if i >= 2 {
			break
		}
		if len(r) < 2 {
			log.Infof("score body: %t", r)
			return 0, 0, errors.New("No score")
		}

		// 解析分数
		score[i], err = strconv.ParseFloat(r[1], 32)
		if err != nil {
			return 0, 0, err
		}

		// 分数为0，需要判错
		if score[i] == 0 {
			return 0, 0, errors.New("score is zero")
		}
	}

	return float32(score[0]), float32(score[1]), nil
}
