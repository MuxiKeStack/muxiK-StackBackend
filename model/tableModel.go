package model

type ClassTableModel struct {
	Id      uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	UserId  uint32 `gorm:"column:user_id"`
	Courses string `gorm:"column:courses"`
}

type ClassTableInfo struct {
	TableId   uint32       `json:"table_id"`
	TableName string       `json:"table_name"`
	Rank      uint8        `json:"rank"`
	ClassNum  uint32       `json:"class_num"`
	ClassList *[]ClassInfo `json:"class_list"`
}

type ClassInfo struct {
	CourseId   string           `json:"course_id"`
	ClassId    string           `json:"class_id"`
	CourseName string           `json:"course_name"`
	ClassName  string           `json:"class_name"`
	Teacher    string           `json:"teacher"`
	Places     *[]string        `json:"places"`
	Times      *[]ClassTimeInfo `json:"times"`
}

type ClassTimeInfo struct {
	Start  int8    `json:"start"`
	During int8    `json:"during"`
	Day    int8    `json:"day"`
	Weeks  *[]int8 `json:"weeks"`
}
