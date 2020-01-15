package service

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/lexkong/log"
)

// 课堂信息
type SearchCourseInfo struct {
	Id       uint32  `json:"id"`        //主键
	Name     string  `json:"name"`      //课程名称
	Credit   float32 `json:"credit"`    //学分
	Teacher  string  `json:"teacher"`   //任课教师姓名
	CourseId string  `json:"course_id"` //课程编号
	ClassId  string  `json:"class_id"`  //课堂编号
	Type     uint8   `json:"type"`      //课程类型
}

// 历史课堂信息
type SearchHistoryCourseInfo struct {
	Id       uint32   `json:"id"`        //主键
	Hash     string   `json:"hash"`      //教师名和课程hash成的唯一标识
	Name     string   `json:"name"`      //课程名称
	Teacher  string   `json:"teacher"`   //教师姓名
	Type     uint8    `json:"type"`      //课程类型，公共课为0，专业课为1
	Rate     float32  `json:"rate"`      //课程评价星级
	StarsNum uint32   `json:"stars_num"` //参与评分人数
	Credit   float32  `json:"credit"`    //学分
	Tags     []string `json:"tags"`      //前四的tag
}

// 选课助手的课程列表（现用课程信息）
type CourseInfoForAssistant struct {
	Id       uint32  `json:"id"` //数据库表中记录的id，自增id
	Name     string  `json:"course_name"`
	Academy  string  `json:"academy"`
	Teacher  string  `json:"teacher"`
	CourseId string  `json:"course_id"`
	Rate     float32 `json:"rate"`
	StarsNum uint32  `json:"stars_num"`
	Time1    string  `json:"time1"`
	Time2    string  `json:"time2"`
	Time3    string  `json:"time3"`
	Place1   string  `json:"place1"`
	Place2   string  `json:"place2"`
	Place3   string  `json:"place3"`
	Region   uint8   `json:"region"`
}

// 搜索查询的课程列表（历史课堂信息）
type CourseInfoForSearch struct {
	Id         uint32   `json:"id"` //数据库表中记录的id，自增id
	CourseName string   `json:"course_name"`
	Teacher    string   `json:"teacher"`
	Rate       float32  `json:"rate"`
	StarsNum   uint32   `json:"stars_num"`
	Tags       []string `json:"tags"`
}

func kwReplace(kw string) string {
	if res, existed := KeywordMap[kw]; existed {
		kw = res
	}
	return kw
}

func SearchCourses(keyword string, page, limit uint64, t, a, w, p string) ([]CourseInfoForAssistant, error) {
	keyword = kwReplace(keyword)
	courseRows, _ := model.AgainstAndMatchCourses(keyword, page, limit, t, a, w, p)
	/*	if err != nil {
			return courses, err
		}
		defer courseRows.Close()

		for courseRows.Next() {
			var course SearchCourseInfo
			courseRows.Scan(&course)
			courses = append(courses, course)
		}*/
	courses := make([]CourseInfoForAssistant, len(courseRows))
	for i, row := range courseRows {
		class := &model.HistoryCourseModel{Hash: row.Hash}
		if err := class.GetHistoryByHash(); err != nil {
			log.Info("course.GetHistoryByHash() error.")
		}

		courses[i] = CourseInfoForAssistant{
			Id:       row.Id,
			Name:     row.Name,
			Academy:  row.Academy,
			Teacher:  row.Teacher,
			CourseId: row.CourseId,
			Rate:     class.Rate,
			StarsNum: class.StarsNum,
			Time1:    row.Time1,
			Time2:    row.Time2,
			Time3:    row.Time3,
			Place1:   row.Place1,
			Place2:   row.Place2,
			Place3:   row.Place3,
			Region:   row.Region,
		}
	}
	return courses, nil
}

func SearchHistoryCourses(keyword string, page, limit uint64, t string) ([]SearchHistoryCourseInfo, error) {
	keyword = kwReplace(keyword)
	courseRows, _ := model.AgainstAndMatchHistoryCourses(keyword, page, limit, t)
	/*	if err != nil {
			return courses, err
		}
		defer courseRows.Close()

		for courseRows.Next() {
			var course SearchHistoryCourseInfo
			courseRows.Scan(&course)
			courses = append(courses, course)
		}*/
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

func GetAllCourses(page, limit uint64, t, a, w, p string) ([]CourseInfoForAssistant, error) {
	courseRows, err := model.AllCourses(page, limit, t, a, w, p)
	if err != nil {
		return nil, err
	}
	courses := make([]CourseInfoForAssistant, len(courseRows))
	for i, row := range courseRows {
		class := &model.HistoryCourseModel{Hash: row.Hash}
		if err := class.GetHistoryByHash(); err != nil {
			log.Info("course.GetHistoryByHash() error.")
		}

		courses[i] = CourseInfoForAssistant{
			Id:       row.Id,
			Name:     row.Name,
			Teacher:  row.Teacher,
			CourseId: row.CourseId,
			Rate:     class.Rate,
			StarsNum: class.StarsNum,
			Time1:    row.Time1,
			Time2:    row.Time2,
			Time3:    row.Time3,
			Place1:   row.Place1,
			Place2:   row.Place2,
			Place3:   row.Place3,
			Region:   row.Region,
		}
	}
	return courses, nil
}

func GetAllHistoryCourses(page, limit uint64, t string) ([]SearchHistoryCourseInfo, error) {

	courseRows, err := model.AllHistoryCourses(page, limit, t)
	if err != nil {
		return nil, err
	}
	courses := make([]SearchHistoryCourseInfo, len(courseRows))
	for i, row := range courseRows {
		course := &model.UsingCourseModel{Hash: row.Hash}
		if err := course.GetByHash(); err != nil {
			log.Info("course.GetByHash() error.")
			return nil, err
		}

		tagIds, err := model.GetFourMostTagIdsOfCourseByHashId(course.CourseId)
		if err != nil {
			log.Error("GetFourMostTagsOfCourseByHashId function error", err)
			return nil, err
		}

		var result []string
		for _, id := range tagIds {
			tag, err := model.GetTagNameById(id)
			if err != nil {
				log.Error("GetTagNameById function error", err)
				return nil, err
			}
			result = append(result, tag)
		}
		courses[i] = SearchHistoryCourseInfo{
			Id:       row.Id,
			Hash:     row.Hash,
			Name:     row.Name,
			Teacher:  row.Teacher,
			Type:     row.Type,
			Rate:     row.Rate,
			StarsNum: row.StarsNum,
			Credit:   row.Credit,
			Tags:     result,
		}
	}
	return courses, nil
}
