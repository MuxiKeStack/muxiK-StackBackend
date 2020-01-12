package service

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/lexkong/log"
)

type ProducedCourseItem struct {
	CourseId     string `json:"course_id"`
	Name         string `json:"name"`
	Teacher      string `json:"teacher"`
	Academic     string `json:"academic"`
	HasEvaluated bool   `json:"has_evaluated"`
}

func GetSelfCourseList(userId uint32, sid, pwd, year, term string) (*[]ProducedCourseItem, error) {
	originalCourses, err := util.GetSelfCoursesFromXK(sid, pwd, year, term)
	if err != nil {
		log.Error("GetSelfCoursesFromXK function error", err)
		return nil, err
	}

	var list []ProducedCourseItem
	//(*originalCourses.Items)[0].Jsxx = "2006982627/葛非,2006982646/彭熙,2006982670/刘明,2007980066/姚华雄"

	for _, item := range *originalCourses.Items {
		teacher := util.GetTeachersSqStrBySplitting(item.Jsxx)
		hashId := util.HashCourseId(item.Kch, teacher)
		info := ProducedCourseItem{
			CourseId:     hashId,
			Name:         item.Kcmc,
			Teacher:      teacher,
			Academic:     item.Kkxymc,
			HasEvaluated: model.HasEvaluated(userId, hashId),
		}
		list = append(list, info)
	}

	return &list, nil
}
