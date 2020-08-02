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
	ClassList []*ClassInfo `json:"class_list"`
}

// 增加Type，用于前端颜色分配
// 课表-课堂信息，用于response
type ClassInfo struct {
	CourseId  string           `json:"course_id"`
	ClassId   string           `json:"class_id"` // 教学班编号
	ClassName string           `json:"class_name"`
	Teacher   string           `json:"teacher"`
	Places    []string         `json:"places"`
	Times     []*ClassTimeInfo `json:"times"`
	Type      int8             `json:"type"` // 0-通必,1-专必,2-专选,3-通选,4-专业课,5-通核
}

type ClassTimeInfo struct {
	Start     int8   `json:"start"`      // 开始节数
	Duration  int8   `json:"duration"`   // 持续节数，若为1，则该课占两节
	Day       int8   `json:"day"`        // 星期
	Weeks     string `json:"weeks"`      // 周次，2-17
	WeekState int8   `json:"week_state"` // 全周0,单周1,双周2
}
