package service

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/lexkong/log"
	"sync"
)

// 现用课堂信息
type CourseInfoForUsing struct {
	Id         uint32   `json:"id"`         //主键
	Hash       string   `json:"hash"`       //教师名和课程hash成的唯一标识，用于getinfo
	CourseId   string   `json:"course_id"`  //仅用于在UI上进行展示
	Name       string   `json:"name"`       //课程名称
	Teacher    string   `json:"teacher"`    //教师姓名
	Rate       float32  `json:"rate"`       //课程评价星级
	StarsNum   uint32   `json:"stars_num"`  //参与评分人数
	Attendance string   `json:"attendance"` //点名方式
	Exam       string   `json:"exam"`       //考核方式
	Tags       []string `json:"tags"`       //前二的tag
}

// 搜索查询的课程列表（历史课堂信息）
type CourseInfoForHistory struct {
	Id         uint32   `json:"id"`   //数据库表中记录的id，自增id
	Hash       string   `json:"hash"` //教师名和课程hash成的唯一标识，用于getinfo
	Name       string   `json:"name"`
	Teacher    string   `json:"teacher"`
	Rate       float32  `json:"rate"`
	StarsNum   uint32   `json:"stars_num"`
	Attendance string   `json:"attendance"` //点名方式
	Exam       string   `json:"exam"`       //考核方式
	Tags       []string `json:"tags"`       //前二的tag
}

// 获取点名方式
func GetAttendanceTypeMax(hash string) string {
	i := model.GetAttendanceType(hash)
	return GetAttendanceCheckTypeByCode(i)
}

// 获取期末考核方式
func GetExamCheckTypeMax(hash string) string {
	i := model.GetExamCheckType(hash)
	return GetExamCheckTypeByCode(i)
}

func kwReplace(kw string) string {
	if res, existed := KeywordMap[kw]; existed {
		kw = res
	}
	return kw
}

// Using
func SearchCourses(keyword string, page, limit uint64, t, a, w, p string) ([]CourseInfoForUsing, error) {
	keyword = kwReplace(keyword)
	courseRows, _ := model.AgainstAndMatchCourses(keyword, page, limit, t, a, w, p)

	courses := make([]CourseInfoForUsing, len(courseRows))
	locker := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i, row := range courseRows {
		wg.Add(1)
		go func(index int, row model.UsingCourseSearchModel) {
			result, err := model.GetTwoMostTagNamesOfCourseByHashId(row.Hash)
			if err != nil {
				log.Error("Get Tag Name error", err)
			}
			attendance := GetAttendanceTypeMax(row.Hash)
			exam := GetExamCheckTypeMax(row.Hash)

			locker.Lock()
			defer locker.Unlock()
			courses[index] = CourseInfoForUsing{
				Id:         row.Id,
				Hash:       row.Hash,
				Name:       row.Name,
				CourseId:   row.CourseId,
				Teacher:    row.Teacher,
				Rate:       row.Rate,
				StarsNum:   row.StarsNum,
				Attendance: attendance,
				Exam:       exam,
				Tags:       result,
			}
			wg.Done()
		}(i, row)
	}
	wg.Wait()

	return courses, nil
}

// History
func SearchHistoryCourses(keyword string, page, limit uint64, t string) ([]CourseInfoForHistory, error) {
	keyword = kwReplace(keyword)
	courseRows, _ := model.AgainstAndMatchHistoryCourses(keyword, page, limit, t)

	courses := make([]CourseInfoForHistory, len(courseRows))
	locker := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i, row := range courseRows {
		wg.Add(1)
		go func(index int, row model.HistoryCourseModel) {
			result, err := model.GetTwoMostTagNamesOfCourseByHashId(row.Hash)
			attendance := GetAttendanceTypeMax(row.Hash)
			exam := GetExamCheckTypeMax(row.Hash)
			if err != nil {
				log.Error("Get Tag Name error", err)
			}

			locker.Lock()
			defer locker.Unlock()
			courses[index] = CourseInfoForHistory{
				Id:         row.Id,
				Hash:       row.Hash,
				Name:       row.Name,
				Teacher:    row.Teacher,
				Rate:       row.Rate,
				StarsNum:   row.StarsNum,
				Attendance: attendance,
				Exam:       exam,
				Tags:       result,
			}
			wg.Done()
		}(i, row)
	}
	wg.Wait()

	return courses, nil
}

// Using
func GetAllCourses(page, limit uint64, t, a, w, p string) ([]CourseInfoForUsing, error) {
	courseRows, err := model.AllCourses(page, limit, t, a, w, p)
	if err != nil {
		return nil, err
	}

	courses := make([]CourseInfoForUsing, len(courseRows))
	locker := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i, row := range courseRows {
		wg.Add(1)
		go func(index int, row model.UsingCourseSearchModel) {
			result, err := model.GetTwoMostTagNamesOfCourseByHashId(row.Hash)
			if err != nil {
				log.Error("Get Tag Name error", err)
			}
			attendance := GetAttendanceTypeMax(row.Hash)
			exam := GetExamCheckTypeMax(row.Hash)

			locker.Lock()
			defer locker.Unlock()
			courses[index] = CourseInfoForUsing{
				Id:         row.Id,
				Hash:       row.Hash,
				Name:       row.Name,
				Teacher:    row.Teacher,
				CourseId:   row.CourseId,
				Rate:       row.Rate,
				StarsNum:   row.StarsNum,
				Attendance: attendance,
				Exam:       exam,
				Tags:       result,
			}
			wg.Done()
		}(i, row)
	}
	wg.Wait()

	return courses, nil
}

// History
func GetAllHistoryCourses(page, limit uint64, t string) ([]CourseInfoForHistory, error) {
	courseRows, err := model.AllHistoryCourses(page, limit, t)
	if err != nil {
		return nil, err
	}

	courses := make([]CourseInfoForHistory, len(courseRows))
	locker := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i, row := range courseRows {
		wg.Add(1)
		go func(index int, row model.HistoryCourseModel) {
			result, err := model.GetTwoMostTagNamesOfCourseByHashId(row.Hash)
			if err != nil {
				log.Error("Get Tag Name error", err)
			}
			attendance := GetAttendanceTypeMax(row.Hash)
			exam := GetExamCheckTypeMax(row.Hash)

			locker.Lock()
			defer locker.Unlock()
			courses[index] = CourseInfoForHistory{
				Id:         row.Id,
				Hash:       row.Hash,
				Name:       row.Name,
				Teacher:    row.Teacher,
				Rate:       row.Rate,
				StarsNum:   row.StarsNum,
				Attendance: attendance,
				Exam:       exam,
				Tags:       result,
			}
			wg.Done()
		}(i, row)
	}
	wg.Wait()

	return courses, nil
}
