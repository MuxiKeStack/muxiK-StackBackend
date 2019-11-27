package service

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

// 课堂信息
type SearchCourseInfo struct {
	Id         uint32  `json:"id"`          //主键
	Name       string  `json:"name"`        //课程名称
	Credit     float32 `json:"credit"`      //学分
	Teacher    string  `json:"teacher"`     //任课教师姓名
	CourseId   string  `json:"course_id"`   //课程编号
	ClassId    uint64  `json:"class_id"`    //课堂编号
	Type       uint8   `json:"type"`        //课程类型
	CreditType uint8   `json:"credit_type"` //学分类型
}

// 历史课堂信息
type SearchHistoryCourseInfo struct {
	Id       uint32  `json:"id"`        //主键
	Hash     string  `json:"hash"`      //教师名和课程hash成的唯一标识
	Name     string  `json:"name"`      //课程名称
	Teacher  string  `json:"teacher"`   //教师姓名
	Type     uint8   `json:"type"`      //课程类型，公共课为0，专业课为1
	Rate     float32 `json:"rate"`      //课程评价星级
	StarsNum uint32  `json:"stars_num"` //参与评分人数
	Credit   float32 `json:"credit"`    //学分
}

func SearchCourses(keyword string, page, limit uint64, th bool) ([]SearchCourseInfo, error) {
	courses := make([]SearchCourseInfo, 0)
	courseRows, err := model.AgainstAndMatchCourses(keyword, page, limit, th)
	if err != nil {
		return courses, err
	}
	defer courseRows.Close()

	for courseRows.Next() {
		var course SearchCourseInfo
		courseRows.Scan(&course)
		courses = append(courses, course)
	}
	return courses, nil
}

func SearchHistoryCourses(keyword string, page, limit uint64) ([]SearchHistoryCourseInfo, error) {
	courses := make([]SearchHistoryCourseInfo, 0)
	courseRows, err := model.AgainstAndMatchHistoryCourses(keyword, page, limit)
	if err != nil {
		return courses, err
	}
	defer courseRows.Close()

	for courseRows.Next() {
		var course SearchHistoryCourseInfo
		courseRows.Scan(&course)
		courses = append(courses, course)
	}
	return courses, nil
}

func GetAllCourses(page, limit uint64, th bool) ([]SearchCourseInfo, error) {
	courseRows, err := model.AllCourses(page, limit, th)
	if err != nil {
		return nil, err
	}
	courses := make([]SearchCourseInfo, len(courseRows))
	for i, row := range courseRows {
		courses[i] = SearchCourseInfo{
			Id:         row.Id,
			Name:       row.Name,
			Credit:     row.Credit,
			Teacher:    row.Teacher,
			CourseId:   row.CourseId,
			ClassId:    row.ClassId,
			Type:       row.Type,
			CreditType: row.CreditType,
		}
	}
	return courses, nil
}

func GetAllHistoryCourses(page, limit uint64) ([]SearchHistoryCourseInfo, error) {
	courseRows, err := model.AllHistoryCourses(page, limit)
	if err != nil {
		return nil, err
	}
	courses := make([]SearchHistoryCourseInfo, len(courseRows))
	for i, row := range courseRows {
		courses[i] = SearchHistoryCourseInfo{
			Id:       row.Id,
			Hash:     row.Hash,
			Name:     row.Name,
			Teacher:  row.Teacher,
			Type:     row.Type,
			Rate:     row.Rate,
			StarsNum: row.StarsNum,
			Credit:   row.Credit,
		}
	}
	return courses, nil
}
