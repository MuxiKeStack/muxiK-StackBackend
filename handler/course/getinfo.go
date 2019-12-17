package course

import (
	//	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type TPList struct {
	time  string
	place string
}

/*type ClassList struct {
	list2 *[]TPList
}

type InfoClass struct {
	list1 *[]ClassList
}*/

type ResponseInfo struct {
	CourseName     string
	TeacherName    string
	CourseCategory uint8
	CourseCredit   float32
	Rate           float32
	StarsNum       uint32
	Attendance     map[string]uint32
	Exam           map[string]uint32
	ClassInfo      [][]TPList
	CourseFeature1 uint32
	CourseFeature2 uint32
	CourseFeature3 uint32
	CourseFeature4 uint32
	CourseFeature5 uint32
	CourseFeature6 uint32
}

//获取课程信息
func GetCourseInfo(c *gin.Context) {
	/*	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		id32 := uint32(id)
		if err != nil {
			handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
			return
		}
		course := &model.UsingCourseModel{Id: id32}
		if err := course.GetById(); err != nil {
			log.Info("course.GetById() error.")
			return
		}

		handler.SendResponse(c, nil, course)*/

	hash := c.Param("hash")
	if hash == "" {
		log.Info("Get Param error")
		return
	}

	course := &model.UsingCourseModel{Hash: hash}
	if err := course.GetByHash(); err != nil {
		log.Info("course.GetByHash() error.")
		return
	}

	courseid := course.CourseId

	class := &model.HistoryCourseModel{Hash: hash}
	if err := class.GetHistoryByHash(); err != nil {
		log.Info("course.GetHistoryByHash() error.")
	}

	var attendanceMap = service.GetAttendanceCheckTypeNumForCourseInfo(course.CourseId)

	var examMap = service.GetExamCheckTypeNumForCourseInfo(course.CourseId)
	//var test InfoClass
	test := make([][]TPList, 10, 60)

	tag1, _ := model.GetTagsNumber(1, course.CourseId)
	tag2, _ := model.GetTagsNumber(2, course.CourseId)
	tag3, _ := model.GetTagsNumber(3, course.CourseId)
	tag4, _ := model.GetTagsNumber(4, course.CourseId)
	tag5, _ := model.GetTagsNumber(5, course.CourseId)
	tag6, _ := model.GetTagsNumber(6, course.CourseId)

	var i int
	for i = 0; i < 10; i++ {
		list := make([]TPList, 10, 20)
		//var list ClassList
		var list1, list2, list3 TPList
		//list1 := make([]TPList, 2)
		//list2 := make([]TPList, 2)
		//list3 := make([]TPList, 2)
		aclass := &model.UsingCourseModel{Hash: hash}
		if err := aclass.GetClass(courseid, uint64(i)); err != nil {
			log.Info("course.GetClass() error.")
		}
		list1.time = aclass.Time1
		list1.place = aclass.Place1
		list2.time = aclass.Time2
		list2.place = aclass.Place2
		list3.time = aclass.Time3
		list3.place = aclass.Place3
		list = append(list, list1, list2, list3)
		test = append(test, list)
	}

	courseResponse := ResponseInfo{
		CourseName:     course.Name,
		TeacherName:    course.Teacher,
		CourseCategory: course.Type,
		CourseCredit:   course.Credit,
		Rate:           class.Rate,
		StarsNum:       class.StarsNum,
		Attendance:     attendanceMap,
		Exam:           examMap,
		ClassInfo:      test,
		CourseFeature1: tag1, //model.GetTagsNumber(1, course.CourseId),
		CourseFeature2: tag2, //model.GetTagsNumber(2, course.CourseId),
		CourseFeature3: tag3, //model.GetTagsNumber(3, course.CourseId),
		CourseFeature4: tag4, //model.GetTagsNumber(4, course.CourseId),
		CourseFeature5: tag5, //model.GetTagsNumber(5, course.CourseId),
		CourseFeature6: tag6, //model.GetTagsNumber(6, course.CourseId),
	}

	handler.SendResponse(c, nil, courseResponse) //*
}
