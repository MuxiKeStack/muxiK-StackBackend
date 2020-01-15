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
	Id    uint32
	Time  string
	Place string
	Week  string
}

type TTPL struct {
	Id   uint32   `json:"id"`
	List []TPList `json:"list"`
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
	Score          map[uint32]uint32 `json:"score"`
	ScoreNum       uint32            `json:"score_num"`
	Attendance     map[string]uint32 `json:"attendance"`
	Exam           map[string]uint32 `json:"exam"`
	ClassInfo      []TTPL            `json:"class_info"`
	TotalScore     float32           `json:"total_score"`
	OrdinaryScore  float32           `json:"ordinary_score"`
	CourseFeature1 uint32            `json:"course_feature_1"`
	CourseFeature2 uint32            `json:"course_feature_2"`
	CourseFeature3 uint32            `json:"course_feature_3"`
	CourseFeature4 uint32            `json:"course_feature_4"`
	CourseFeature5 uint32            `json:"course_feature_5"`
	CourseFeature6 uint32            `json:"course_feature_6"`
}

//获取课程信息
func GetCourseInfo(c *gin.Context) {
	log.Info("GetInfo function is called")

	var n1 uint32

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

	var score = map[uint32]uint32{70: 11, 7085: 76, 85: 13}

	var attendanceMap = service.GetAttendanceCheckTypeNumForCourseInfoEnglish(hash)

	var examMap = service.GetExamCheckTypeNumForCourseInfoEnglish(hash)
	//var test InfoClass
	test := make([]TTPL, 0, 60)

	tag1, _ := model.GetTagsNumber(1, course.CourseId)
	tag2, _ := model.GetTagsNumber(2, course.CourseId)
	tag3, _ := model.GetTagsNumber(3, course.CourseId)
	tag4, _ := model.GetTagsNumber(4, course.CourseId)
	tag5, _ := model.GetTagsNumber(5, course.CourseId)
	tag6, _ := model.GetTagsNumber(6, course.CourseId)

	var i int
	for i = 1; i <= 10; i++ {
		list := make([]TPList, 0, 3)
		//var list ClassList
		//var list2, list3 TPList
		//list1 := make([]TPList, 2)
		//list2 := make([]TPList, 2)
		//list3 := make([]TPList, 2)
		aclass := &model.UsingCourseModel{Hash: hash}
		if err := aclass.GetClass("45677654", uint64(i)); err != nil {
			log.Info("course.GetClass() error.")
			handler.SendError(c, err, nil, "")
			log.Info(courseid)
			//log.Println("%s", courseid)
			//fmt.Print(courseid, uint64(i))
		}
		//log.Info(aclass.Time1)
		//log.Info(aclass.Place1)
		a := aclass.Time1
		b := aclass.Place1
		x := aclass.Weeks1
		c := aclass.Time2
		d := aclass.Place2
		y := aclass.Weeks2
		e := aclass.Time3
		f := aclass.Place3
		z := aclass.Weeks3
		list1 := TPList{1, a, b, x}
		list2 := TPList{2, c, d, y}
		list3 := TPList{3, e, f, z}
		//fmt.Print(list4)
		//list1.time = aclass.Time1
		//list1.place = aclass.Place1
		if a != "" || b != "" {
			list = append(list, list1)
		}
		if c != "" || d != "" {
			list = append(list, list2)
		}
		if e != "" || f != "" {
			list = append(list, list3)
		}
		//fmt.Print(list)
		if len(list) != 0 {
			n1++
			LIST := TTPL{n1, list}
			test = append(test, LIST)
			//fmt.Print(test)
		}
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
		Score:          score,
		ScoreNum:       89,
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
