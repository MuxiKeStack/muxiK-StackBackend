package course

import (
	"fmt"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type TPList struct {
	Time  string
	Place string
}

/*type ClassList struct {
	list2 *[]TPList
}

type InfoClass struct {
	list1 *[]ClassList
}*/

type ResponseInfo struct {
	CourseName     string            `json:"course_name"`
	TeacherName    string            `json:"teacher_name"`
	CourseCategory uint8             `json:"course_category"`
	CourseCredit   float32           `json:"course_credit"`
	Rate           float32           `json:"rate"`
	StarsNum       uint32            `json:"stars_num"`
	Attendance     map[string]uint32 `json:"attendance"`
	Exam           map[string]uint32 `json:"exam"`
	ClassInfo      [][]TPList        `json:"class_info"`
	TotalScore     float32           `json:"total_score"`
	OrdinaryScore  float32           `json:"ordinary_score"`
	CourseFeature1 uint32            `json:"course_feature_1"`
	CourseFeature2 uint32            `json:"course_feature_2"`
	CourseFeature3 uint32            `json:"course_feature_3"`
	CourseFeature4 uint32            `json:"course_feature_4"`
	CourseFeature5 uint32            `json:"course_feature_5"`
	CourseFeature6 uint32            `json:"course_feature_6"`
}

/*
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendInfo(c *gin.Context, err error, data ResponseInfo) {
	code, message := errno.DecodeErr(err)

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}*/

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
	log.Info("GetInfo function is called")

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
	fmt.Println(courseid)

	class := &model.HistoryCourseModel{Hash: hash}
	if err := class.GetHistoryByHash(); err != nil {
		log.Info("course.GetHistoryByHash() error.")
	}

	var attendanceMap = service.GetAttendanceCheckTypeNumForCourseInfo(hash)

	var examMap = service.GetExamCheckTypeNumForCourseInfo(hash)
	//var test InfoClass
	test := make([][]TPList, 0, 60)

	tag1, _ := model.GetTagsNumber(1, course.CourseId)
	tag2, _ := model.GetTagsNumber(2, course.CourseId)
	tag3, _ := model.GetTagsNumber(3, course.CourseId)
	tag4, _ := model.GetTagsNumber(4, course.CourseId)
	tag5, _ := model.GetTagsNumber(5, course.CourseId)
	tag6, _ := model.GetTagsNumber(6, course.CourseId)

	var i int
	for i = 0; i < 10; i++ {
		list := make([]TPList, 0, 3)
		//var list ClassList
		var list2, list3 TPList
		//list1 := make([]TPList, 2)
		//list2 := make([]TPList, 2)
		//list3 := make([]TPList, 2)
		aclass := &model.UsingCourseModel{Hash: hash}
		if err := aclass.GetClass("45677654", 1); err != nil {
			log.Info("course.GetClass() error.")
			handler.SendError(c, err, nil, "")
			log.Info(courseid)
			//log.Println("%s", courseid)
			//fmt.Print(courseid, uint64(i))
		}
		//log.Info(aclass.Time1)
		//log.Info(aclass.Place1)
		a := aclass.Time1
		b := aclass.Time2
		list4 := TPList{a, b}
		//fmt.Print(list4)
		//list1.time = aclass.Time1
		//list1.place = aclass.Place1
		list2.Time = aclass.Time2
		list2.Place = aclass.Place2
		list3.Time = aclass.Time3
		list3.Place = aclass.Place3
		list = append(list, list4) //, list2, list3)
		//fmt.Print(list)
		//if len(list) != 0 {
		test = append(test, list)
		//fmt.Print(test)
		//}
	}
	//var test2 *[][]TPList
	//*test2 = make([][]TPList, 60)
	//test2 = &test
	//fmt.Print(test)

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
		TotalScore:     80.0,
		OrdinaryScore:  80.0,
		CourseFeature1: tag1,
		CourseFeature2: tag2,
		CourseFeature3: tag3,
		CourseFeature4: tag4,
		CourseFeature5: tag5,
		CourseFeature6: tag6,
	}

	// fmt.Print(courseResponse)

	//SendInfo(c, nil, courseResponse)
	handler.SendResponse(c, nil, courseResponse) //*
}
