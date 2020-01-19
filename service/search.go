package service //改需求了。。。（待优化）

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/lexkong/log"
)

// 现用课堂信息
type CourseInfoForUsing struct {
	Id       uint32   `json:"id"`        //主键
	Hash     string   `json:"hash"`      //教师名和课程hash成的唯一标识，用于getinfo
	CourseId string   `json:"course_id"` //仅用于在UI上进行展示
	Name     string   `json:"name"`      //课程名称
	Teacher  string   `json:"teacher"`   //教师姓名
	Rate     float32  `json:"rate"`      //课程评价星级
	StarsNum uint32   `json:"stars_num"` //参与评分人数
	Tags     []string `json:"tags"`      //前二的tag
}

// 搜索查询的课程列表（历史课堂信息）
type CourseInfoForHistory struct {
	Id       uint32   `json:"id"`   //数据库表中记录的id，自增id
	Hash     string   `json:"hash"` //教师名和课程hash成的唯一标识，用于getinfo
	Name     string   `json:"name"`
	Teacher  string   `json:"teacher"`
	Rate     float32  `json:"rate"`
	StarsNum uint32   `json:"stars_num"`
	Tags     []string `json:"tags"` //前二的tag
}

func kwReplace(kw string) string {
	if res, existed := KeywordMap[kw]; existed {
		kw = res
	}
	return kw
}

func GetTag(CourseId string) ([]string, error) {
	tagIds, err := model.GetTwoMostTagIdsOfCourseByHashId(CourseId)
	if err != nil {
		log.Error("GetTwoMostTagsOfCourseByHashId function error", err)
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
	return result, nil
}

//Using
func SearchCourses(keyword string, page, limit uint64, t, a, w, p string) ([]CourseInfoForUsing, error) {
	keyword = kwReplace(keyword)
	courseRows, _ := model.AgainstAndMatchCourses(keyword, page, limit, t, a, w, p)
	courses := make([]CourseInfoForUsing, len(courseRows))
	for i, row := range courseRows {
		class := &model.HistoryCourseModel{Hash: row.Hash}
		if err := class.GetHistoryByHash(); err != nil {
			log.Info("course.GetHistoryByHash() error.")
		}

		result, err := GetTag(row.Hash)
		if err != nil {
			log.Error("GetTag error", err)
		}

		courses[i] = CourseInfoForUsing{
			Id:       row.Id,
			Hash:     row.Hash,
			Name:     row.Name,
			CourseId: row.CourseId,
			Teacher:  row.Teacher,
			Rate:     class.Rate,
			StarsNum: class.StarsNum,
			Tags:     result,
		}
	}
	return courses, nil
}

//History
func SearchHistoryCourses(keyword string, page, limit uint64, t string) ([]CourseInfoForHistory, error) {
	keyword = kwReplace(keyword)
	courseRows, _ := model.AgainstAndMatchHistoryCourses(keyword, page, limit, t)
	courses := make([]CourseInfoForHistory, len(courseRows))
	for i, row := range courseRows {
		result, err := GetTag(row.Hash)
		if err != nil {
			log.Error("GetTag error", err)
		}

		courses[i] = CourseInfoForHistory{
			Id:       row.Id,
			Hash:     row.Hash,
			Name:     row.Name,
			Teacher:  row.Teacher,
			Rate:     row.Rate,
			StarsNum: row.StarsNum,
			Tags:     result,
		}
	}
	return courses, nil
}

//Using
func GetAllCourses(page, limit uint64, t, a, w, p string) ([]CourseInfoForUsing, error) {
	courseRows, err := model.AllCourses(page, limit, t, a, w, p)
	if err != nil {
		return nil, err
	}
	courses := make([]CourseInfoForUsing, len(courseRows))
	for i, row := range courseRows {
		class := &model.HistoryCourseModel{Hash: row.Hash}
		if err := class.GetHistoryByHash(); err != nil {
			log.Info("course.GetHistoryByHash() error.")
		}

		result, err := GetTag(row.Hash)
		if err != nil {
			log.Error("GetTag error", err)
		}

		courses[i] = CourseInfoForUsing{
			Id:       row.Id,
			Hash:     row.Hash,
			Name:     row.Name,
			Teacher:  row.Teacher,
			CourseId: row.CourseId,
			Rate:     class.Rate,
			StarsNum: class.StarsNum,
			Tags:     result,
		}
	}
	return courses, nil
}

//History
func GetAllHistoryCourses(page, limit uint64, t string) ([]CourseInfoForHistory, error) {

	courseRows, err := model.AllHistoryCourses(page, limit, t)
	if err != nil {
		return nil, err
	}
	courses := make([]CourseInfoForHistory, len(courseRows))
	for i, row := range courseRows {
		result, err := GetTag(row.Hash)
		if err != nil {
			log.Error("GetTag error", err)
		}

		courses[i] = CourseInfoForHistory{
			Id:       row.Id,
			Hash:     row.Hash,
			Name:     row.Name,
			Teacher:  row.Teacher,
			Rate:     row.Rate,
			StarsNum: row.StarsNum,
			Tags:     result,
		}
	}
	return courses, nil
}
