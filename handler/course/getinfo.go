package course

import (
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"github.com/gin-gonic/gin"
)

type TagList struct {
	Id   uint32 `json:"id"`
	Data OneTag `json:"data"`
}

type OneTag struct {
	Name string `json:"name"`
	Num  uint32 `json:"num"`
}

type TPList struct {
	Id    uint32
	Time  string
	Place string
	Week  string
}

type TTPL struct {
	Id      uint32   `json:"id"`
	ClassId string   `json:"class_id"`
	List    []TPList `json:"list"`
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
	CourseType     string            `json:"course_type"` // using or history
	CourseCredit   float32           `json:"course_credit"`
	Rate           float32           `json:"rate"`
	StarsNum       uint32            `json:"stars_num"`
	Attendance     map[string]uint32 `json:"attendance"`
	Exam           map[string]uint32 `json:"exam"`
	ClassInfo      []TTPL            `json:"class_info"`
	Tag            []TagList         `json:"tag"`
	LikeState      bool              `json:"like_state"`
}

type HistoryResponseInfo struct {
	CourseName     string            `json:"course_name"`
	TeacherName    string            `json:"teacher_name"`
	CourseCategory uint8             `json:"course_category"`
	CourseType     string            `json:"course_type"` // using or history
	Rate           float32           `json:"rate"`
	StarsNum       uint32            `json:"stars_num"`
	Attendance     map[string]uint32 `json:"attendance"`
	Exam           map[string]uint32 `json:"exam"`
	Tag            []TagList         `json:"tag"`
	LikeState      bool              `json:"like_state"`
}

func judge(classid string) string {
	if strings.Contains(classid, "课堂") {
		split := strings.Index(classid, "堂")
		var finstr string
		finstr = classid[split+4:] + classid[:split-3]
		return finstr
	}
	return classid
}

// 获取现用课程信息
func GetCourseInfo(c *gin.Context) {
	log.Info("GetInfo function is called")

	var n1 uint32

	// 获取hash
	hash := c.Param("hash")
	if hash == "" {
		log.Info("Get Param error")
		return
	}

	// userId获取及游客模式判断
	var userId uint32
	isLike := false

	userIdInterface, ok := c.Get("id")
	if ok {
		userId = userIdInterface.(uint32)
		log.Info("This User have token.")
		_, isLike = model.HasFavorited(userId, hash)
	}

	// 获取tag
	tags, err := model.GetCourseTags(hash)
	if err != nil {
		log.Info("Get Tags error")
	}
	tag := make([]TagList, 0)
	for n, v := range tags {
		a, err := model.GetTagNameById32(v.TagId)
		if err != nil {
			log.Info("Get TagNameById error")
		}
		b := v.Num
		in := OneTag{a, b}
		out := TagList{uint32(n + 1), in}
		tag = append(tag, out)
	}

	// 判断是否为历史课程
	course := &model.UsingCourseModel{Hash: hash}
	j, err := course.GetByHash2()
	if j == 1 {
		historyInfo := GetHistoryCourseInfo(hash, tag, isLike)
		handler.SendResponse(c, nil, historyInfo)
		log.Info("This is a HistoryCourse")
	} else {
		if err != nil {
			log.Info("course.GetByHash() error.")
			return
		}
		log.Info("This is a UsingCourse")

		courseid := course.CourseId

		// 获取所有课堂
		classid, err := model.GetAllClass(hash)
		if err != nil {
			log.Info("course.GetAllClass() error.")
		}

		class := &model.HistoryCourseModel{Hash: hash}
		if err := class.GetHistoryByHash(); err != nil {
			log.Info("course.GetHistoryByHash() error.")
		}

		// var score = map[uint32]uint32{70: 11, 7085: 76, 85: 13}

		var attendanceMap = service.GetAttendanceCheckTypeNumForCourseInfoEnglish(hash)

		var examMap = service.GetExamCheckTypeNumForCourseInfoEnglish(hash)
		// var test InfoClass
		test := make([]TTPL, 0, 60)

		//	var i int
		//	for i = 0; i <= len(classid); i++ {
		for _, v := range classid {
			list := make([]TPList, 0, 3)
			aclass := &model.UsingCourseModel{Hash: hash}
			if err := aclass.GetClass(courseid, v.ClassId); err != nil {
				log.Info("course.GetClass() error.")
				handler.SendError(c, err, nil, "")
				log.Info(courseid)
			}
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
			if a != "" || b != "" {
				list = append(list, list1)
			}
			if c != "" || d != "" {
				list = append(list, list2)
			}
			if e != "" || f != "" {
				list = append(list, list3)
			}
			// fmt.Print(list)
			if len(list) != 0 {
				n1++
				LIST := TTPL{n1, judge(v.ClassId), list}
				test = append(test, LIST)
				// fmt.Print(test)
			}
		}

		courseResponse := ResponseInfo{
			CourseName:     course.Name,
			TeacherName:    course.Teacher,
			CourseCategory: course.Type,
			CourseType:     "using",
			CourseCredit:   course.Credit,
			Rate:           class.Rate,
			StarsNum:       class.StarsNum,
			Attendance:     attendanceMap,
			Exam:           examMap,
			ClassInfo:      test,
			Tag:            tag,
			LikeState:      isLike,
		}

		handler.SendResponse(c, nil, courseResponse) // *
	}
}

// 获取历史课程信息
func GetHistoryCourseInfo(hash string, tag []TagList, isLike bool) HistoryResponseInfo {
	course := &model.HistoryCourseModel{Hash: hash}
	if err := course.GetHistoryByHash(); err != nil {
		log.Info("course.GetHistoryByHash() error.")
	}

	var attendanceMap = service.GetAttendanceCheckTypeNumForCourseInfoEnglish(hash)

	var examMap = service.GetExamCheckTypeNumForCourseInfoEnglish(hash)

	courseResponse := HistoryResponseInfo{
		CourseName:     course.Name,
		TeacherName:    course.Teacher,
		CourseCategory: course.Type,
		CourseType:     "history",
		Rate:           course.Rate,
		StarsNum:       course.StarsNum,
		Attendance:     attendanceMap,
		Exam:           examMap,
		Tag:            tag,
		LikeState:      isLike,
	}

	return courseResponse
}
