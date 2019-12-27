package model

type CourseListModel struct {
	Id           uint32 `gorm:"column:id"`
	UserId       uint32 `gorm:"column:user_id"`
	CourseHashId string `gorm:"column:course_hash_id"`
}

type CourseInfoInCollections struct {
	CourseId   string                    `json:"course_id"` // 课程hash id
	CourseName string                    `json:"course_name"`
	ClassSum   int                       `json:"class_sum"` // 课堂数
	Classes    *[]ClassInfoInCollections `json:"classes"`
}

// 选课清单内的课堂（教学班）信息
type ClassInfoInCollections struct {
	ClassId         string                        `json:"class_id"` // 课堂hash id
	ClassName       string                        `json:"class_name"`
	TeachingClassId uint64                        `json:"teaching_class_id"` // 教学班编号
	Teacher         string                        `json:"teacher"`
	Times           *[]ClassTimeInfoInCollections `json:"times"`
	Places          *[]string                     `json:"places"`
}

type ClassTimeInfoInCollections struct {
	Time      string `json:"time"`
	Day       int8   `json:"day"`
	Weeks     string `json:"weeks"`
	WeekState int8   `json:"week_state"` // 全周0,单周1,双周2
}

// 选课清单页面的课程信息
type CourseInfoForCollections struct {
	Id                  uint32    `json:"id"` // 数据库表中记录的id，自增id
	CourseId            string    `json:"course_id"`
	CourseName          string    `json:"course_name"`
	Teacher             string    `json:"teacher"`
	EvaluationNum       uint32    `json:"evaluation_num"`
	Rate                float32   `json:"rate"`
	AttendanceCheckType string    `json:"attendance_check_type"`
	ExamCheckType       string    `json:"exam_check_type"`
	Tags                *[]string `json:"tags"`
}
