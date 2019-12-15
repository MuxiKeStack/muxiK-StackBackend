package model

type CourseListModel struct {
	Id           uint32 `gorm:"column:id"`
	UserId       uint32 `gorm:"column:user_id"`
	CourseHashId string `gorm:"column:course_hash_id"`
}

type CourseInfoInCollections struct {
	CourseId   string                    `json:"course_id"`
	CourseName string                    `json:"course_name"`
	ClassSum   int                       `json:"class_sum"`
	Classes    *[]ClassInfoInCollections `json:"classes"`
}

type ClassInfoInCollections struct {
	ClassId         string                        `json:"class_id"`
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
