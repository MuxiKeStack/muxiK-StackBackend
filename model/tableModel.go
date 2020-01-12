package model

type ClassTableModel struct {
	Id      uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	UserId  uint32 `gorm:"column:user_id"`
	Name    string `gorm:"column:name"`
	Classes string `gorm:"column:classes"`
}

type ClassTableInfo struct {
	TableId   uint32       `json:"table_id"`
	TableName string       `json:"table_name"`
	ClassNum  uint32       `json:"class_num"`
	ClassList *[]ClassInfo `json:"class_list"`
}

type ClassInfo struct {
	CourseId  string           `json:"course_id"`
	ClassId   uint64           `json:"class_id"` // 教学班编号
	ClassName string           `json:"class_name"`
	Teacher   string           `json:"teacher"`
	Places    *[]string        `json:"places"`
	Times     *[]ClassTimeInfo `json:"times"`
}

type ClassTimeInfo struct {
	Start     int8   `json:"start"`
	Duration  int8   `json:"duration"`
	Day       int8   `json:"day"`
	Weeks     string `json:"weeks"`
	WeekState int8   `json:"week_state"` // 全周0,单周1,双周2
}
