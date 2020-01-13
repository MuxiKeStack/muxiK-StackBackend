package model

// 搜索查询的课程列表
type CourseInfoForSearch struct {
	Id         uint32    `json:"id"` //数据库表中记录的id，自增id
	CourseName string    `json:"course_name"`
	Teacher    string    `json:"teacher"`
	Rate       float32   `json:"rate"`
	StarsNum   uint32    `json:"stars_num"`
	Tags       *[]string `json:"tags"`
}

// 选课助手的课程列表
type CourseInfoForAssistant struct {
	Id         uint32  `json:"id"` //数据库表中记录的id，自增id
	CourseName string  `json:"course_name"`
	Teacher    string  `json:"teacher"`
	CourseId   string  `json:"course_id"`
	Rate       float32 `json:"rate"`
	StarsNum   uint32  `json:"stars_num"`
	Time       string  `json:"time"`
	Place      string  `json:"place"`
	Region     uint8   `json:"region"`
}
